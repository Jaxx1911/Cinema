package response

import (
	"TTCS/src/core/domain"
	"github.com/google/uuid"
)

type Cinema struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	Phone        string    `json:"phone"`
	OpeningHours string    `json:"opening_hours"`
}

type CinemaWithFacilities struct {
	Cinema
	Rooms []Room `json:"rooms"`
}

func ToCinemaResponse(cinema *domain.Cinema) Cinema {
	return Cinema{
		ID:           cinema.ID,
		Name:         cinema.Name,
		Address:      cinema.Address,
		Phone:        cinema.Phone,
		OpeningHours: cinema.OpeningHours,
	}
}

func ToListCinemaResponse(cinemas []*domain.Cinema) []Cinema {
	var listCinema []Cinema
	for _, cinema := range cinemas {
		listCinema = append(listCinema, ToCinemaResponse(cinema))
	}
	return listCinema
}

func ToCinemaWithFacilitiesResponse(cinema *domain.Cinema) CinemaWithFacilities {
	return CinemaWithFacilities{
		ToCinemaResponse(cinema),
		ToListRoomResponse(cinema.Rooms),
	}
}
func ToListCinemaWithFacilitiesResponse(cinemas []*domain.Cinema) []CinemaWithFacilities {
	var listCinemaWithFacilities []CinemaWithFacilities
	for _, cinema := range cinemas {
		listCinemaWithFacilities = append(listCinemaWithFacilities, ToCinemaWithFacilitiesResponse(cinema))
	}
	return listCinemaWithFacilities
}
