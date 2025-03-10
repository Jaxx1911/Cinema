package domain

import (
	"github.com/google/uuid"
	"time"
)

type Cinema struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Address   string    `gorm:"type:text;not null"`
	Phone     string    `gorm:"type:varchar(20)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Rooms []Room `gorm:"foreignKey:CinemaID"`
}

type CinemaRepository interface {
}

func (*Cinema) TableName() string {
	return "cinema"
}
