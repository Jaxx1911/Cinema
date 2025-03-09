package service

import (
	"TTCS/src/core/domain"
)

type UserService struct {
	userRepo domain.UserRepo
}

func NewUserService(userRepo domain.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}
