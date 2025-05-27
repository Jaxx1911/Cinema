package response

import (
	"time"

	"github.com/google/uuid"
)

// Movie Revenue Statistics
type MovieRevenueItem struct {
	MovieID       uuid.UUID `json:"movie_id"`
	MovieTitle    string    `json:"movie_title"`
	TicketsSold   int       `json:"tickets_sold"`
	AveragePrice  float64   `json:"average_price"`
	ShowtimeCount int       `json:"showtime_count"`
	OccupancyRate float64   `json:"occupancy_rate"`
}

type MovieRevenueSummary struct {
	TotalRevenue     float64 `json:"total_revenue"`
	TotalTicketsSold int     `json:"total_tickets_sold"`
}

type MovieRevenueResponse struct {
	Movies    []MovieRevenueItem  `json:"movies"`
	Summary   MovieRevenueSummary `json:"summary"`
	StartDate time.Time           `json:"start_date"`
	EndDate   time.Time           `json:"end_date"`
}

// Cinema Revenue Statistics
type CinemaRevenueItem struct {
	CinemaID      uuid.UUID `json:"cinema_id"`
	CinemaName    string    `json:"cinema_name"`
	TicketRevenue float64   `json:"ticket_revenue"`
	ComboRevenue  float64   `json:"combo_revenue"`
	TicketsSold   int       `json:"tickets_sold"`
	ShowtimeCount int       `json:"showtime_count"`
	OccupancyRate float64   `json:"occupancy_rate"`
}

type CinemaRevenueSummary struct {
	TotalRevenue       float64 `json:"total_revenue"`
	TotalTicketRevenue float64 `json:"total_ticket_revenue"`
	TotalComboRevenue  float64 `json:"total_combo_revenue"`
}

type CinemaRevenueResponse struct {
	Cinemas   []CinemaRevenueItem  `json:"cinemas"`
	Summary   CinemaRevenueSummary `json:"summary"`
	StartDate time.Time            `json:"start_date"`
	EndDate   time.Time            `json:"end_date"`
}

// Combo Statistics
type ComboStatisticItem struct {
	ComboID      uuid.UUID `json:"combo_id"`
	ComboName    string    `json:"combo_name"`
	QuantitySold int       `json:"quantity_sold"`
	TotalRevenue float64   `json:"total_revenue"`
}

type ComboStatisticSummary struct {
	TotalRevenue      float64 `json:"total_revenue"`
	TotalQuantitySold int     `json:"total_quantity_sold"`
}

type ComboStatisticResponse struct {
	Combos  []ComboStatisticItem  `json:"combos"`
	Summary ComboStatisticSummary `json:"summary"`
}
