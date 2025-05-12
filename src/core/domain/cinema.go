package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Cinema struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name         string    `gorm:"type:varchar(255);not null"`
	Address      string    `gorm:"type:text;not null"`
	Phone        string    `gorm:"type:varchar(20)"`
	OpeningHours string    `gorm:"type:text"`
	IsActive     bool      `gorm:"type:boolean;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`

	Rooms []Room `gorm:"foreignKey:CinemaID"`
}

type CinemaRepo interface {
	Create(ctx context.Context, cinema *Cinema) (*Cinema, error)
	GetList(ctx context.Context) ([]*Cinema, error)
	GetListByCity(ctx context.Context, city string) ([]*Cinema, error)
	GetWithRoomsByCity(ctx context.Context, city string) ([]*Cinema, error)
	GetDetail(ctx context.Context, id uuid.UUID) (*Cinema, error)
	Update(ctx context.Context, cinema *Cinema) (*Cinema, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*Cinema, error)
}

func (*Cinema) TableName() string {
	return "cinema"
}
