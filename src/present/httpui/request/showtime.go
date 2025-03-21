package request

type CreateShowtime struct {
	MovieId   string  `json:"movie_id"`
	RoomId    string  `json:"room_id"`
	StartTime string  `json:"start_time"`
	Price     float64 `json:"price"`
}

type GetShowtimesByFilter struct {
	CinemaId string `json:"cinema_id,omitempty"`
	MovieId  string `json:"movie_id,omitempty"`
	RoomId   string `json:"room_id,omitempty"`
	Date     string `json:"date,omitempty"`
}
