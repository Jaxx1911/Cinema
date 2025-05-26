package dto

import (
	"TTCS/src/core/domain"
)

type ShowtimeAvailabilityResponse struct {
	IsAvailable bool              `json:"is_available"`
	Conflicts   []domain.Showtime `json:"conflicts,omitempty"`
}

type ShowtimeAvailabilityResult struct {
	ShowtimeAvailabilityResponse
	MovieId   string `json:"movie_id"`
	RoomId    string `json:"room_id"`
	StartTime string `json:"start_time"`
}

type ShowtimesAvailabilityResponse struct {
	Results []ShowtimeAvailabilityResult `json:"results"`
}

type CreateShowtimeResult struct {
	Success   bool             `json:"success"`
	Showtime  *domain.Showtime `json:"showtime,omitempty"`
	Error     string           `json:"error,omitempty"`
	MovieId   string           `json:"movie_id"`
	RoomId    string           `json:"room_id"`
	StartTime string           `json:"start_time"`
}

type CreateShowtimesResponse struct {
	Results []CreateShowtimeResult `json:"results"`
	Summary struct {
		Total   int `json:"total"`
		Success int `json:"success"`
		Failed  int `json:"failed"`
	} `json:"summary"`
}
