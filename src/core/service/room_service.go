package service

import (
	"TTCS/src/core/domain"
	"TTCS/src/present/httpui/request"
	"context"
	"github.com/google/uuid"
)

type RoomService struct {
	roomRepo domain.RoomRepo
	seatRepo domain.SeatRepo
}

func NewRoomService(roomRepo domain.RoomRepo) *RoomService {
	return &RoomService{roomRepo: roomRepo}
}

func (r *RoomService) Create(ctx context.Context, req request.CreateRoomReq) (*domain.Room, error) {
	room, err := r.roomRepo.Create(ctx, &domain.Room{
		CinemaID:    req.CinemaId,
		Name:        req.Name,
		Capacity:    req.ColumnCount * req.RowCount,
		Type:        req.Type,
		RowCount:    req.RowCount,
		ColumnCount: req.ColumnCount,
		IsActive:    true,
	})

	if err != nil {
		return nil, err
	}
	for _, v := range room.Seats {
		err := r.seatRepo.Create(ctx, &domain.Seat{
			RoomID:     room.ID,
			RowNumber:  v.RowNumber,
			SeatNumber: v.SeatNumber,
			Type:       v.Type,
		})
		if err != nil {
			return nil, err
		}
	}
	return room, nil
}

func (r *RoomService) Deactivate(ctx context.Context, id uuid.UUID, isActive bool) error {
	if err := r.roomRepo.Deactivate(ctx, id, isActive); err != nil {
		return err
	}
	return nil
}
