package domain

import "time"

type Order struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"not null"`
	DiscountID *uint     `gorm:"default:null"`
	Status     string    `gorm:"type:varchar(20);not null;default:'pending'"`
	TotalPrice float64   `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	Tickets     []Ticket     `gorm:"foreignKey:OrderID"`
	OrderCombos []OrderCombo `gorm:"foreignKey:OrderID"`
	Discount    *Discount    `gorm:"foreignKey:DiscountID"`
}

type OrderRepository interface{}

func (*Order) TableName() string {
	return "orders"
}

type OrderCombo struct {
	ID         uint    `gorm:"primaryKey"`
	OrderID    uint    `gorm:"not null"`
	ComboID    uint    `gorm:"not null"`
	Quantity   int     `gorm:"not null"`
	TotalPrice float64 `gorm:"not null"`
}

func (*OrderCombo) TableName() string {
	return "order_combo"
}
