package request

import "github.com/google/uuid"

type CreateOrderRequest struct {
	ShowtimeID uuid.UUID           `json:"showtime_id"`
	DiscountId *uuid.UUID          `json:"discount_id"`
	TotalPrice float64             `json:"total_price"`
	Tickets    []uuid.UUID         `json:"tickets"`
	Combos     []ComboOrderRequest `json:"combos"`
}

type ComboOrderRequest struct {
	Id         uuid.UUID `json:"id"`
	Quantity   int       `json:"quantity"`
	TotalPrice float64   `json:"total_price"`
}
