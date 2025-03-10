package domain

import (
	"github.com/google/uuid"
	"time"
)

type Movie struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title       string    `gorm:"type:varchar(255);not null"`
	Duration    int       `gorm:"not null"`
	PosterURL   string    `gorm:"type:text"`
	Director    string    `gorm:"type:varchar(255)"`
	Caster      string    `gorm:"type:text"`
	Description string    `gorm:"type:text"`
	ReleaseDate time.Time `gorm:"not null"`
	Rating      string    `gorm:"type:varchar(10);not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	Showtimes []Showtime `gorm:"foreignKey:MovieID;constraint:OnDelete:CASCADE"`
	Genres    []Genre    `gorm:"many2many:movie_genre"`
}

type Genre struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"type:varchar(255);unique;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Movies []Movie `gorm:"many2many:movie_genre"`
}

func (*Movie) TableName() string {
	return "movie"
}

func (*Genre) TableName() string {
	return "genre"
}
