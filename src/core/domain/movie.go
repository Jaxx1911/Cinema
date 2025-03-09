package domain

import "time"

type Movie struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Duration    int       `gorm:"not null" json:"duration"`
	PosterURL   string    `gorm:"type:text" json:"poster_url"`
	Director    string    `gorm:"type:varchar(255)" json:"director"`
	Caster      string    `gorm:"type:text" json:"cast"`
	Description string    `gorm:"type:text" json:"description"`
	ReleaseDate time.Time `gorm:"not null" json:"release_date"`
	Rating      string    `gorm:"type:varchar(10);not null" json:"rating"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Showtimes []Showtime `gorm:"foreignKey:MovieID" json:"showtimes"`
	Genres    []Genre    `gorm:"many2many:movie_genre" json:"genres"`
}

type Genre struct {
	ID     uint    `gorm:"primaryKey"`
	Name   string  `gorm:"type:varchar(255);unique;not null"`
	Movies []Movie `gorm:"many2many:movie_genre"`
}

func (*Movie) TableName() string {
	return "movie"
}

func (*Genre) TableName() string {
	return "genre"
}
