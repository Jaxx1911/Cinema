package repo

import (
	"TTCS/src/core/domain"
	"context"
	"github.com/google/uuid"
)

type SeatRepo struct {
	*BaseRepo
}

func (s SeatRepo) Create(ctx context.Context, seat *domain.Seat) error {
	if err := s.db.WithContext(ctx).Create(seat).Error; err != nil {
		return s.returnError(ctx, err)
	}
	return nil
}

func (s SeatRepo) GetById(ctx context.Context, seatID uuid.UUID) (*domain.Seat, error) {
	var seat domain.Seat
	if err := s.db.WithContext(ctx).Where("id = ?", seatID).First(&seat).Error; err != nil {
		return nil, s.returnError(ctx, err)
	}
	return &seat, nil
}

func (s SeatRepo) GetByRoomID(ctx context.Context, roomID uuid.UUID) ([]domain.Seat, error) {
	var seat []domain.Seat
	if err := s.db.WithContext(ctx).Where("room_id = ?", roomID).Find(&seat).Error; err != nil {
		return nil, s.returnError(ctx, err)
	}
	return seat, nil
}

func NewSeatRepo(baseRepo *BaseRepo) domain.SeatRepo {
	return &SeatRepo{
		BaseRepo: baseRepo,
	}
}
