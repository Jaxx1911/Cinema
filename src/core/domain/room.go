package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Room struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CinemaID    uuid.UUID `gorm:"type:uuid;not null"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Capacity    int       `gorm:"not null"`
	Type        string    `gorm:"type:varchar(50);not null"` // 2D, 3D, IMAX
	RowCount    int       `gorm:"not null"`
	ColumnCount int       `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	Showtimes []Showtime `gorm:"foreignKey:RoomID"`
	Seats     []Seat     `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE"`
}

func (*Room) TableName() string {
	return "room"
}

type RoomRepo interface {
	Create(ctx context.Context, room *Room) error
	GetById(ctx context.Context, roomID string) (*Room, error)
}

type Seat struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RoomID     uuid.UUID `gorm:"type:uuid;not null"`
	RowNumber  string    `gorm:"type:varchar(5);not null"` // Hàng: A, B, C...
	SeatNumber int       `gorm:"not null"`                 // Số ghế trong hàng
	Type       string    `gorm:"type:varchar(50);not null;check:type IN ('standard', 'VIP', 'couple', 'disabled')"`

	Room Room `gorm:"foreignKey:RoomID"`
}

func (*Seat) TableName() string {
	return "seat"
}
