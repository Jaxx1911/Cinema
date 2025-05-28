package service

import (
	"TTCS/src/common/fault"
	"TTCS/src/core/domain"
	"TTCS/src/present/httpui/request"
	"TTCS/src/present/httpui/response"
	"context"
	"time"

	"github.com/google/uuid"
)

type StatisticService struct {
	orderRepo      domain.OrderRepo
	ticketRepo     domain.TicketRepo
	showtimeRepo   domain.ShowtimeRepo
	movieRepo      domain.MovieRepo
	cinemaRepo     domain.CinemaRepo
	roomRepo       domain.RoomRepo
	seatRepo       domain.SeatRepo
	comboRepo      domain.ComboRepository
	orderComboRepo domain.OrderComboRepository
}

func NewStatisticService(
	orderRepo domain.OrderRepo,
	ticketRepo domain.TicketRepo,
	showtimeRepo domain.ShowtimeRepo,
	movieRepo domain.MovieRepo,
	cinemaRepo domain.CinemaRepo,
	roomRepo domain.RoomRepo,
	seatRepo domain.SeatRepo,
	comboRepo domain.ComboRepository,
	orderComboRepo domain.OrderComboRepository,
) *StatisticService {
	return &StatisticService{
		orderRepo:      orderRepo,
		ticketRepo:     ticketRepo,
		showtimeRepo:   showtimeRepo,
		movieRepo:      movieRepo,
		cinemaRepo:     cinemaRepo,
		roomRepo:       roomRepo,
		seatRepo:       seatRepo,
		comboRepo:      comboRepo,
		orderComboRepo: orderComboRepo,
	}
}

func (s *StatisticService) GetMovieRevenue(ctx context.Context, req request.StatisticDateRange) (*response.MovieRevenueResponse, error) {
	caller := "StatisticService.GetMovieRevenue"

	// Lấy tất cả orders trong khoảng thời gian
	orders, err := s.orderRepo.GetOrdersByDateRange(ctx, req.StartDate, req.EndDate)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to get orders", caller)
	}

	movieStats := make(map[uuid.UUID]*movieStatData)
	var totalRevenue float64
	var totalTicketsSold int

	for _, order := range orders {
		if order.Status != "success" {
			continue
		}

		// Lấy showtime để biết movie
		showtime, err := s.showtimeRepo.GetById(ctx, order.ShowtimeID)
		if err != nil {
			continue
		}

		// Lấy tickets của order này
		tickets, err := s.ticketRepo.FindByOrderID(ctx, order.ID)
		if err != nil {
			continue
		}

		ticketCount := len(tickets)
		if ticketCount == 0 {
			continue
		}

		// Tính doanh thu vé (loại trừ combo)
		orderCombos, _ := s.orderComboRepo.GetByOrderID(ctx, order.ID)
		comboRevenue := float64(0)
		for _, combo := range orderCombos {
			comboRevenue += combo.TotalPrice
		}
		ticketRevenue := order.TotalPrice - comboRevenue

		if _, exists := movieStats[showtime.MovieID]; !exists {
			movieStats[showtime.MovieID] = &movieStatData{
				MovieID:       showtime.MovieID,
				TicketsSold:   0,
				TotalRevenue:  0,
				ShowtimeCount: 0,
				ShowtimeIDs:   make(map[uuid.UUID]bool),
			}
		}

		stat := movieStats[showtime.MovieID]
		stat.TicketsSold += ticketCount
		stat.TotalRevenue += ticketRevenue
		if !stat.ShowtimeIDs[showtime.ID] {
			stat.ShowtimeCount++
			stat.ShowtimeIDs[showtime.ID] = true
		}

		totalRevenue += ticketRevenue
		totalTicketsSold += ticketCount
	}

	// Chuyển đổi sang response format
	var movieItems []response.MovieRevenueItem
	for movieID, stat := range movieStats {
		movie, err := s.movieRepo.GetById(ctx, movieID)
		if err != nil {
			continue
		}

		// Tính tỷ lệ lấp đầy
		occupancyRate := s.calculateMovieOccupancyRate(ctx, movieID, req.StartDate, req.EndDate, stat.TicketsSold)

		averagePrice := float64(0)
		if stat.TicketsSold > 0 {
			averagePrice = stat.TotalRevenue / float64(stat.TicketsSold)
		}

		movieItems = append(movieItems, response.MovieRevenueItem{
			MovieID:       movieID,
			MovieTitle:    movie.Title,
			TicketsSold:   stat.TicketsSold,
			AveragePrice:  averagePrice,
			ShowtimeCount: stat.ShowtimeCount,
			OccupancyRate: occupancyRate,
		})
	}

	return &response.MovieRevenueResponse{
		Movies: movieItems,
		Summary: response.MovieRevenueSummary{
			TotalRevenue:     totalRevenue,
			TotalTicketsSold: totalTicketsSold,
		},
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}, nil
}

func (s *StatisticService) GetCinemaRevenue(ctx context.Context, req request.StatisticDateRange) (*response.CinemaRevenueResponse, error) {
	caller := "StatisticService.GetCinemaRevenue"

	orders, err := s.orderRepo.GetOrdersByDateRange(ctx, req.StartDate, req.EndDate)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to get orders", caller)
	}

	cinemaStats := make(map[uuid.UUID]*cinemaStatData)
	var totalRevenue, totalTicketRevenue, totalComboRevenue float64

	for _, order := range orders {
		if order.Status != "success" {
			continue
		}

		showtime, err := s.showtimeRepo.GetById(ctx, order.ShowtimeID)
		if err != nil {
			continue
		}

		room, err := s.roomRepo.GetById(ctx, showtime.RoomID)
		if err != nil {
			continue
		}

		tickets, err := s.ticketRepo.FindByOrderID(ctx, order.ID)
		if err != nil {
			continue
		}

		ticketCount := len(tickets)
		if ticketCount == 0 {
			continue
		}

		// Tính doanh thu combo
		orderCombos, _ := s.orderComboRepo.GetByOrderID(ctx, order.ID)
		comboRevenue := float64(0)
		for _, combo := range orderCombos {
			comboRevenue += combo.TotalPrice
		}
		ticketRevenue := order.TotalPrice - comboRevenue

		if _, exists := cinemaStats[room.CinemaID]; !exists {
			cinemaStats[room.CinemaID] = &cinemaStatData{
				CinemaID:      room.CinemaID,
				TicketRevenue: 0,
				ComboRevenue:  0,
				TicketsSold:   0,
				ShowtimeCount: 0,
				ShowtimeIDs:   make(map[uuid.UUID]bool),
			}
		}

		stat := cinemaStats[room.CinemaID]
		stat.TicketRevenue += ticketRevenue
		stat.ComboRevenue += comboRevenue
		stat.TicketsSold += ticketCount
		if !stat.ShowtimeIDs[showtime.ID] {
			stat.ShowtimeCount++
			stat.ShowtimeIDs[showtime.ID] = true
		}

		totalTicketRevenue += ticketRevenue
		totalComboRevenue += comboRevenue
		totalRevenue += order.TotalPrice
	}

	// Chuyển đổi sang response format
	var cinemaItems []response.CinemaRevenueItem
	for cinemaID, stat := range cinemaStats {
		cinema, err := s.cinemaRepo.FindByID(ctx, cinemaID)
		if err != nil {
			continue
		}

		// Tính tỷ lệ lấp đầy
		occupancyRate := s.calculateCinemaOccupancyRate(ctx, cinemaID, req.StartDate, req.EndDate, stat.TicketsSold)

		cinemaItems = append(cinemaItems, response.CinemaRevenueItem{
			CinemaID:      cinemaID,
			CinemaName:    cinema.Name,
			TicketRevenue: stat.TicketRevenue,
			ComboRevenue:  stat.ComboRevenue,
			TicketsSold:   stat.TicketsSold,
			ShowtimeCount: stat.ShowtimeCount,
			OccupancyRate: occupancyRate,
		})
	}

	return &response.CinemaRevenueResponse{
		Cinemas: cinemaItems,
		Summary: response.CinemaRevenueSummary{
			TotalRevenue:       totalRevenue,
			TotalTicketRevenue: totalTicketRevenue,
			TotalComboRevenue:  totalComboRevenue,
		},
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}, nil
}

func (s *StatisticService) GetComboStatistics(ctx context.Context, req request.StatisticDateRange) (*response.ComboStatisticResponse, error) {
	caller := "StatisticService.GetComboStatistics"

	// Lấy tất cả orders trong khoảng thời gian
	orders, err := s.orderRepo.GetOrdersByDateRange(ctx, req.StartDate, req.EndDate)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to get orders", caller)
	}

	comboStats := make(map[uuid.UUID]*comboStatData)
	var totalRevenue float64
	var totalQuantitySold int

	for _, order := range orders {
		// Chỉ xử lý orders thành công
		if order.Status != "success" {
			continue
		}

		// Lấy order combos của order này
		orderCombos, err := s.orderComboRepo.GetByOrderID(ctx, order.ID)
		if err != nil {
			continue
		}

		for _, orderCombo := range orderCombos {
			if _, exists := comboStats[orderCombo.ComboID]; !exists {
				comboStats[orderCombo.ComboID] = &comboStatData{
					ComboID:      orderCombo.ComboID,
					QuantitySold: 0,
					TotalRevenue: 0,
				}
			}

			stat := comboStats[orderCombo.ComboID]
			stat.QuantitySold += orderCombo.Quantity
			stat.TotalRevenue += orderCombo.TotalPrice

			totalQuantitySold += orderCombo.Quantity
			totalRevenue += orderCombo.TotalPrice
		}
	}

	// Chuyển đổi sang response format
	var comboItems []response.ComboStatisticItem
	for comboID, stat := range comboStats {
		combo, err := s.comboRepo.FindByID(ctx, comboID)
		if err != nil {
			continue
		}

		// Tính phần trăm của tổng doanh thu
		percentageOfTotal := float64(0)
		if totalRevenue > 0 {
			percentageOfTotal = (stat.TotalRevenue / totalRevenue) * 100
		}

		comboItems = append(comboItems, response.ComboStatisticItem{
			ID:                comboID,
			Name:              combo.Name,
			Description:       combo.Description,
			Price:             combo.Price,
			Quantity:          stat.QuantitySold,
			Revenue:           stat.TotalRevenue,
			PercentageOfTotal: percentageOfTotal,
		})
	}

	return &response.ComboStatisticResponse{
		Combos: comboItems,
		Summary: response.ComboStatisticSummary{
			TotalRevenue:      totalRevenue,
			TotalQuantitySold: totalQuantitySold,
		},
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}, nil
}

// Helper functions
func (s *StatisticService) calculateMovieOccupancyRate(ctx context.Context, movieID uuid.UUID, startDate, endDate time.Time, ticketsSold int) float64 {
	// Lấy tất cả showtimes của movie trong khoảng thời gian
	showtimes, err := s.showtimeRepo.GetByMovieIDAndDateRange(ctx, movieID, startDate, endDate)
	if err != nil {
		return 0
	}

	totalCapacity := 0
	for _, showtime := range showtimes {
		room, err := s.roomRepo.GetById(ctx, showtime.RoomID)
		if err != nil {
			continue
		}
		totalCapacity += room.Capacity
	}

	if totalCapacity == 0 {
		return 0
	}

	return float64(ticketsSold) / float64(totalCapacity) * 100
}

func (s *StatisticService) calculateCinemaOccupancyRate(ctx context.Context, cinemaID uuid.UUID, startDate, endDate time.Time, ticketsSold int) float64 {
	// Lấy tất cả rooms của cinema
	rooms, err := s.roomRepo.GetListByCinemaId(ctx, cinemaID)
	if err != nil {
		return 0
	}

	totalCapacity := 0
	for _, room := range rooms {
		// Lấy showtimes của room trong khoảng thời gian
		showtimes, err := s.showtimeRepo.GetByRoomIDAndDateRange(ctx, room.ID, startDate, endDate)
		if err != nil {
			continue
		}
		totalCapacity += len(showtimes) * room.Capacity
	}

	if totalCapacity == 0 {
		return 0
	}

	return float64(ticketsSold) / float64(totalCapacity) * 100
}

// Helper structs
type movieStatData struct {
	MovieID       uuid.UUID
	TicketsSold   int
	TotalRevenue  float64
	ShowtimeCount int
	ShowtimeIDs   map[uuid.UUID]bool
}

type cinemaStatData struct {
	CinemaID      uuid.UUID
	TicketRevenue float64
	ComboRevenue  float64
	TicketsSold   int
	ShowtimeCount int
	ShowtimeIDs   map[uuid.UUID]bool
}

type comboStatData struct {
	ComboID      uuid.UUID
	QuantitySold int
	TotalRevenue float64
}
