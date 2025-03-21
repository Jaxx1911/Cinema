package response

import (
	"TTCS/src/core/domain"
	"github.com/google/uuid"
	"time"
)

type MovieResponse struct {
	ID             uuid.UUID `json:"id"`
	Title          string    `json:"title"`
	Duration       int       `json:"duration"`
	PosterURL      string    `json:"poster_url"`
	LargePosterURL string    `json:"large_poster_url"`
	Director       string    `json:"director"`
	Caster         string    `json:"caster"`
	Description    string    `json:"description"`
	ReleaseDate    string    `json:"release_date"`
	TrailerURL     string    `json:"trailer_url"`
	Status         string    `json:"status"`
	Genres         []string  `json:"genres"`
}

func ToMovieResponse(movie *domain.Movie) *MovieResponse {
	var status string
	if movie.ReleaseDate.After(time.Now()) {
		status = "incoming"
	} else {
		status = "new"
	}
	genres := make([]string, 0)
	for _, genre := range movie.Genres {
		genres = append(genres, genre.Name)
	}
	return &MovieResponse{
		ID:             movie.ID,
		Title:          movie.Title,
		Duration:       movie.Duration,
		PosterURL:      movie.PosterURL,
		LargePosterURL: movie.LargePosterURL,
		Director:       movie.Director,
		Caster:         movie.Caster,
		Description:    movie.Description,
		ReleaseDate:    movie.ReleaseDate.Format("2006-01-02"),
		TrailerURL:     movie.TrailerURL,
		Status:         status,
		Genres:         genres,
	}
}

func ToListMoviesResponse(movies []*domain.Movie) []*MovieResponse {
	var movieResponses []*MovieResponse
	for _, movie := range movies {
		var movieResponse = ToMovieResponse(movie)
		movieResponses = append(movieResponses, movieResponse)
	}
	return movieResponses
}

type MovieDetailResponse struct {
	MovieResponse
}
