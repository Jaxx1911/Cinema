package service

import (
	"TTCS/src/common/fault"
	"TTCS/src/core/domain"
	"context"
	"github.com/google/uuid"
)

type SeatService struct {
	seatRepo domain.SeatRepo
}

func NewSeatService(seatRepo domain.SeatRepo) *SeatService {
	return &SeatService{seatRepo: seatRepo}
}

func (s *SeatService) GetSeat(ctx context.Context, id string) (*domain.Seat, error) {
	caller := "SeatService.GetSeat"
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fault.Wrapf(err, "cannot parse id %v", caller)
	}
	seat, err := s.seatRepo.GetById(ctx, uid)
	if err != nil {
		return nil, err
	}
	return seat, nil
}

func (s *SeatService) GetByRoomId(ctx context.Context, roomId string) ([]domain.Seat, error) {
	caller := "SeatService.GetByRoomId"
	uid, err := uuid.Parse(roomId)
	if err != nil {
		return nil, fault.Wrapf(err, "cannot parse id %v", caller)
	}
	seats, err := s.seatRepo.GetByRoomID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return seats, nil
}
