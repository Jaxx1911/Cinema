package response

import (
	"TTCS/src/core/domain"
	"github.com/google/uuid"
	"time"
)

type MovieDetailResponse struct {
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
	Tag            string    `json:"tag"`
}

type MovieOfListResponse struct {
	ID             uuid.UUID `json:"id"`
	Title          string    `json:"title"`
	Duration       int       `json:"duration"`
	PosterURL      string    `json:"poster_url"`
	LargePosterURL string    `json:"large_poster_url"`
	Description    string    `json:"description"`
	ReleaseDate    string    `json:"release_date"`
	TrailerURL     string    `json:"trailer_url"`
	Genres         []string  `json:"genres"`
	Status         string    `json:"status"`
	Tag            string    `json:"tag"`
}

func ToMovieDetailResponse(movie *domain.Movie) *MovieDetailResponse {
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
	return &MovieDetailResponse{
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
		Tag:            movie.Tag,
	}
}

func ToMovieResponse(movie *domain.Movie) *MovieOfListResponse {
	genres := make([]string, 0)
	for _, genre := range movie.Genres {
		genres = append(genres, genre.Name)
	}
	return &MovieOfListResponse{
		ID:             movie.ID,
		Title:          movie.Title,
		Duration:       movie.Duration,
		PosterURL:      movie.PosterURL,
		LargePosterURL: movie.LargePosterURL,
		Description:    movie.Description,
		ReleaseDate:    movie.ReleaseDate.Format("2006-01-02"),
		TrailerURL:     movie.TrailerURL,
		Genres:         genres,
		Status:         movie.Status,
		Tag:            movie.Tag,
	}
}

func ToListMoviesResponse(movies []*domain.Movie) []*MovieOfListResponse {
	var movieResponses []*MovieOfListResponse

	for _, movie := range movies {
		genres := make([]string, 0)
		for _, genre := range movie.Genres {
			genres = append(genres, genre.Name)
		}
		var movieResponse = &MovieOfListResponse{
			ID:             movie.ID,
			Title:          movie.Title,
			Duration:       movie.Duration,
			PosterURL:      movie.PosterURL,
			LargePosterURL: movie.LargePosterURL,
			Description:    movie.Description,
			ReleaseDate:    movie.ReleaseDate.Format("2006-01-02"),
			TrailerURL:     movie.TrailerURL,
			Genres:         genres,
			Status:         movie.Status,
			Tag:            movie.Tag,
		}
		movieResponses = append(movieResponses, movieResponse)
	}
	return movieResponses
}
