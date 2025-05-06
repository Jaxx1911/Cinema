package domain

import (
	"TTCS/src/present/httpui/request"
	"context"
	"github.com/google/uuid"
	"time"
)

type Movie struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title          string    `gorm:"type:varchar(255);not null"`
	Duration       int       `gorm:"not null"`
	PosterURL      string    `gorm:"type:text"`
	LargePosterURL string    `gorm:"type:text"`
	Director       string    `gorm:"type:varchar(255)"`
	Caster         string    `gorm:"type:text"`
	Description    string    `gorm:"type:text"`
	ReleaseDate    time.Time `gorm:"not null"`
	TrailerURL     string    `gorm:"type:text"`
	Status         string    `gorm:"not null;default:new"`
	Tag            string    `gorm:"type:varchar(255)"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`

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

type MovieRepo interface {
	GetList(ctx context.Context, page request.Page) ([]*Movie, error)
	GetListByStatus(ctx context.Context, page request.Page, showingStatus string) ([]*Movie, error)
	GetById(ctx context.Context, id uuid.UUID) (*Movie, error)
	GetDetail(ctx context.Context, id uuid.UUID) (*Movie, error)
	Create(ctx context.Context, movie *Movie) (*Movie, error)
	Update(ctx context.Context, movie *Movie) (*Movie, error)
	GetListInDateRange(ctx context.Context, startDate time.Time, endDate time.Time) ([]*Movie, error)
}

type GenreRepo interface {
	GetList(ctx context.Context) ([]*Genre, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Genre, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]Genre, error)
	Create(ctx context.Context, genre *Genre) (*Genre, error)
}
