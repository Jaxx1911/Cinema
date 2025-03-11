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

func (r *UserRepo) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	err := r.db.Create(user).Error
	if err != nil {
		return nil, r.returnError(ctx, err)
	}
	return user, nil
}

func (r *UserRepo) GetList(ctx context.Context) ([]*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepo) GetById(ctx context.Context, id uint) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, r.returnError(ctx, err)
	}
	return &user, nil
}
