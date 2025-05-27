package adapters

import "TTCS/src/core/domain"

type Adapters struct {
	OrderRepo    domain.OrderRepo
	ShowtimeRepo domain.ShowtimeRepo
	CinemaRepo   domain.CinemaRepo
	ComboRepo    domain.ComboRepository
	DiscountRepo domain.DiscountRepository
	GenreRepo    domain.GenreRepo
	MovieRepo    domain.MovieRepo
	Payment      domain.PaymentRepo
	RoomRepo     domain.RoomRepo
	SeatRepo     domain.SeatRepo
	TicketRepo   domain.TicketRepo
	UserRepo     domain.UserRepo
}
