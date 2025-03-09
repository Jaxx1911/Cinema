package domain

import (
	"time"
)

type Combo struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	Price       float64   `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	OrderCombos []OrderCombo `gorm:"foreignKey:ComboID"`
}

type ComboRepository interface {
}

func (*Combo) TableName() string {
	return "combo"
}
