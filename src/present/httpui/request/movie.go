package request

import (
	"github.com/google/uuid"
	"mime/multipart"
)

type CreateMovieRequest struct {
	Title            string                `form:"title" binding:"required"`
	Duration         int                   `form:"duration" binding:"required"`
	PosterImage      *multipart.FileHeader `form:"poster_image" binding:"required"`
	LargePosterImage *multipart.FileHeader `form:"large_poster_image" binding:"required"`
	Director         string                `form:"director" binding:"required"`
	Caster           string                `form:"caster" binding:"required"`
	Description      string                `form:"description" binding:"required"`
	ReleaseDate      string                `form:"release_date" binding:"required"`
	TrailerURL       string                `form:"trailer_url" binding:"required"`
	Status           string                `form:"status" binding:"required"`
	Genres           []uuid.UUID           `form:"genres" binding:"required"`
	Tag              string                `form:"tag" binding:"required"`
}

type UpdateMovieRequest struct {
	Id          uuid.UUID   `form:"id,omitempty"`
	Title       string      `form:"title" binding:"required"`
	Duration    int         `form:"duration" binding:"required"`
	Director    string      `form:"director" binding:"required"`
	Caster      string      `form:"caster" binding:"required"`
	Description string      `form:"description" binding:"required"`
	ReleaseDate string      `form:"release_date" binding:"required"`
	TrailerURL  string      `form:"trailer_url" binding:"required"`
	Status      string      `form:"status" binding:"required"`
	Genres      []uuid.UUID `form:"genres" binding:"required"`
}
