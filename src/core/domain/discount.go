package domain

import (
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

type DiscountRepository interface{}

func (*Discount) TableName() string {
	return "discount"
}
