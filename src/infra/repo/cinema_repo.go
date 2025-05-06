package repo

import (
	"TTCS/src/core/domain"
	"context"
	"github.com/google/uuid"
)

type CinemaRepo struct {
	*BaseRepo
}

func NewCinemaRepo(baseRepo *BaseRepo) domain.CinemaRepo {
	return &CinemaRepo{
		BaseRepo: baseRepo,
	}
}

func (c *CinemaRepo) Create(ctx context.Context, cinema *domain.Cinema) (*domain.Cinema, error) {
	if err := c.db.WithContext(ctx).Create(cinema).Error; err != nil {
		return nil, c.returnError(ctx, err)
	}
	return cinema, nil
}

func (c *CinemaRepo) GetList(ctx context.Context) ([]*domain.Cinema, error) {
	var cinema []*domain.Cinema
	if err := c.db.WithContext(ctx).Find(&cinema).Error; err != nil {
		return nil, c.returnError(ctx, err)
	}
	return cinema, nil
}

func (c *CinemaRepo) GetListByCity(ctx context.Context, city string) ([]*domain.Cinema, error) {
	var cinema []*domain.Cinema
	if err := c.db.WithContext(ctx).Where("address LIKE ?", "%"+city+"%").Find(&cinema).Error; err != nil {
		return nil, c.returnError(ctx, err)
	}
	return cinema, nil
}

func (c *CinemaRepo) GetWithRoomsByCity(ctx context.Context, city string) ([]*domain.Cinema, error) {
	var cinema []*domain.Cinema
	if err := c.db.WithContext(ctx).Preload("Rooms").Where("address LIKE ?", "%"+city+"%").Find(&cinema).Error; err != nil {
		return nil, c.returnError(ctx, err)
	}
	return cinema, nil
}

func (c *CinemaRepo) GetDetail(ctx context.Context, id uuid.UUID) (*domain.Cinema, error) {
	var cinema domain.Cinema
	if err := c.db.WithContext(ctx).Preload("Rooms").Where("id = ?", id).Find(&cinema).Error; err != nil {
		return nil, c.returnError(ctx, err)
	}
	return &cinema, nil
}
