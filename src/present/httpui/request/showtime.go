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
