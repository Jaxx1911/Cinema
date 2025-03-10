package repo

import (
	"TTCS/src/core/domain"
	"context"
	"gorm.io/gorm"
)

type UserRepo struct {
	*BaseRepo
	db *gorm.DB
}

func NewUserRepo(baseRepo *BaseRepo, db *gorm.DB) domain.UserRepo {
	return &UserRepo{
		BaseRepo: baseRepo,
		db:       db,
	}
}

func (u UserRepo) Create(ctx context.Context, user *domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (u UserRepo) GetList(ctx context.Context) ([]*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepo) GetById(ctx context.Context, id uint) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user *domain.User
	if err := u.db.WithContext(ctx).Preload("Auth").Where("email = ?", email).Scan(&user).Error; err != nil {
		return nil, u.returnError(ctx, err)
	}
	return user, nil
}
