package service

import (
	"TTCS/src/common/crypto"
	"TTCS/src/common/fault"
	"TTCS/src/common/log"
	"TTCS/src/common/mail"
	"TTCS/src/core/domain"
	"TTCS/src/infra/upload"
	"TTCS/src/present/httpui/request"
	"context"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo     domain.UserRepo
	hashProvider crypto.HashProvider
	upload       *upload.UploadService
	mailService  *mail.GmailService
}

func NewUserService(userRepo domain.UserRepo, hashProvider crypto.HashProvider, upload *upload.UploadService, mailService *mail.GmailService) *UserService {
	return &UserService{
		userRepo:     userRepo,
		hashProvider: hashProvider,
		upload:       upload,
		mailService:  mailService,
	}
}

func (u *UserService) Update(ctx context.Context, id uuid.UUID, req *request.UserInfo) (*domain.User, error) {
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

func (u *UserService) Create(ctx context.Context, req *request.UserInfo) (*domain.User, error) {
	caller := "UserService.Create"
	user := u.buildModelUser(req, &domain.User{})

	password := strings.TrimSuffix(req.Email, "@gmail.com")
	hashedPw, err := u.hashProvider.Hash(password)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to hash email", caller).SetKey(fault.KeyUser)
	}
	user.PasswordHash = hashedPw
	user, err = u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Gửi email thông báo tài khoản được tạo nếu role là customer
	if req.Role == "customer" {
		emailData := struct {
			Name     string
			Email    string
			Password string
		}{
			Name:     user.Name,
			Email:    user.Email,
			Password: password,
		}

		err = u.mailService.SendEmailOAuth2("Tài khoản được tạo thành công", user.Email, emailData, "account-created.txt")
		if err != nil {
			log.Error(ctx, "[%v] failed to send account creation email %+v", caller, err)
			// Không return error vì tạo user đã thành công, chỉ log lỗi gửi email
		}
	}

	return user, nil
}

func (u *UserService) buildModelUser(req *request.UserInfo, user *domain.User) *domain.User {
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	return user
}

func (u *UserService) GetList(ctx context.Context, page request.GetListUser) ([]*domain.User, int64, error) {
	_ = "UserService.GetList"
	users, total, err := u.userRepo.GetList(ctx, page)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (u *UserService) GetById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	_ = "UserService.GetById"
	user, err := u.userRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) GetPayments(ctx context.Context, id uuid.UUID) ([]domain.Payment, error) {
	_ = "UserService.GetPayment"
	payments, err := u.userRepo.GetPaymentsById(ctx, id)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (u *UserService) GetOrders(ctx context.Context, id uuid.UUID) ([]domain.Order, error) {
	_ = "UserService.GetPayment"
	orders, err := u.userRepo.GetOrdersById(ctx, id)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (u *UserService) ChangeAvatar(ctx context.Context, file *multipart.FileHeader, user *domain.User) (*domain.User, error) {
	caller := "UserService.ChangeAvatar"
	url, err := u.upload.UploadFile(ctx, file)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to upload file %v", caller, err)
	}
	user.AvatarUrl = url
	user, err = u.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	caller := "UserService.Delete"

	// First check if user exists
	user, err := u.userRepo.GetById(ctx, id)
	if err != nil {
		return fault.Wrapf(err, "[%v] failed to get user", caller)
	}

	// Use GORM's soft delete
	if err := u.userRepo.Delete(ctx, user); err != nil {
		return fault.Wrapf(err, "[%v] failed to delete user", caller)
	}

	return nil
}
