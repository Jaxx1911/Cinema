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
