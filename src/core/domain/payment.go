package domain

import (
	"github.com/google/uuid"
	"time"
)

type Payment struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID     uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	Amount      float64   `gorm:"not null"`
	Status      string    `gorm:"type:varchar(20);not null"`
	PaymentTime time.Time `gorm:"autoCreateTime"`

	Order Order `gorm:"foreignKey:OrderID"`
}

func (*Payment) TableName() string {
	return "payment"
}
