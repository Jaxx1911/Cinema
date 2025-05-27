package request

import "github.com/google/uuid"

type CreateShowtime struct {
	MovieId   uuid.UUID `json:"movie_id"`
	RoomId    uuid.UUID `json:"room_id"`
	StartTime string    `json:"start_time"`
	Price     float64   `json:"price"`
}

type GetShowtimesByUserFilter struct {
	CinemaId string `json:"cinema_id" binding:"required" form:"cinema_id"`
	MovieId  string `json:"movie_id" binding:"required" form:"movie_id"`
	Day      string `json:"day" binding:"required" form:"day"`
}

type GetShowtimesByCinemaIdFilter struct {
	CinemaId string `json:"cinema_id" binding:"required" form:"cinema_id"`
	Day      string `json:"day" binding:"required" form:"day"`
}

type GetShowtimesByRoomIdFilter struct {
	RoomId string `json:"room_id" binding:"required" form:"room_id"`
	Day    string `json:"day" binding:"required" form:"day"`
}

type GetListShowtime struct {
	Page
	MovieID  string `form:"movie_id"`
	RoomID   string `form:"room_id"`
	CinemaID string `form:"cinema_id"`
	FromDate string `form:"from_date"`
	ToDate   string `form:"to_date"`
}

type UpdateShowtime struct {
	MovieId   uuid.UUID `json:"movie_id"`
	RoomId    uuid.UUID `json:"room_id"`
	StartTime string    `json:"start_time"`
	Price     float64   `json:"price"`
}

type CheckShowtimeAvailability struct {
	Id        *uuid.UUID `json:"id"`
	MovieId   uuid.UUID  `json:"movie_id" binding:"required"`
	RoomId    uuid.UUID  `json:"room_id" binding:"required"`
	StartTime string     `json:"start_time" binding:"required"`
}

type CheckShowtimesAvailability struct {
	Showtimes []CheckShowtimeAvailability `json:"showtimes" binding:"required"`
}

type CreateShowtimes struct {
	Showtimes []CreateShowtime `json:"showtimes" binding:"required"`
}
