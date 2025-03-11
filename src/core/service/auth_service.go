package service

import (
	"TTCS/src/common/configs"
	"TTCS/src/common/constant"
	"TTCS/src/common/crypto"
	"TTCS/src/common/fault"
	"TTCS/src/core/domain"
	"TTCS/src/core/dto"
	"TTCS/src/present/httpui/request"
	"context"
	"fmt"
)

type AuthService struct {
	userRepo     domain.UserRepo
	otpRepo      domain.OtpRepo
	hashProvider crypto.HashProvider
	otpProvider  crypto.OTPProvider
}

func NewAuthService(userRepo domain.UserRepo, otpRepo domain.OtpRepo, hashProvider crypto.HashProvider, otpProvider crypto.OTPProvider) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		otpRepo:      otpRepo,
		hashProvider: hashProvider,
		otpProvider:  otpProvider,
	}
}

func (s *AuthService) GenOTP(ctx context.Context, email string) (string, error) {
	_ = "AuthService.GenOTP"
	otp := s.otpProvider.GenerateCode()
	if err := s.otpRepo.Create(ctx, &domain.Otp{
		Email: email,
		Otp:   otp,
	}); err != nil {
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
		Name:         constant.DefaultUserName,
		Email:        req.Email,
		Phone:        "",
		PasswordHash: hasedPw,
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

	accessToken, err := crypto.Generate(configs.GetConfig().Jwt.AccessSecret, payload, configs.GetConfig().Jwt.ExpireAccess)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to gen token", caller).SetTag(fault.TagInternalServer)
	}
	refreshToken, err := crypto.Generate(configs.GetConfig().Jwt.RefreshSecret, payload, configs.GetConfig().Jwt.ExpireRefresh)
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
