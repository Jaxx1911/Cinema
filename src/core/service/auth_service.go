package service

import (
	"TTCS/src/common/constant"
	"TTCS/src/common/crypto"
	"TTCS/src/common/fault"
	"TTCS/src/common/mail"
	"TTCS/src/core/domain"
	"TTCS/src/core/dto"
	"TTCS/src/present/httpui/request"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo     domain.UserRepo
	otpRepo      domain.OtpRepo
	hashProvider crypto.HashProvider
	otpProvider  crypto.OTPProvider
	jwtProvider  crypto.JwtProvider
	mailService  *mail.GmailService
}

func NewAuthService(userRepo domain.UserRepo, otpRepo domain.OtpRepo, hashProvider crypto.HashProvider, otpProvider crypto.OTPProvider, jwtProvider crypto.JwtProvider, mailService *mail.GmailService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		otpRepo:      otpRepo,
		hashProvider: hashProvider,
		otpProvider:  otpProvider,
		jwtProvider:  jwtProvider,
		mailService:  mailService,
	}
}

func (s *AuthService) SignUpOTP(ctx context.Context, email string) (string, error) {
	caller := "AuthService.SignUpOTP"
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}
	if user != nil {
		err := errors.New("user already exists")
		return "", fault.Wrapf(err, "[%v] user already exists", caller).SetTag(fault.TagDuplicate).SetKey(fault.KeyUser)
	}
	otpDb, err := s.otpRepo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}
	if otpDb != nil {
		if otpDb.CreatedAt.After(time.Now().Add(-5 * time.Minute)) {
			err := errors.New("otp already exists")
			return "", fault.Wrapf(err, "[%v] otp already exists", caller).SetTag(fault.TagBadRequest).SetKey(fault.KeyOtp)
		}
		_ = s.otpRepo.DeleteByEmail(ctx, email)
	}
	otp := s.otpProvider.GenerateOTP()
	if err := s.otpRepo.Create(ctx, &domain.Otp{
		Email: email,
		Otp:   otp,
	}); err != nil {
		return "", err
	}
	if err := s.mailService.SendEmailOAuth2("Xác thực OTP", email, domain.Otp{
		Email: email,
		Otp:   otp,
	}, "otp.txt"); err != nil {
		_ = s.otpRepo.DeleteByEmail(ctx, email)
		return "", err
	}
	return otp, nil
}

func (s *AuthService) ResetOTP(ctx context.Context, email string) (string, error) {
	caller := "AuthService.ResetOTP"
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	otpDb, err := s.otpRepo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}
	if otpDb != nil {
		if otpDb.CreatedAt.After(time.Now().Add(-10 * time.Minute)) {
			err := errors.New("otp already exists")
			return "", fault.Wrapf(err, "[%v] otp already exists", caller).SetTag(fault.TagBadRequest).SetKey(fault.KeyOtp)
		}
		_ = s.otpRepo.DeleteByEmail(ctx, email)
	}
	otp := s.otpProvider.GenerateOTP()
	if err := s.otpRepo.Create(ctx, &domain.Otp{
		Email: email,
		Otp:   otp,
	}); err != nil {
		return "", err
	}
	if err := s.mailService.SendEmailOAuth2("Xác thực OTP", email, domain.Otp{
		Email: email,
		Otp:   otp,
	}, "reset-otp.txt"); err != nil {
		_ = s.otpRepo.DeleteByEmail(ctx, email)
		return "", err
	}
	return otp, nil
}

func (s *AuthService) SignUp(ctx context.Context, req request.SignUpRequest) (*dto.TokenDto, *domain.User, error) {
	caller := "AuthService.SignUp"

	isValidOtp, err := s.verifyOTP(ctx, req.Otp, req.Email)
	if err != nil {
		return nil, nil, err
	}
	if !isValidOtp {
		err := fmt.Errorf("otp is invalid")
		return nil, nil, fault.Wrap(err).SetTag(fault.TagUnAuthorize).SetKey(fault.KeyOtp)
	}
	_ = s.otpRepo.DeleteByEmail(ctx, req.Email)
	hasedPw, err := s.hashProvider.Hash(req.Password)
	if err != nil {
		return nil, nil, fault.Wrapf(err, "[%v] failed to hash password", caller).SetTag(fault.TagInternalServer).SetKey(fault.KeyAuth)
	}
	user, err := s.userRepo.Create(ctx, &domain.User{
		Name:         req.Name,
		Email:        req.Email,
		Phone:        "",
		PasswordHash: hasedPw,
		AvatarUrl:    constant.DefaultAvatarUrl,
		Role:         constant.DefaultRole,
	})
	if err != nil {
		return nil, nil, err
	}
	jwtToken, err := s.generateToken(ctx, user)
	if err != nil {
		return nil, nil, err
	}
	return jwtToken, user, nil
}

func (s *AuthService) Login(ctx context.Context, req request.LoginRequest) (*dto.TokenDto, *domain.User, error) {
	caller := "AuthService.Login"

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, nil, err
	}
	err = s.hashProvider.ComparePassword(req.Password, user.PasswordHash)
	if err != nil {
		return nil, nil, fault.Wrapf(err, "[%v] wrong password", caller).SetTag(fault.TagUnAuthorize).SetKey(fault.KeyAuth)
	}
	jwtToken, err := s.generateToken(ctx, user)
	if err != nil {
		return nil, nil, err
	}
	return jwtToken, user, nil
}

func (s *AuthService) LoginAdmin(ctx context.Context, req request.LoginRequest) (*dto.TokenDto, *domain.User, error) {
	caller := "AuthService.LoginAdmin"

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, nil, err
	}

	// Check if user is admin
	if user.Role != constant.AdminRole && user.Role != constant.StaffRole {
		err := fmt.Errorf("user is not an admin")
		return nil, nil, fault.Wrapf(err, "[%v] unauthorized access", caller).SetTag(fault.TagForbidden).SetKey(fault.KeyAuth)
	}

	err = s.hashProvider.ComparePassword(req.Password, user.PasswordHash)
	if err != nil {
		return nil, nil, fault.Wrapf(err, "[%v] wrong password", caller).SetTag(fault.TagUnAuthorize).SetKey(fault.KeyAuth)
	}

	jwtToken, err := s.generateToken(ctx, user)
	if err != nil {
		return nil, nil, err
	}
	return jwtToken, user, nil
}

func (s *AuthService) generateToken(ctx context.Context, user *domain.User) (*dto.TokenDto, error) {
	caller := "AuthService.generateToken"

	payload := crypto.Payload{
		Id:       user.ID.String(),
		Username: user.Email,
	}
	exA, _ := strconv.ParseInt(os.Getenv("JWT_EXPIRE_ACCESS"), 10, 64)
	accessToken, err := s.jwtProvider.Generate(os.Getenv("JWT_ACCESS_SECRET"), payload, exA)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to gen token", caller).SetTag(fault.TagInternalServer).SetKey(fault.KeyAuth)
	}

	exR, _ := strconv.ParseInt(os.Getenv("JWT_EXPIRE_REFRESH"), 10, 64)
	refreshToken, err := s.jwtProvider.Generate(os.Getenv("JWT_REFRESH_SECRET"), payload, exR)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to gen token", caller).SetTag(fault.TagInternalServer).SetKey(fault.KeyAuth)
	}

	return &dto.TokenDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) verifyOTP(ctx context.Context, otp string, email string) (bool, error) {
	caller := "AuthService.verifyOTP"

	otpDb, err := s.otpRepo.GetByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	if otpDb.Otp != otp {
		err := fmt.Errorf("otp validation failed")
		return false, fault.Wrapf(err, "[%v] wrong Otp ", caller).SetTag(fault.TagUnAuthorize).SetKey(fault.KeyOtp)
	}
	return true, nil
}

func (s *AuthService) VerifyToken(ctx context.Context, token string) (*domain.User, error) {
	caller := "AuthUserUseCase.VerifyToken"

	payload, err := s.jwtProvider.Verify(os.Getenv("JWT_ACCESS_SECRET"), token)
	if errors.Is(err, crypto.ErrTokenExpired) {
		return nil, fault.Wrapf(err, "[%v] token expired", caller).SetTag(fault.TagUnAuthorize).SetKey(fault.KeyAuth)

	}
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to verify token", caller).SetTag(fault.TagUnAuthorize).SetKey(fault.KeyAuth)
	}
	username := payload.Username
	user, err := s.userRepo.GetByEmail(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, id uuid.UUID, req request.ChangePasswordRequest) error {
	caller := "AuthService.ChangePassword"

	user, err := s.userRepo.GetById(ctx, id)
	if err != nil {
		return err
	}
	err = s.hashProvider.ComparePassword(req.OldPassword, user.PasswordHash)
	if err != nil {
		return fault.Wrapf(err, "[%v] old password is wrong", caller).SetTag(fault.TagUnAuthorize).SetKey(fault.KeyAuth)
	}
	hashedPw, err := s.hashProvider.Hash(req.NewPassword)
	if err != nil {
		return fault.Wrapf(err, "[%v] failed to hash password", caller).SetTag(fault.TagInternalServer).SetKey(fault.KeyAuth)
	}
	user.PasswordHash = hashedPw
	user, err = s.userRepo.Update(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, req *request.ResetPasswordRequest) error {
	caller := "AuthService.resetPassword"

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	isValidOtp, err := s.verifyOTP(ctx, req.Otp, req.Email)
	if err != nil {
		return err
	}
	if !isValidOtp {
		err := fmt.Errorf("otp is invalid")
		return fault.Wrap(err).SetTag(fault.TagUnAuthorize).SetKey(fault.KeyOtp)
	}
	_ = s.otpRepo.DeleteByEmail(ctx, req.Email)
	hasedPw, err := s.hashProvider.Hash(req.NewPassword)
	if err != nil {
		return fault.Wrapf(err, "[%v] failed to hash password", caller).SetTag(fault.TagInternalServer).SetKey(fault.KeyAuth)
	}

	user.PasswordHash = hasedPw
	user, err = s.userRepo.Update(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
