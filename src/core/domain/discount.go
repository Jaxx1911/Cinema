package domain

import "time"

type Discount struct {
	ID         uint      `gorm:"primaryKey"`
	Code       string    `gorm:"type:varchar(50);unique;not null"`
	Percentage float64   `gorm:"not null"`
	StartDate  time.Time `gorm:"not null"`
	EndDate    time.Time `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

type DiscountRepository interface{}

func (*Discount) TableName() string {
	return "discount"
}
