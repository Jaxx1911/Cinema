package domain

import (
	"TTCS/src/common"
	"context"
)

type Auth struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	UserID       uint   `gorm:"uniqueIndex;not null" json:"user_id"`
	PasswordHash string `gorm:"type:text;not null" json:"password_hash"`
}

type AuthRepo interface {
	GetByEmail(ctx context.Context, email string) (*Auth, *common.Error)
	GetByID(ctx context.Context, id uint) (*Auth, *common.Error)
	Create(ctx context.Context, auth *Auth) *common.Error
}

func (*Auth) TableName() string {
	return "auth"
}
