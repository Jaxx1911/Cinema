package response

import (
	"TTCS/src/core/domain"
	"github.com/google/uuid"
)

type Genre struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func ToGenreResponse(genre domain.Genre) Genre {
	return Genre{
		ID:   genre.ID,
		Name: genre.Name,
	}
}

func ToGenresResponse(genres []*domain.Genre) []Genre {
	var list []Genre
	for _, g := range genres {
		list = append(list, ToGenreResponse(*g))
	}
	return list
}
