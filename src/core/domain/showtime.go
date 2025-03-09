package domain

import "time"

type Showtime struct {
	ID        uint      `gorm:"primaryKey"`
	MovieID   uint      `gorm:"not null"`
	RoomID    uint      `gorm:"not null"`
	StartTime time.Time `gorm:"not null"`
	Price     float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Tickets []Ticket `gorm:"foreignKey:ShowtimeID"`
}

func (*Showtime) TableName() string {
	return "showtime"
}
