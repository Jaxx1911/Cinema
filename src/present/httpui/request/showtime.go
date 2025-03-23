package request

type CreateShowtime struct {
	MovieId   string  `json:"movie_id"`
	RoomId    string  `json:"room_id"`
	StartTime string  `json:"start_time"`
	Price     float64 `json:"price"`
}

type GetShowtimesByUserFilter struct {
	CinemaId string `json:"cinema_id" binding:"required" form:"cinema_id"`
	MovieId  string `json:"movie_id" binding:"required" form:"movie_id"`
	Day      string `json:"day" binding:"required" form:"day"`
}
