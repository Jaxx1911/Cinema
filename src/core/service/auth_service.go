package service

import (
	"TTCS/src/common/configs"
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
	"gorm.io/gorm"
	"time"
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

func (s *AuthService) GenOTP(ctx context.Context, email string) (string, error) {
	caller := "AuthService.GenOTP"
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}
	if user != nil {
		err := errors.New("user already exists")
		return "", fault.Wrapf(err, "[%v] user already exists", caller).SetTag(fault.TagBadRequest)
	}
	otpDb, err := s.otpRepo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}
	if otpDb != nil {
		if otpDb.CreatedAt.After(time.Now().Add(-5 * time.Minute)) {
			err := errors.New("otp already exists")
			return "", fault.Wrapf(err, "[%v] otp already exists", caller).SetTag(fault.TagBadRequest)
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
	if err := s.mailService.SendEmaiOAuth2("Xác thực OTP", email, domain.Otp{
		Email: email,
		Otp:   otp,
	}, "otp.txt"); err != nil {
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
		return nil, nil, fault.Wrap(err).SetTag(fault.TagUnAuthorize)
	}
	hasedPw, err := s.hashProvider.Hash(req.Password)
	if err != nil {
		return nil, nil, fault.Wrapf(err, "[%v] failed to hash password", caller).SetTag(fault.TagInternalServer)
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
		return nil, nil, fault.Wrapf(err, "[%v] wrong password", caller).SetTag(fault.TagUnAuthorize)
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

	accessToken, err := s.jwtProvider.Generate(configs.GetConfig().Jwt.AccessSecret, payload, configs.GetConfig().Jwt.ExpireAccess)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to gen token", caller).SetTag(fault.TagInternalServer)
	}
	refreshToken, err := s.jwtProvider.Generate(configs.GetConfig().Jwt.RefreshSecret, payload, configs.GetConfig().Jwt.ExpireRefresh)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to gen token", caller).SetTag(fault.TagInternalServer)
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
		return false, fault.Wrapf(err, "[%v] wrong Otp ", caller).SetTag(fault.TagUnAuthorize)
	}
	return true, nil
}

func (s *AuthService) VerifyToken(ctx context.Context, token string) (*domain.User, error) {
	caller := "AuthUserUseCase.VerifyToken"

	payload, err := s.jwtProvider.Verify(configs.GetConfig().Jwt.AccessSecret, token)
	if errors.Is(err, crypto.ErrTokenExpired) {
		return nil, fault.Wrapf(err, "[%v] token expired", caller).SetTag(fault.TagUnAuthorize)

	}
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to verify token", caller).SetTag(fault.TagUnAuthorize)
	}
	username := payload.Username
	user, err := s.userRepo.GetByEmail(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, req request.ChangePasswordRequest) error {
	caller := "AuthService.ChangePassword"

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	err = s.hashProvider.ComparePassword(req.OldPassword, user.PasswordHash)
	if err != nil {
		return fault.Wrapf(err, "[%v] old password is wrong", caller).SetTag(fault.TagUnAuthorize)
	}
	hashedPw, err := s.hashProvider.Hash(req.NewPassword)
	if err != nil {
		return fault.Wrapf(err, "[%v] failed to hash password", caller).SetTag(fault.TagInternalServer)
	}
	user.PasswordHash = hashedPw
	user, err = s.userRepo.Update(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
