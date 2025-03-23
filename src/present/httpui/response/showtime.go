package response

import (
	"TTCS/src/core/domain"
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

type ShowtimeWithRoom struct {
	Id        string    `json:"id"`
	MovieId   string    `json:"movie_id"`
	RoomId    string    `json:"room_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Price     float64   `json:"price"`

	Room Room `json:"room"`
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
