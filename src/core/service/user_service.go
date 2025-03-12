package service

import (
	"TTCS/src/core/domain"
	"TTCS/src/present/httpui/request"
	"context"
)

type UserService struct {
	userRepo domain.UserRepo
}

func NewUserService(userRepo domain.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (u *UserService) Update(ctx context.Context, id string, req *request.UserInfo) (*domain.User, error) {
	user, err := u.userRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	user, err = u.userRepo.Update(ctx, u.buildModelUser(req, user))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) buildModelUser(req *request.UserInfo, user *domain.User) *domain.User {
	if req.Email == "" {
		user.Email = req.Email
	}
	if req.Name == "" {
		user.Name = req.Name
	}
	if req.Phone == "" {
		user.Phone = req.Phone
	}
	return user
}
