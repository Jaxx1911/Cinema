package domain

import (
	"github.com/google/uuid"
	"time"
)

type Showtime struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	MovieID   uuid.UUID `gorm:"type:uuid;not null"`
	RoomID    uuid.UUID `gorm:"type:uuid;not null"`
	StartTime time.Time `gorm:"not null"`
	Price     float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (*Showtime) TableName() string {
	return "showtime"
}
