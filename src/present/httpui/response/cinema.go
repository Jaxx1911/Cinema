package response

import (
	"TTCS/src/core/domain"
	"github.com/google/uuid"
)

type Cinema struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
	Phone   string    `json:"phone"`
}

func ToCinemaResponse(cinema *domain.Cinema) Cinema {
	return Cinema{
		ID:      cinema.ID,
		Name:    cinema.Name,
		Address: cinema.Address,
		Phone:   cinema.Phone,
	}
}

func ToListCinemaResponse(cinemas []*domain.Cinema) []Cinema {
	var listCinema []Cinema
	for _, cinema := range cinemas {
		listCinema = append(listCinema, ToCinemaResponse(cinema))
	}
	return listCinema
}
