package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null"`
	ShowtimeID uuid.UUID  `gorm:"type:uuid;not null"`
	DiscountID *uuid.UUID `gorm:"type:uuid"`
	Status     string     `gorm:"type:varchar(20);not null;default:'pending'"`
	TotalPrice float64    `gorm:"not null"`
	CreatedAt  time.Time  `gorm:"autoCreateTime"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime"`

	User        User         `gorm:"foreignKey:UserID"`
	Showtime    Showtime     `gorm:"foreignKey:ShowtimeID"`
	Discount    *Discount    `gorm:"foreignKey:DiscountID"`
	Tickets     []Ticket     `gorm:"foreignKey:OrderID"`
	OrderCombos []OrderCombo `gorm:"foreignKey:OrderID"`
}

type OrderRepo interface {
	Create(ctx context.Context, order *Order) (*Order, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Order, error)
	GetDetailByID(ctx context.Context, id uuid.UUID) (*Order, error)
	Update(ctx context.Context, order *Order) (*Order, error)
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

type OrderComboRepository interface {
	Create(ctx context.Context, orderCombo *OrderCombo) (*OrderCombo, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]OrderCombo, error)
}

func (*OrderCombo) TableName() string {
	return "order_combo"
}
