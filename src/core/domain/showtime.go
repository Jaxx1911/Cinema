package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Showtime struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	MovieID   uuid.UUID `gorm:"type:uuid;not null"`
	RoomID    uuid.UUID `gorm:"type:uuid;not null"`
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`
	Price     float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Movie  Movie    `gorm:"foreignKey:MovieID"`
	Room   Room     `gorm:"foreignKey:RoomID"`
	Ticket []Ticket `gorm:"foreignKey:ShowtimeID"`
}

func (*Showtime) TableName() string {
	return "showtime"
}

type ShowtimeRepo interface {
	Create(ctx context.Context, showtime *Showtime) (*Showtime, error)
	FindConflictByRoomId(ctx context.Context, roomId uuid.UUID, startTime, endTime time.Time) ([]Showtime, error)
	GetListByFilter(ctx context.Context, movieId string, cinemaId string, day time.Time) ([]*Showtime, error)
}
