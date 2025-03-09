package service

import (
	"TTCS/src/common"
	"TTCS/src/common/configs"
	"TTCS/src/common/jwt"
	"TTCS/src/core/domain"
	"TTCS/src/infra/repo"
	"TTCS/src/present/httpui/request"
	"TTCS/src/present/httpui/response"
	"context"
	"fmt"
	"strconv"
)

type AuthService struct {
	authRepo     *repo.AuthRepo
	userRepo     *repo.UserRepo
	hashProvider jwt.HashProvider
}

func (s *AuthService) Login(ctx context.Context, req request.LoginRequest) (*response.Token, *domain.User, *common.Error) {
	caller := "AuthService.Login"

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, nil, err
	}
	ierr := s.hashProvider.ComparePassword(req.Password, user.Auth.PasswordHash)
	if ierr != nil {
		return nil, nil, common.ErrBadRequest(ctx).SetDetail(fmt.Sprintf("[%v] wrong password", caller))
	}
	jwtToken, ierr := s.generateToken(ctx, user)
	if ierr != nil {
		return nil, nil, common.ErrBadRequest(ctx).SetDetail(fmt.Sprintf("[%v] failed to gen token", caller))
	}
	return jwtToken, user, nil
}

func NewAuthService(authRepo *repo.AuthRepo, userRepo *repo.UserRepo, hashProvider jwt.HashProvider) *AuthService {
	return &AuthService{
		authRepo:     authRepo,
		userRepo:     userRepo,
		hashProvider: hashProvider,
	}
}

func (s *AuthService) generateToken(ctx context.Context, user *domain.User) (*response.Token, error) {
	caller := "AuthService.generateToken"

	payload := jwt.Payload{
		Id:       strconv.FormatUint(uint64(user.ID), 10),
		Username: user.Email,
	}

	accessToken, err := jwt.Generate(configs.GetConfig().Jwt.AccessSecret, payload, configs.GetConfig().Jwt.ExpireAccess)
	if err != nil {
		return nil, common.ErrBadRequest(ctx).SetDetail(fmt.Sprintf("[%v] failed to gen token", caller))
	}
	refreshToken, err := jwt.Generate(configs.GetConfig().Jwt.RefreshSecret, payload, configs.GetConfig().Jwt.ExpireRefresh)
	if err != nil {
		return nil, common.ErrBadRequest(ctx).SetDetail(fmt.Sprintf("[%v] failed to gen token", caller))
	}

	return &response.Token{
		AccessToken:      accessToken.Token,
		AccessExpiredAt:  accessToken.Expire,
		RefreshToken:     refreshToken.Token,
		RefreshExpiredAt: refreshToken.Expire,
	}, nil
}
