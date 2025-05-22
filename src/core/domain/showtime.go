package domain

import (
	"TTCS/src/present/httpui/request"
	"context"
	"time"

	"github.com/google/uuid"
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

	Movie   Movie    `gorm:"foreignKey:MovieID"`
	Room    Room     `gorm:"foreignKey:RoomID"`
	Tickets []Ticket `gorm:"foreignKey:ShowtimeID"`
	Order   []Order  `gorm:"foreignKey:ShowtimeID"`
}

func (*Showtime) TableName() string {
	return "showtime"
}

type ShowtimeRepo interface {
	Create(ctx context.Context, showtime *Showtime) (*Showtime, error)
	FindConflictByRoomId(ctx context.Context, roomId uuid.UUID, startTime, endTime time.Time) ([]Showtime, error)
	GetListByFilter(ctx context.Context, movieId uuid.UUID, cinemaId uuid.UUID, day time.Time) ([]*Showtime, error)
	GetListByCinemaFilter(ctx context.Context, id uuid.UUID, day time.Time) ([]*Showtime, error)
	GetListByRoomFilter(ctx context.Context, id uuid.UUID, day time.Time) ([]*Showtime, error)
	GetById(ctx context.Context, id uuid.UUID) (*Showtime, error)
	GetList(ctx context.Context, page request.GetListShowtime) ([]*Showtime, int64, error)
	Update(ctx context.Context, id uuid.UUID, showtime *Showtime) (*Showtime, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
