package domain

import "time"

type Room struct {
	ID          uint      `gorm:"primaryKey"`
	CinemaID    uint      `gorm:"not null"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Capacity    int       `gorm:"not null"`
	Type        string    `gorm:"type:varchar(50);not null"` // 2D, 3D, IMAX, 4DX
	RowCount    int       `gorm:"not null"`                  // Số hàng ghế
	ColumnCount int       `gorm:"not null"`                  // Số cột ghế
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	Seats     []Seat     `gorm:"foreignKey:RoomID"`
	Showtimes []Showtime `gorm:"foreignKey:RoomID"`
}

func (*Room) TableName() string {
	return "room"
}

type Seat struct {
	ID         uint   `gorm:"primaryKey"`
	RoomID     uint   `gorm:"not null"`
	RowNumber  string `gorm:"type:char(1);not null"` // A, B, C...
	SeatNumber int    `gorm:"not null"`
	Type       string `gorm:"type:varchar(50);not null"` // Standard, VIP, Couple
}

func (*Seat) TableName() string {
	return "seat"
}
