package repo

import (
	"TTCS/src/core/domain"
	"context"
	"github.com/google/uuid"
)

type RoomRepo struct {
	*BaseRepo
}

func (r RoomRepo) Create(ctx context.Context, room *domain.Room) (*domain.Room, error) {
	if err := r.db.Create(room).Error; err != nil {
		return nil, r.returnError(ctx, err)
	}
	return room, nil
}

func (r RoomRepo) GetById(ctx context.Context, id uuid.UUID) (*domain.Room, error) {
	room := &domain.Room{}
	if err := r.db.Preload("Seats").First(room, "id = ?", id).Error; err != nil {
		return nil, r.returnError(ctx, err)
	}
	return room, nil
}

func (r RoomRepo) GetListByCinemaId(ctx context.Context, id uuid.UUID) ([]*domain.Room, error) {
	var rooms []*domain.Room
	if err := r.db.Where("cinema_id = ?", id).Where("is_active = ?", true).Find(&rooms).Error; err != nil {
		return nil, r.returnError(ctx, err)
	}
	return rooms, nil
}

func (r RoomRepo) Deactivate(ctx context.Context, id uuid.UUID, isActive bool) error {
	if err := r.db.Where("id = ?", id).Update("is_active", isActive).Error; err != nil {
		return r.returnError(ctx, err)
	}
	return nil
}

func NewRoomRepo(baseRepo *BaseRepo) domain.RoomRepo {
	return &RoomRepo{BaseRepo: baseRepo}
}
