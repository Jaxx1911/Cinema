package service

import (
	"TTCS/src/common/fault"
	"TTCS/src/core/domain"
	"TTCS/src/present/httpui/request"
	"context"

	"github.com/google/uuid"
)

type RoomService struct {
	roomRepo domain.RoomRepo
	seatRepo domain.SeatRepo
}

func NewRoomService(roomRepo domain.RoomRepo, seatRepo domain.SeatRepo) *RoomService {
	return &RoomService{
		roomRepo: roomRepo,
		seatRepo: seatRepo,
	}
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
	for _, v := range req.Seats {
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

func (r *RoomService) Update(ctx context.Context, id uuid.UUID, req request.UpdateRoomReq) (*domain.Room, error) {
	// Get the existing room
	room, err := r.roomRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	room.Name = req.Name
	room.Type = req.Type
	room.IsActive = req.IsActive

	// Save changes for the room
	room, err = r.roomRepo.Update(ctx, room)
	if err != nil {
		return nil, err
	}

	// Update seats if provided
	if len(req.Seats) > 0 {
		for _, seatReq := range req.Seats {
			seat := &domain.Seat{
				ID:         seatReq.ID,
				RoomID:     room.ID,
				RowNumber:  seatReq.RowNumber,
				SeatNumber: seatReq.SeatNumber,
				Type:       seatReq.Type,
			}

			if err := r.seatRepo.UpdateSeat(ctx, seat); err != nil {
				return nil, err
			}
		}

		// Reload room with updated seats
		room, err = r.roomRepo.GetById(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	return room, nil
}

func (r *RoomService) GetRoomById(ctx context.Context, id uuid.UUID) (*domain.Room, error) {
	return r.roomRepo.GetById(ctx, id)
}

func (r *RoomService) GetList(ctx context.Context, page request.GetListRoom) ([]*domain.Room, int64, error) {
	caller := "RoomService.GetList"

	rooms, total, err := r.roomRepo.GetList(ctx, page)
	if err != nil {
		return nil, 0, fault.Wrapf(err, "[%v] failed to get rooms", caller)
	}

	return rooms, total, nil
}

func (r *RoomService) GetListByCinemaId(ctx context.Context, cinemaId uuid.UUID) ([]domain.Room, error) {
	return r.roomRepo.GetListByCinemaId(ctx, cinemaId)
}
