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
	discountRepo   domain.DiscountRepository
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
	discountRepo domain.DiscountRepository,
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
		discountRepo:   discountRepo,
	}
}

func (s *StatisticService) GetMovieRevenue(ctx context.Context, req request.StatisticDateRange) (*response.MovieRevenueResponse, error) {
	caller := "StatisticService.GetMovieRevenue"

	// Lấy tất cả orders trong khoảng thời gian
	orders, err := s.orderRepo.GetOrdersByDateRange(ctx, req.StartDate, req.EndDate)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to get orders", caller)
	}

	movieStats := make(map[uuid.UUID]*domain.MovieStatistic)
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

		// Tính doanh thu combo (xử lý discount nếu có)
		orderCombos, _ := s.orderComboRepo.GetByOrderID(ctx, order.ID)
		comboRevenue := float64(0)
		for _, combo := range orderCombos {
			comboRevenue += combo.TotalPrice
		}

		var ticketRevenue float64
		// Nếu có discount, cần tính lại comboRevenue gốc trước khi áp discount
		if order.DiscountID != nil {
			// Lấy thông tin discount
			discount, err := s.discountRepo.GetDiscount(ctx, *order.DiscountID)
			if err != nil {
				return nil, fault.Wrapf(err, "[%v] failed to get discount", caller)
			}

			ticketRevenue = order.TotalPrice - comboRevenue*(1-discount.Percentage/100)
		} else {
			ticketRevenue = order.TotalPrice - comboRevenue
		}

		if _, exists := movieStats[showtime.MovieID]; !exists {
			movie, err := s.movieRepo.GetById(ctx, showtime.MovieID)
			if err != nil {
				continue
			}

			movieStats[showtime.MovieID] = &domain.MovieStatistic{
				MovieID:       showtime.MovieID,
				MovieTitle:    movie.Title,
				TicketsSold:   0,
				TotalRevenue:  0,
				AveragePrice:  0,
				ShowtimeCount: 0,
				OccupancyRate: 0,
				StartDate:     req.StartDate,
				EndDate:       req.EndDate,
			}
		}

		stat := movieStats[showtime.MovieID]
		stat.TicketsSold += ticketCount
		stat.TotalRevenue += ticketRevenue

		totalRevenue += ticketRevenue
		totalTicketsSold += ticketCount
	}

	// Chuyển đổi sang response format
	var movieItems []response.MovieRevenueItem
	for movieID, stat := range movieStats {
		// Tính số suất chiếu thực tế (tất cả suất chiếu trong khoảng thời gian)
		actualShowtimes, err := s.showtimeRepo.GetByMovieIDAndDateRange(ctx, movieID, req.StartDate, req.EndDate)
		if err != nil {
			continue
		}
		stat.ShowtimeCount = len(actualShowtimes)

		// Tính tỷ lệ lấp đầy chính xác
		stat.OccupancyRate = s.calculateMovieOccupancyRate(ctx, movieID, req.StartDate, req.EndDate, stat.TicketsSold)

		// Tính giá trung bình
		if stat.TicketsSold > 0 {
			stat.AveragePrice = stat.TotalRevenue / float64(stat.TicketsSold)
		}

		movieItems = append(movieItems, response.MovieRevenueItem{
			MovieID:       stat.MovieID,
			MovieTitle:    stat.MovieTitle,
			TicketsSold:   stat.TicketsSold,
			AveragePrice:  stat.AveragePrice,
			ShowtimeCount: stat.ShowtimeCount,
			OccupancyRate: stat.OccupancyRate,
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

	cinemaStats := make(map[uuid.UUID]*domain.CinemaStatistic)
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

		// Tính doanh thu combo (xử lý discount nếu có)
		orderCombos, _ := s.orderComboRepo.GetByOrderID(ctx, order.ID)
		comboRevenue := float64(0)
		for _, combo := range orderCombos {
			comboRevenue += combo.TotalPrice
		}

		var ticketRevenue float64
		if order.DiscountID != nil {
			discount, err := s.discountRepo.GetDiscount(ctx, *order.DiscountID)
			if err != nil {
				return nil, fault.Wrapf(err, "[%v] failed to get discount", caller)
			}

			ticketRevenue = order.TotalPrice - comboRevenue*(1-discount.Percentage/100)
		} else {
			ticketRevenue = order.TotalPrice - comboRevenue
		}

		if _, exists := cinemaStats[room.CinemaID]; !exists {
			cinema, err := s.cinemaRepo.FindByID(ctx, room.CinemaID)
			if err != nil {
				continue
			}

			cinemaStats[room.CinemaID] = &domain.CinemaStatistic{
				CinemaID:      room.CinemaID,
				CinemaName:    cinema.Name,
				TicketRevenue: 0,
				ComboRevenue:  0,
				TotalRevenue:  0,
				TicketsSold:   0,
				ShowtimeCount: 0,
				OccupancyRate: 0,
				StartDate:     req.StartDate,
				EndDate:       req.EndDate,
			}
		}

		stat := cinemaStats[room.CinemaID]
		stat.TicketRevenue += ticketRevenue
		stat.ComboRevenue += comboRevenue // combo revenue đã bao gồm discount
		stat.TicketsSold += ticketCount

		totalTicketRevenue += ticketRevenue
		totalComboRevenue += comboRevenue
		totalRevenue += order.TotalPrice
	}

	// Chuyển đổi sang response format
	var cinemaItems []response.CinemaRevenueItem
	for cinemaID, stat := range cinemaStats {
		// Tính số suất chiếu thực tế của cinema
		rooms, err := s.roomRepo.GetListByCinemaId(ctx, cinemaID)
		if err != nil {
			continue
		}

		actualShowtimeCount := 0
		for _, room := range rooms {
			roomShowtimes, err := s.showtimeRepo.GetByRoomIDAndDateRange(ctx, room.ID, req.StartDate, req.EndDate)
			if err != nil {
				continue
			}
			actualShowtimeCount += len(roomShowtimes)
		}
		stat.ShowtimeCount = actualShowtimeCount

		// Tính tỷ lệ lấp đầy và tổng doanh thu
		stat.OccupancyRate = s.calculateCinemaOccupancyRate(ctx, cinemaID, req.StartDate, req.EndDate, stat.TicketsSold)
		stat.TotalRevenue = stat.TicketRevenue + stat.ComboRevenue

		cinemaItems = append(cinemaItems, response.CinemaRevenueItem{
			CinemaID:      stat.CinemaID,
			CinemaName:    stat.CinemaName,
			TicketRevenue: stat.TicketRevenue,
			ComboRevenue:  stat.ComboRevenue,
			TicketsSold:   stat.TicketsSold,
			ShowtimeCount: stat.ShowtimeCount,
			OccupancyRate: stat.OccupancyRate,
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

	comboStats := make(map[uuid.UUID]*domain.ComboStatistic)
	var totalRevenue float64
	var totalQuantitySold int

	for _, order := range orders {
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
				combo, err := s.comboRepo.FindByID(ctx, orderCombo.ComboID)
				if err != nil {
					continue
				}

				comboStats[orderCombo.ComboID] = &domain.ComboStatistic{
					ComboID:           orderCombo.ComboID,
					ComboName:         combo.Name,
					ComboDescription:  combo.Description,
					ComboPrice:        combo.Price,
					QuantitySold:      0,
					TotalRevenue:      0,
					PercentageOfTotal: 0,
					StartDate:         req.StartDate,
					EndDate:           req.EndDate,
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
	for _, stat := range comboStats {
		// Tính phần trăm tổng doanh thu
		if totalRevenue > 0 {
			stat.PercentageOfTotal = (stat.TotalRevenue / totalRevenue) * 100
		}

		comboItems = append(comboItems, response.ComboStatisticItem{
			ID:                stat.ComboID,
			Name:              stat.ComboName,
			Description:       stat.ComboDescription,
			Price:             stat.ComboPrice,
			Quantity:          stat.QuantitySold,
			Revenue:           stat.TotalRevenue,
			PercentageOfTotal: stat.PercentageOfTotal,
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
	if err != nil || len(showtimes) == 0 {
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

	occupancyRate := float64(ticketsSold) / float64(totalCapacity) * 100

	// Đảm bảo tỷ lệ không vượt quá 100% (có thể do lỗi dữ liệu)
	if occupancyRate > 100 {
		occupancyRate = 100
	}

	return occupancyRate
}

func (s *StatisticService) calculateCinemaOccupancyRate(ctx context.Context, cinemaID uuid.UUID, startDate, endDate time.Time, ticketsSold int) float64 {
	// Lấy tất cả rooms của cinema
	rooms, err := s.roomRepo.GetListByCinemaId(ctx, cinemaID)
	if err != nil || len(rooms) == 0 {
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

	occupancyRate := float64(ticketsSold) / float64(totalCapacity) * 100

	// Đảm bảo tỷ lệ không vượt quá 100% (có thể do lỗi dữ liệu)
	if occupancyRate > 100 {
		occupancyRate = 100
	}

	return occupancyRate
}
