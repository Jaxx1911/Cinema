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

func (r RoomRepo) GetById(ctx context.Context, roomID string) (*domain.Room, error) {
	uid, err := uuid.Parse(roomID)
	if err != nil {
		return nil, r.returnError(ctx, err)
	}
	room := &domain.Room{}
	if err = r.db.Preload("Seats").First(room, "id = ?", uid).Error; err != nil {
		return nil, r.returnError(ctx, err)
	}
	return room, nil
}

func (r RoomRepo) GetListByCinemaId(ctx context.Context, cinemaId string) ([]*domain.Room, error) {
	uid, err := uuid.Parse(cinemaId)
	if err != nil {
		return nil, r.returnError(ctx, err)
	}
	var rooms []*domain.Room
	if err = r.db.Where("cinema_id = ?", uid).Where("is_active = ?", true).Find(&rooms).Error; err != nil {
		return nil, r.returnError(ctx, err)
	}
	return rooms, nil
}

func (r RoomRepo) Deactivate(ctx context.Context, id uuid.UUID, isActive bool) error {
	if err := r.db.Where("id = ?", id).Set("is_active", isActive).Error; err != nil {
		return r.returnError(ctx, err)
	}
	return nil
}

func NewRoomRepo(baseRepo *BaseRepo) domain.RoomRepo {
	return &RoomRepo{BaseRepo: baseRepo}
}
