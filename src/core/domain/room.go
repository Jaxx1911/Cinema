package domain

import (
	"TTCS/src/present/httpui/request"
	"context"
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CinemaID    uuid.UUID `gorm:"type:uuid;not null"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Capacity    int       `gorm:"not null"`
	Type        string    `gorm:"type:varchar(50);not null"` // 2D, 3D, IMAX
	RowCount    int       `gorm:"not null"`
	ColumnCount int       `gorm:"not null"`
	IsActive    bool      `gorm:"not null;default:true"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	Showtimes []Showtime `gorm:"foreignKey:RoomID"`
	Seats     []Seat     `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE"`
}

func (*Room) TableName() string {
	return "room"
}

type RoomRepo interface {
	Create(ctx context.Context, room *Room) (*Room, error)
	GetById(ctx context.Context, roomID uuid.UUID) (*Room, error)
	GetList(ctx context.Context, page request.GetListRoom) ([]*Room, int64, error)
	GetListByCinemaId(ctx context.Context, cinemaId uuid.UUID) ([]Room, error)
	Deactivate(ctx context.Context, id uuid.UUID, isActive bool) error
	Update(ctx context.Context, room *Room) (*Room, error)
	Delete(ctx context.Context, id uuid.UUID) error
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

type SeatRepo interface {
	Create(ctx context.Context, seat *Seat) error
	GetById(ctx context.Context, seatID uuid.UUID) (*Seat, error)
	GetByRoomID(ctx context.Context, roomID uuid.UUID) ([]Seat, error)
	UpdateSeat(ctx context.Context, seat *Seat) error
}
