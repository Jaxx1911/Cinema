package request

import "github.com/google/uuid"

type GetCinemaRequest struct {
	City string `json:"city" binding:"required" form:"city"`
}

func (g *GetCinemaRequest) MappingCity() {
	switch g.City {
	case "hanoi":
		g.City = "Hà Nội"
	case "hcm":
		g.City = "HCM"
	case "danang":
		g.City = "Đà Nẵng"
	}
}

type CreateCinemaRequest struct {
	Name         string `json:"name"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	OpeningHours string `json:"opening_hours"`
}

type UpdateCinemaRequest struct {
	Id           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	Phone        string    `json:"phone"`
	OpeningHours string    `json:"opening_hours"`
	IsActive     bool      `json:"is_active"`
}
