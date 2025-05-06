package service

import (
	"TTCS/src/core/domain"
	"context"
	"github.com/google/uuid"
)

type CinemaService struct {
	cinemaRepo domain.CinemaRepo
}

func NewCinemaService(cinemaRepo domain.CinemaRepo) *CinemaService {
	return &CinemaService{
		cinemaRepo: cinemaRepo,
	}
}

func (c *CinemaService) Create(ctx context.Context, cinema domain.Cinema) error {
	return nil
}

func (c *CinemaService) GetList(ctx context.Context) ([]*domain.Cinema, error) {
	cinemas, err := c.cinemaRepo.GetList(ctx)
	if err != nil {
		return nil, err
	}
	return cinemas, nil
}

func (c *CinemaService) GetListByCity(ctx context.Context, city string) ([]*domain.Cinema, error) {
	cinemas, err := c.cinemaRepo.GetListByCity(ctx, city)
	if err != nil {
		return nil, err
	}
	return cinemas, nil
}

func (c *CinemaService) GetWithRoomsByCity(ctx context.Context, city string) ([]*domain.Cinema, error) {
	cinemas, err := c.cinemaRepo.GetWithRoomsByCity(ctx, city)
	if err != nil {
		return nil, err
	}
	return cinemas, nil
}

func (c *CinemaService) GetDetail(ctx context.Context, id uuid.UUID) (*domain.Cinema, error) {
	cinema, err := c.cinemaRepo.GetDetail(ctx, id)
	if err != nil {
		return nil, err
	}
	return cinema, nil
}
