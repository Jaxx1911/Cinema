package repo

import (
	"TTCS/src/core/domain"
	"TTCS/src/present/httpui/request"
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

func (r RoomRepo) GetListByCinemaId(ctx context.Context, id uuid.UUID) ([]domain.Room, error) {
	var rooms []domain.Room
	if err := r.db.Where("cinema_id = ?", id).Where("is_active = ?", true).Find(&rooms).Error; err != nil {
		return nil, r.returnError(ctx, err)
	}
	return rooms, nil
}

func (r RoomRepo) Deactivate(ctx context.Context, id uuid.UUID, isActive bool) error {
	if err := r.db.Model(&domain.Room{}).Where("id = ?", id).Update("is_active", isActive).Error; err != nil {
		return r.returnError(ctx, err)
	}
	return nil
}

func (r RoomRepo) Update(ctx context.Context, room *domain.Room) (*domain.Room, error) {
	if err := r.db.WithContext(ctx).Model(room).Updates(map[string]interface{}{
		"name":      room.Name,
		"type":      room.Type,
		"is_active": room.IsActive,
	}).Error; err != nil {
		return nil, r.returnError(ctx, err)
	}

	return r.GetById(ctx, room.ID)
}

func (r *RoomRepo) GetList(ctx context.Context, page request.GetListRoom) ([]*domain.Room, int64, error) {
	var rooms []*domain.Room
	var total int64

	query := r.db.Model(&domain.Room{}).
		Preload("Cinema").
		Preload("Seats")

	if page.CinemaID != "" {
		query = query.Where("cinema_id = ?", uuid.MustParse(page.CinemaID))
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, r.returnError(ctx, err)
	}

	limit, offset := r.toLimitOffset(ctx, page.Page)
	if err := query.WithContext(ctx).Offset(offset).Limit(limit).Find(&rooms).Error; err != nil {
		return nil, 0, r.returnError(ctx, err)
	}

	return rooms, total, nil
}

func NewRoomRepo(baseRepo *BaseRepo) domain.RoomRepo {
	return &RoomRepo{BaseRepo: baseRepo}
}
