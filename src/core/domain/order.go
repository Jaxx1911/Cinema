package domain

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null"`
	DiscountID *uuid.UUID `gorm:"type:uuid"`
	Status     string     `gorm:"type:varchar(20);not null;default:'pending'"`
	TotalPrice float64    `gorm:"not null"`
	CreatedAt  time.Time  `gorm:"autoCreateTime"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime"`

	User     User      `gorm:"foreignKey:UserID"`
	Discount *Discount `gorm:"foreignKey:DiscountID"`
	Tickets  []Ticket  `gorm:"foreignKey:OrderID"`
}

type OrderRepository interface {
	GetDetail(id uuid.UUID) (*Order, *Movie, *Showtime, error)
}

func (*Order) TableName() string {
	return "orders"
}

type OrderCombo struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID    uuid.UUID `gorm:"type:uuid;not null"`
	ComboID    uuid.UUID `gorm:"type:uuid;not null"`
	Quantity   int       `gorm:"not null"`
	TotalPrice float64   `gorm:"not null"`

	Order Order `gorm:"foreignKey:OrderID"`
	Combo Combo `gorm:"foreignKey:ComboID"`
}

func (*OrderCombo) TableName() string {
	return "order_combo"
}
