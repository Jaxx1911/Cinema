package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Discount struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Code       string    `gorm:"type:varchar(50);unique;not null"`
	Percentage float64   `gorm:"not null"`
	StartDate  time.Time `gorm:"not null"`
	EndDate    time.Time `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	Orders []Order `gorm:"foreignKey:DiscountID"`
}

type DiscountRepository interface {
	GetDiscount(ctx context.Context, id uuid.UUID) (*Discount, error)
	GetListDiscount(ctx context.Context) ([]Discount, error)
	GetDiscountByCode(ctx context.Context, code string) (*Discount, error)
	CreateDiscount(ctx context.Context, discount Discount) (*Discount, error)
	UpdateDiscount(ctx context.Context, discount Discount) (*Discount, error)
}

func (*Discount) TableName() string {
	return "discount"
}
