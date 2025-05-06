package service

import (
	"TTCS/src/core/domain"
	"TTCS/src/present/httpui/request"
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

func (c *CinemaService) Create(ctx context.Context, req request.CreateCinemaRequest) (*domain.Cinema, error) {
	cinema, err := c.cinemaRepo.Create(ctx, &domain.Cinema{
		Name:         req.Name,
		Address:      req.Address,
		Phone:        req.Phone,
		OpeningHours: req.OpeningHours,
	})
	if err != nil {
		return nil, err
	}
	return cinema, nil
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

func (c *CinemaService) Update(ctx context.Context, req request.UpdateCinemaRequest) (*domain.Cinema, error) {
	cinema, err := c.cinemaRepo.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	cinema.Name = req.Name
	cinema.Address = req.Address
	cinema.Phone = req.Phone
	cinema.OpeningHours = req.OpeningHours
	cinema, err = c.cinemaRepo.Update(ctx, cinema)
	if err != nil {
		return nil, err
	}
	return cinema, nil
}

func (c *CinemaService) Delete(ctx context.Context, id uuid.UUID) error {
	err := c.cinemaRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
