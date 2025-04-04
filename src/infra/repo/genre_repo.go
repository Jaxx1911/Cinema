package repo

import (
	"TTCS/src/common/fault"
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

func (g GenreRepo) GetByID(ctx context.Context, id string) (*domain.Genre, error) {
	parse, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	genre := &domain.Genre{}
	if err := g.db.WithContext(ctx).Where("id = ?", parse).First(&genre).Error; err != nil {
		return nil, err
	}
	return genre, nil
}

func (g GenreRepo) GetByIDs(ctx context.Context, ids []string) ([]domain.Genre, error) {
	var uuids []uuid.UUID
	for _, id := range ids {
		parse, err := uuid.Parse(id)
		if err != nil {
			return nil, fault.Wrapf(err, "fail to parse id: %s", id)
		}
		uuids = append(uuids, parse)
	}

	var genres []domain.Genre
	if err := g.db.WithContext(ctx).Where("id in (?)", uuids).Find(&genres).Error; err != nil {
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
