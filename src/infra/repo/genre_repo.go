package repo

import (
	"TTCS/src/core/domain"
	"context"
	"github.com/google/uuid"
)

type GenreRepo struct {
	*BaseRepo
}

func (g GenreRepo) GetList(ctx context.Context) ([]*domain.Genre, error) {
	var genres []*domain.Genre
	if err := g.db.WithContext(ctx).Find(&genres).Error; err != nil {
		return nil, g.returnError(ctx, err)
	}
	return genres, nil
}

func (g GenreRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Genre, error) {
	genre := &domain.Genre{}
	if err := g.db.WithContext(ctx).Where("id = ?", id).First(&genre).Error; err != nil {
		return nil, err
	}
	return genre, nil
}

func (g GenreRepo) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.Genre, error) {
	var genres []domain.Genre
	if err := g.db.WithContext(ctx).Where("id in (?)", ids).Find(&genres).Error; err != nil {
		return nil, g.returnError(ctx, err)
	}
	return genres, nil

}

func (g GenreRepo) Create(ctx context.Context, genre *domain.Genre) (*domain.Genre, error) {
	if err := g.db.WithContext(ctx).Create(genre).Error; err != nil {
		return nil, g.returnError(ctx, err)
	}
	return genre, nil
}

func NewGenreRepo(baseRepo *BaseRepo) domain.GenreRepo {
	return &GenreRepo{
		BaseRepo: baseRepo,
	}
}
