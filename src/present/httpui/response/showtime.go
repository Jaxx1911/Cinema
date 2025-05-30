package response

import (
	"TTCS/src/core/domain"
	"TTCS/src/core/dto"
	"time"

	"github.com/google/uuid"
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

// New response struct for showtime list with sold tickets count
type ShowtimeWithRoomAndSoldTickets struct {
	Id          string    `json:"id"`
	MovieId     string    `json:"movie_id"`
	RoomId      string    `json:"room_id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Price       float64   `json:"price"`
	SoldTickets int       `json:"sold_tickets"`

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

type ShowtimeAvailabilityResponse struct {
	IsAvailable bool                `json:"is_available"`
	Conflicts   []ShowtimeWithMovie `json:"conflicts,omitempty"`
}

type ShowtimeWithMovie struct {
	Id        string    `json:"id"`
	MovieId   string    `json:"movie_id"`
	RoomId    string    `json:"room_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Price     float64   `json:"price"`
	MovieName string    `json:"movie_name"`
}

func ToShowtimeAvailabilityResponse(resp *dto.ShowtimeAvailabilityResponse) ShowtimeAvailabilityResponse {
	var conflicts []ShowtimeWithMovie
	for _, conflict := range resp.Conflicts {
		conflicts = append(conflicts, ShowtimeWithMovie{
			Id:        conflict.ID.String(),
			MovieId:   conflict.MovieID.String(),
			RoomId:    conflict.RoomID.String(),
			StartTime: conflict.StartTime,
			EndTime:   conflict.EndTime,
			Price:     conflict.Price,
			MovieName: conflict.Movie.Title,
		})
	}

	return ShowtimeAvailabilityResponse{
		IsAvailable: resp.IsAvailable,
		Conflicts:   conflicts,
	}
}

type ShowtimeAvailabilityResult struct {
	IsAvailable bool                `json:"is_available"`
	Conflicts   []ShowtimeWithMovie `json:"conflicts,omitempty"`
	MovieId     string              `json:"movie_id"`
	RoomId      string              `json:"room_id"`
	StartTime   string              `json:"start_time"`
}

type ShowtimesAvailabilityResponse struct {
	Results []ShowtimeAvailabilityResult `json:"results"`
}

type CreateShowtimeResult struct {
	Success   bool              `json:"success"`
	Showtime  *ShowtimeResponse `json:"showtime,omitempty"`
	Error     string            `json:"error,omitempty"`
	MovieId   string            `json:"movie_id"`
	RoomId    string            `json:"room_id"`
	StartTime string            `json:"start_time"`
}

type CreateShowtimesResponse struct {
	Results []CreateShowtimeResult `json:"results"`
	Summary struct {
		Total   int `json:"total"`
		Success int `json:"success"`
		Failed  int `json:"failed"`
	} `json:"summary"`
}

func ToShowtimesAvailabilityResponse(resp *dto.ShowtimesAvailabilityResponse) ShowtimesAvailabilityResponse {
	var results []ShowtimeAvailabilityResult
	for _, result := range resp.Results {
		var conflicts []ShowtimeWithMovie
		for _, conflict := range result.Conflicts {
			conflicts = append(conflicts, ShowtimeWithMovie{
				Id:        conflict.ID.String(),
				MovieId:   conflict.MovieID.String(),
				RoomId:    conflict.RoomID.String(),
				StartTime: conflict.StartTime,
				EndTime:   conflict.EndTime,
				Price:     conflict.Price,
				MovieName: conflict.Movie.Title,
			})
		}

		results = append(results, ShowtimeAvailabilityResult{
			IsAvailable: result.IsAvailable,
			Conflicts:   conflicts,
			MovieId:     result.MovieId,
			RoomId:      result.RoomId,
			StartTime:   result.StartTime,
		})
	}

	return ShowtimesAvailabilityResponse{
		Results: results,
	}
}

func ToCreateShowtimesResponse(resp *dto.CreateShowtimesResponse) CreateShowtimesResponse {
	var results []CreateShowtimeResult
	for _, result := range resp.Results {
		var showtimeResp *ShowtimeResponse
		if result.Showtime != nil {
			resp := ToShowtimeResponse(*result.Showtime)
			showtimeResp = &resp
		}

		results = append(results, CreateShowtimeResult{
			Success:   result.Success,
			Showtime:  showtimeResp,
			Error:     result.Error,
			MovieId:   result.MovieId,
			RoomId:    result.RoomId,
			StartTime: result.StartTime,
		})
	}

	return CreateShowtimesResponse{
		Results: results,
		Summary: struct {
			Total   int `json:"total"`
			Success int `json:"success"`
			Failed  int `json:"failed"`
		}{
			Total:   resp.Summary.Total,
			Success: resp.Summary.Success,
			Failed:  resp.Summary.Failed,
		},
	}
}

func ToShowtimeWithRoomAndSoldTickets(showtime *domain.Showtime) ShowtimeWithRoomAndSoldTickets {
	// Count sold tickets (tickets with status "success")
	soldTickets := 0
	for _, ticket := range showtime.Tickets {
		if ticket.Status == "success" {
			soldTickets++
		}
	}

	return ShowtimeWithRoomAndSoldTickets{
		Id:          showtime.ID.String(),
		MovieId:     showtime.MovieID.String(),
		RoomId:      showtime.RoomID.String(),
		StartTime:   showtime.StartTime,
		EndTime:     showtime.EndTime,
		Price:       showtime.Price,
		SoldTickets: soldTickets,

		Room: ToRoomResponse(&showtime.Room),
	}
}

func ToListShowtimeWithRoomAndSoldTickets(showtimes []*domain.Showtime) []ShowtimeWithRoomAndSoldTickets {
	var list []ShowtimeWithRoomAndSoldTickets
	for _, v := range showtimes {
		list = append(list, ToShowtimeWithRoomAndSoldTickets(v))
	}
	return list
}
