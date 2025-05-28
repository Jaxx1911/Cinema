package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID            uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID        *uuid.UUID `gorm:"type:uuid"`
	OrderID       *uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	TransactionID string     `gorm:"type:varchar(50);not null;uniqueIndex"`
	Status        string     `gorm:"type:varchar(20);not null"`
	Amount        float64    `gorm:"not null"`

	PaymentTime time.Time `gorm:"autoCreateTime"`

	Order Order `gorm:"foreignKey:OrderID"`
	User  User  `gorm:"foreignKey:UserID"`
}

type PaymentRepo interface {
	Create(ctx context.Context, payment *Payment) (*Payment, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]Payment, error)
	GetByCinemaID(ctx context.Context, cinemaID uuid.UUID) ([]Payment, error)
	GetByCinemaIDAndDateRange(ctx context.Context, cinemaID uuid.UUID, startDate, endDate time.Time) ([]Payment, error)
}

func (*Payment) TableName() string {
	return "payment"
}
