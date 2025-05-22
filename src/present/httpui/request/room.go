package request

import "github.com/google/uuid"

type CreateRoomReq struct {
	CinemaId    uuid.UUID    `json:"cinema_id"`
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	RowCount    int          `json:"row_count"`
	ColumnCount int          `json:"column_count"`
	Seats       []CreateSeat `json:"seats"`
}

type UpdateRoomReq struct {
	Name     string       `json:"name"`
	Type     string       `json:"type"`
	IsActive bool         `json:"is_active"`
	Seats    []UpdateSeat `json:"seats,omitempty"`
}

type CreateSeat struct {
	RowNumber  string `json:"row_number"` // Hàng: A, B, C...
	SeatNumber int    `json:"seat_number"`
	Type       string `json:"type"` // Standard, Premium, Couple
}

type UpdateSeat struct {
	ID         uuid.UUID `json:"id"`
	RowNumber  string    `json:"row_number"` // Hàng: A, B, C...
	SeatNumber int       `json:"seat_number"`
	Type       string    `json:"type"` // Standard, Premium, Couple
}

type GetListRoom struct {
	Page
	CinemaID string `form:"cinema_id"`
}
