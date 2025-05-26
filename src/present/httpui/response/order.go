package response

import (
	"TTCS/src/core/domain"

	"github.com/google/uuid"
)

type OrderResponse struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"user_id"`
	ShowtimeID uuid.UUID  `json:"showtime_id"`
	DiscountID *uuid.UUID `json:"discount_id"`
	Status     string     `json:"status"`
	TotalPrice float64    `json:"total_price"`

	Tickets []TicketWithSeat `json:"tickets"`
}

type OrderWithQrResponse struct {
	Order  OrderResponse `json:"order"`
	QrText string        `json:"qr_text"`
}

func ToOrderResponse(order domain.Order) OrderResponse {
	return OrderResponse{
		ID:         order.ID,
		UserID:     order.UserID,
		ShowtimeID: order.ShowtimeID,
		DiscountID: order.DiscountID,
		Status:     order.Status,
		TotalPrice: order.TotalPrice,

		Tickets: ToListTicketWithSeat(order.Tickets),
	}
}

func ToOrdersResponse(orders []domain.Order) []OrderResponse {
	var orderResponses []OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, ToOrderResponse(order))
	}
	return orderResponses
}

type OrderDetailResponse struct {
	OrderResponse
	ShowtimeFullDetail
	Combos []ComboWithQuantity `json:"combos"`
}

func ToOrderDetailResponse(order domain.Order) OrderDetailResponse {
	var combos []*domain.Combo
	for _, v := range order.OrderCombos {
		combos = append(combos, &v.Combo)
	}

	return OrderDetailResponse{
		OrderResponse:      ToOrderResponse(order),
		ShowtimeFullDetail: ToShowtimeWithMovieAndRoom(order.Showtime),
		Combos:             ToListComboWithQuantity(order.OrderCombos),
	}
}
