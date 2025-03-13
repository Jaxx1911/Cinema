package service

import (
	"TTCS/src/common/crypto"
	"TTCS/src/common/fault"
	"TTCS/src/core/domain"
	"TTCS/src/present/httpui/request"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strings"
)

type UserService struct {
	userRepo     domain.UserRepo
	hashProvider crypto.HashProvider
}

func NewUserService(userRepo domain.UserRepo, hashProvider crypto.HashProvider) *UserService {
	return &UserService{
		userRepo:     userRepo,
		hashProvider: hashProvider,
	}
}

func (u *UserService) Update(ctx context.Context, id string, req *request.UserInfo) (*domain.User, error) {
	_ = "UserService.Update"
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

func (u *UserService) Create(ctx *gin.Context, r *request.UserInfo) (*domain.User, error) {
	caller := "UserService.Create"
	user := u.buildModelUser(r, &domain.User{})

	hashedPw, err := u.hashProvider.Hash(strings.TrimSuffix(r.Email, "@gmail.com"))
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to hash email", caller)
	}
	user.PasswordHash = hashedPw
	user, err = u.userRepo.Create(ctx, user)
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
	if req.Role == "" {
		user.Phone = req.Phone
	}
	return user
}

func (u *UserService) GetList(ctx *gin.Context, page request.Page) ([]*domain.User, error) {
	_ = "UserService.GetList"
	users, err := u.userRepo.GetList(ctx, page)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserService) GetById(ctx *gin.Context, id string) (*domain.User, error) {
	_ = "UserService.GetById"
	user, err := u.userRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) GetPayments(ctx *gin.Context, id uuid.UUID) ([]domain.Payment, error) {
	_ = "UserService.GetPayment"
	payments, err := u.userRepo.GetPaymentsById(ctx, id)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (u *UserService) GetOrders(ctx *gin.Context, id uuid.UUID) ([]domain.Order, error) {
	_ = "UserService.GetPayment"
	orders, err := u.userRepo.GetOrdersById(ctx, id)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
