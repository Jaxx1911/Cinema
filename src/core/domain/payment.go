package domain

import "time"

type Payment struct {
	ID          uint      `gorm:"primaryKey"`
	OrderID     uint      `gorm:"not null;uniqueIndex"`      // Một đơn hàng có một thanh toán
	Amount      float64   `gorm:"not null"`                  // Tổng số tiền thanh toán
	Status      string    `gorm:"type:varchar(20);not null"` // pending, completed, failed
	PaymentTime time.Time `gorm:"autoCreateTime"`            // Thời gian thanh toán

	Order Order `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

func (*Payment) TableName() string {
	return "payment"
}
