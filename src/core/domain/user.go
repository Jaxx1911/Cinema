package domain

import (
	"TTCS/src/present/httpui/request"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name         string         `gorm:"type:varchar(255)"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	Phone        string         `gorm:"type:varchar(20)"`
	PasswordHash string         `gorm:"type:text;not null"`
	AvatarUrl    string         `gorm:"type:varchar(255)"`
	Role         string         `gorm:"type:varchar(50);default:customer;not null"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	Orders   []Order   `gorm:"foreignKey:UserID"`
	Payments []Payment `gorm:"foreignKey:UserID"`
}

type UserRepo interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetList(ctx context.Context, page request.Page) ([]*User, error)
	GetById(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	GetPaymentsById(ctx context.Context, id uuid.UUID) ([]Payment, error)
	GetOrdersById(ctx context.Context, id uuid.UUID) ([]Order, error)
}

func (User) TableName() string {
	return "users"
}

type Otp struct {
	Email     string `gorm:"primaryKey"`
	Otp       string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type OtpRepo interface {
	Create(ctx context.Context, otp *Otp) error
	GetByEmail(ctx context.Context, email string) (*Otp, error)
	DeleteByEmail(ctx context.Context, email string) error
}

func (Otp) TableName() string {
	return "otp"
}
