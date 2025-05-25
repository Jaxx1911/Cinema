package dto

import (
	"TTCS/src/core/domain"
)

type ShowtimeAvailabilityResponse struct {
	IsAvailable bool              `json:"is_available"`
	Conflicts   []domain.Showtime `json:"conflicts,omitempty"`
}
