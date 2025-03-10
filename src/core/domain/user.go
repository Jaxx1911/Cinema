package domain

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name         string         `gorm:"type:varchar(255)"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	Phone        string         `gorm:"type:varchar(20);uniqueIndex;not null"`
	PasswordHash string         `gorm:"type:text;not null"`
	Role         string         `gorm:"type:varchar(50);default:customer;not null"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type UserRepo interface {
	Create(ctx context.Context, user *User) error
	GetList(ctx context.Context) ([]*User, error)
	GetById(ctx context.Context, id uint) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

func (*User) TableName() string {
	return "users"
}
