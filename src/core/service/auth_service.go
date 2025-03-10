package service

import (
	"TTCS/src/common/configs"
	"TTCS/src/common/crypto"
	"TTCS/src/common/fault"
	"TTCS/src/core/domain"
	"TTCS/src/present/httpui/request"
	"TTCS/src/present/httpui/response"
	"context"
)

type AuthService struct {
	userRepo     domain.UserRepo
	hashProvider crypto.HashProvider
}

func NewAuthService(userRepo domain.UserRepo, hashProvider crypto.HashProvider) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		hashProvider: hashProvider,
	}
}

func (s *AuthService) Login(ctx context.Context, req request.LoginRequest) (*response.Token, *domain.User, error) {
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
		return nil, nil, fault.Wrapf(err, "[%v] failed to generate token", caller).SetTag(fault.TagBadRequest)
	}
	return jwtToken, user, nil
}

func (s *AuthService) generateToken(ctx context.Context, user *domain.User) (*response.Token, error) {
	caller := "AuthService.generateToken"

	payload := crypto.Payload{
		Id:       user.ID.String(),
		Username: user.Email,
	}

	accessToken, err := crypto.Generate(configs.GetConfig().Jwt.AccessSecret, payload, configs.GetConfig().Jwt.ExpireAccess)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to gen token", caller)
	}
	refreshToken, err := crypto.Generate(configs.GetConfig().Jwt.RefreshSecret, payload, configs.GetConfig().Jwt.ExpireRefresh)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to gen token", caller)
	}

	return &response.Token{
		AccessToken:      accessToken.Token,
		AccessExpiredAt:  accessToken.Expire,
		RefreshToken:     refreshToken.Token,
		RefreshExpiredAt: refreshToken.Expire,
	}, nil
}
