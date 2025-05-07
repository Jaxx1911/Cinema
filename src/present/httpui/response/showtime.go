package response

import (
	"TTCS/src/core/domain"
	"github.com/google/uuid"
	"time"
)

type ShowtimeResponse struct {
	Id        string    `json:"id"`
	MovieId   string    `json:"movie_id"`
	RoomId    string    `json:"room_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Price     float64   `json:"price"`
}

func ToShowtimeResponse(showtime domain.Showtime) ShowtimeResponse {
	return ShowtimeResponse{
		Id:        showtime.ID.String(),
		MovieId:   showtime.MovieID.String(),
		RoomId:    showtime.RoomID.String(),
		StartTime: showtime.StartTime,
		EndTime:   showtime.EndTime,
		Price:     showtime.Price,
	}
}

func ToListShowtimeResponse(showtimes []*domain.Showtime) []ShowtimeResponse {
	var list []ShowtimeResponse
	for _, v := range showtimes {
		list = append(list, ToShowtimeResponse(*v))
	}
	return list
}

type ShowtimeFullDetail struct {
	Showtime ShowtimeResponse    `json:"showtime"`
	Movie    MovieDetailResponse `json:"movie"`
	Room     Room                `json:"room"`
}

func ToShowtimeWithMovieAndRoom(showtime domain.Showtime) ShowtimeFullDetail {
	return ShowtimeFullDetail{
		Showtime: ToShowtimeResponse(showtime),
		Movie:    *ToMovieDetailResponse(&showtime.Movie),
		Room:     ToRoomResponse(&showtime.Room),
	}
}

func ToListShowtimeWithMovieAndRoom(showtimes []*domain.Showtime) []ShowtimeFullDetail {
	var list []ShowtimeFullDetail
	for _, v := range showtimes {
		list = append(list, ToShowtimeWithMovieAndRoom(*v))
	}
	return list
}

type ShowtimeWithRoom struct {
	Id        string    `json:"id"`
	MovieId   string    `json:"movie_id"`
	RoomId    string    `json:"room_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Price     float64   `json:"price"`

	Room Room `json:"room"`
}

type ShowtimeDetail struct {
	ShowtimeWithRoom
	Tickets []Ticket `json:"tickets"`
}

type Ticket struct {
	ID         uuid.UUID  `json:"id"`
	OrderID    *uuid.UUID `json:"order_id"`
	ShowtimeID uuid.UUID  `json:"showtime_id"`
	SeatID     uuid.UUID  `json:"seat_id"`
	Status     string     `json:"status"`
}

func ToShowtimeWithRoom(showtime *domain.Showtime) ShowtimeWithRoom {
	return ShowtimeWithRoom{
		Id:        showtime.ID.String(),
		MovieId:   showtime.MovieID.String(),
		RoomId:    showtime.RoomID.String(),
		StartTime: showtime.StartTime,
		EndTime:   showtime.EndTime,
		Price:     showtime.Price,

		Room: ToRoomResponse(&showtime.Room),
	}
}

func ToListShowtimeWithRoom(showtimes []*domain.Showtime) []ShowtimeWithRoom {
	var list []ShowtimeWithRoom
	for _, v := range showtimes {
		list = append(list, ToShowtimeWithRoom(v))
	}
	return list
}

func ToShowtimeDetailResponse(showtime *domain.Showtime) ShowtimeDetail {
	return ShowtimeDetail{
		ShowtimeWithRoom: ToShowtimeWithRoom(showtime),
		Tickets:          ToListTicketResponse(showtime.Tickets),
	}
}

func ToTicketResponse(ticket domain.Ticket) Ticket {
	return Ticket{
		ID:         ticket.ID,
		OrderID:    ticket.OrderID,
		ShowtimeID: ticket.ShowtimeID,
		SeatID:     ticket.SeatID,
		Status:     ticket.Status,
	}
}

func ToListTicketResponse(tickets []domain.Ticket) []Ticket {
	var list []Ticket
	for _, v := range tickets {
		list = append(list, ToTicketResponse(v))
	}
	return list
}

type TicketWithSeat struct {
	ID         uuid.UUID  `json:"id"`
	OrderID    *uuid.UUID `json:"order_id"`
	ShowtimeID uuid.UUID  `json:"showtime_id"`
	SeatID     uuid.UUID  `json:"seat_id"`
	Status     string     `json:"status"`

	Seat Seat `json:"seat"`
}

func ToTicketWithSeat(ticket domain.Ticket) TicketWithSeat {
	return TicketWithSeat{
		ID:         ticket.ID,
		OrderID:    ticket.OrderID,
		ShowtimeID: ticket.ShowtimeID,
		SeatID:     ticket.SeatID,
		Status:     ticket.Status,

		Seat: ToSeatResponse(&ticket.Seat),
	}
}

func ToListTicketWithSeat(seats []domain.Ticket) []TicketWithSeat {
	var list []TicketWithSeat
	for _, v := range seats {
		list = append(list, ToTicketWithSeat(v))
	}
	return list
}
