package repo

import (
	"TTCS/src/common"
	"TTCS/src/core/domain"
	"context"
	"gorm.io/gorm"
)

type AuthRepo struct {
	*BaseRepo
	db *gorm.DB
}

func (a *AuthRepo) GetByEmail(ctx context.Context, email string) (*domain.Auth, *common.Error) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthRepo) GetByID(ctx context.Context, id uint) (*domain.Auth, *common.Error) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthRepo) Create(ctx context.Context, auth *domain.Auth) *common.Error {
	//TODO implement me
	panic("implement me")
}

func NewAuthRepo(baseRepo *BaseRepo, db *gorm.DB) domain.AuthRepo {
	return &AuthRepo{
		BaseRepo: baseRepo,
		db:       db,
	}
}
