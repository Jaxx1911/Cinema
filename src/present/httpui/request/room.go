package request

import "github.com/google/uuid"

type CreateRoomReq struct {
	CinemaId    uuid.UUID `json:"cinema_id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	RowCount    int       `json:"row_count"`
	ColumnCount int       `json:"column_count"`
}
