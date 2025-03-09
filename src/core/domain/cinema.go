package domain

import "time"

type Cinema struct {
	ID        uint      `gorm:"primaryKey"`
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
