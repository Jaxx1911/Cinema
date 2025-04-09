package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Combo struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	BannerUrl   string    `gorm:"type:varchar(255);not null"`
	Price       float64   `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type ComboRepository interface {
	Create(ctx context.Context, combo *Combo) error
	FindAll(ctx context.Context) ([]*Combo, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Combo, error)
	Update(ctx context.Context, combo *Combo) error
	Delete(ctx context.Context, id uuid.UUID) error
}

func (*Combo) TableName() string {
	return "combo"
}
