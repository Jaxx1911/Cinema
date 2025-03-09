package domain

import (
	"TTCS/src/common"
	"context"
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Phone     string    `gorm:"type:varchar(20);uniqueIndex;not null"`
	Role      string    `gorm:"type:varchar(50);default:customer;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Auth   Auth    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Orders []Order `gorm:"foreignKey:UserID"`
}

type UserRepo interface {
	Create(ctx context.Context, user *User) *common.Error
	GetList(ctx context.Context) ([]*User, *common.Error)
	GetById(ctx context.Context, id uint) (*User, *common.Error)
	GetByEmail(ctx context.Context, email string) (*User, *common.Error)
}

func (*User) TableName() string {
	return "users"
}
