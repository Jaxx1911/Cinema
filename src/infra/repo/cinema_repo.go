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

func (c *CinemaRepo) Update(ctx context.Context, cinema *domain.Cinema) (*domain.Cinema, error) {
	updates := map[string]interface{}{
		"name":          cinema.Name,
		"address":       cinema.Address,
		"phone":         cinema.Phone,
		"opening_hours": cinema.OpeningHours,
		"is_active":     cinema.IsActive,
	}

	if err := c.db.WithContext(ctx).
		Model(&domain.Cinema{}).
		Where("id = ?", cinema.ID).
		Updates(updates).Error; err != nil {
		return nil, c.returnError(ctx, err)
	}

	return cinema, nil
}

func (c *CinemaRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Cinema, error) {
	cinema := new(domain.Cinema)
	if err := c.db.WithContext(ctx).Where("id = ?", id).First(cinema).Error; err != nil {
		return nil, c.returnError(ctx, err)
	}
	return cinema, nil
}

func (c *CinemaRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if err := c.db.WithContext(ctx).Where("id = ?", id).Update("is_active", false).Error; err != nil {
		return c.returnError(ctx, err)
	}
	return nil
}
