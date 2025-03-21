package repo

import (
	"TTCS/src/core/domain"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomRepo struct {
	*BaseRepo
	db *gorm.DB
}

func (r RoomRepo) Create(ctx context.Context, room *domain.Room) error {
	if err := r.db.Create(room).Error; err != nil {
		return err
	}
	return nil
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
	if err = r.db.Where("cinema_id = ?", uid).Find(&rooms).Error; err != nil {
		return nil, r.returnError(ctx, err)
	}
	return rooms, nil
}

func NewRoomRepo(baseRepo *BaseRepo, db *gorm.DB) domain.RoomRepo {
	return &RoomRepo{BaseRepo: baseRepo, db: db}
}
