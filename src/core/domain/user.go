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
	Phone        string         `gorm:"type:varchar(20);uniqueIndex"`
	PasswordHash string         `gorm:"type:text;not null"`
	Role         string         `gorm:"type:varchar(50);default:customer;not null"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type UserRepo interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetList(ctx context.Context) ([]*User, error)
	GetById(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
}

func (User) TableName() string {
	return "users"
}

type Otp struct {
	Email string `gorm:"primaryKey"`
	Otp   string
}

type OtpRepo interface {
	Create(ctx context.Context, otp *Otp) error
	GetByEmail(ctx context.Context, email string) (*Otp, error)
	DeleteByEmail(ctx context.Context, email string) error
}

func (Otp) TableName() string {
	return "otp"
}
