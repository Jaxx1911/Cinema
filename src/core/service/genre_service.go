package service

import (
	"TTCS/src/core/domain"
	"context"
)

type GenreService struct {
	genreRepo domain.GenreRepo
}

func NewGenreService(genreRepo domain.GenreRepo) *GenreService {
	return &GenreService{genreRepo: genreRepo}
}

func (g *GenreService) GetGenres(ctx context.Context) ([]*domain.Genre, error) {
	genres, err := g.genreRepo.GetList(ctx)
	if err != nil {
		return nil, err
	}
	return genres, err
}
