package response

import (
	"TTCS/src/core/domain"
	"github.com/google/uuid"
)

type Room struct {
	ID          uuid.UUID `json:"id"`
	CinemaID    uuid.UUID `json:"cinema_id"`
	Name        string    `json:"name"`
	Capacity    int       `json:"capacity"`
	Type        string    `json:"type"`
	RowCount    int       `json:"row_count"`
	ColumnCount int       `json:"column_count"`
}

func ToRoomResponse(room *domain.Room) Room {
	return Room{
		ID:          room.ID,
		CinemaID:    room.CinemaID,
		Name:        room.Name,
		Capacity:    room.Capacity,
		Type:        room.Type,
		RowCount:    room.RowCount,
		ColumnCount: room.ColumnCount,
	}
}

func ToListRoomResponse(rooms []domain.Room) []Room {
	var list []Room
	for _, room := range rooms {
		list = append(list, ToRoomResponse(&room))
	}
	return list
}

type Seat struct {
	ID         uuid.UUID `json:"id"`
	RowNumber  string    `json:"row_number"`  // Hàng: A, B, C...
	SeatNumber int       `json:"seat_number"` // Số ghế trong hàng
	Type       string    `json:"type"`
}

func ToSeatResponse(seat *domain.Seat) Seat {
	return Seat{
		ID:         seat.ID,
		RowNumber:  seat.RowNumber,
		SeatNumber: seat.SeatNumber,
		Type:       seat.Type,
	}
}

func ToListSeatResponse(seats []domain.Seat) []Seat {
	var list []Seat
	for _, seat := range seats {
		list = append(list, ToSeatResponse(&seat))
	}
	return list
}
