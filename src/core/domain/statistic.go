package domain

import (
	"time"

	"github.com/google/uuid"
)

// MovieStatistic represents statistical data for a movie
type MovieStatistic struct {
	MovieID       uuid.UUID `json:"movie_id"`
	MovieTitle    string    `json:"movie_title"`
	TicketsSold   int       `json:"tickets_sold"`
	TotalRevenue  float64   `json:"total_revenue"`
	AveragePrice  float64   `json:"average_price"`
	ShowtimeCount int       `json:"showtime_count"`
	OccupancyRate float64   `json:"occupancy_rate"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
}

// CinemaStatistic represents statistical data for a cinema
type CinemaStatistic struct {
	CinemaID      uuid.UUID `json:"cinema_id"`
	CinemaName    string    `json:"cinema_name"`
	TicketRevenue float64   `json:"ticket_revenue"`
	ComboRevenue  float64   `json:"combo_revenue"`
	TotalRevenue  float64   `json:"total_revenue"`
	TicketsSold   int       `json:"tickets_sold"`
	ShowtimeCount int       `json:"showtime_count"`
	OccupancyRate float64   `json:"occupancy_rate"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
}

// ComboStatistic represents statistical data for a combo
type ComboStatistic struct {
	ComboID           uuid.UUID `json:"combo_id"`
	ComboName         string    `json:"combo_name"`
	ComboDescription  string    `json:"combo_description"`
	ComboPrice        float64   `json:"combo_price"`
	QuantitySold      int       `json:"quantity_sold"`
	TotalRevenue      float64   `json:"total_revenue"`
	PercentageOfTotal float64   `json:"percentage_of_total"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
}
