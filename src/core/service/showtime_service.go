package service

import (
	"TTCS/src/common/fault"
	"TTCS/src/core/domain"
	"TTCS/src/core/dto"
	"TTCS/src/present/httpui/request"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type ShowtimeService struct {
	ShowtimeRepo domain.ShowtimeRepo
	MovieRepo    domain.MovieRepo
	RoomRepo     domain.RoomRepo
	TicketRepo   domain.TicketRepo
}

func NewShowtimeService(showtimeRepo domain.ShowtimeRepo, movieRepo domain.MovieRepo, roomRepo domain.RoomRepo, ticketRepo domain.TicketRepo) *ShowtimeService {
	return &ShowtimeService{
		ShowtimeRepo: showtimeRepo,
		MovieRepo:    movieRepo,
		RoomRepo:     roomRepo,
		TicketRepo:   ticketRepo,
	}
}

func (s *ShowtimeService) Create(ctx context.Context, req request.CreateShowtime) (*domain.Showtime, error) {
	caller := "ShowtimeService.Create"

	movie, err := s.MovieRepo.GetById(ctx, req.MovieId)
	if err != nil {
		return nil, err
	}
	room, err := s.RoomRepo.GetById(ctx, req.RoomId)
	if err != nil {
		return nil, err
	}

	startTime, err := time.ParseInLocation("02-01-2006 15:04", req.StartTime, time.FixedZone("UTC+7", 7*60*60))
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to parse start time", caller).SetTag(fault.TagBadRequest).SetKey(fault.KeyShowtime)
	}
	duration := time.Duration(movie.Duration) * time.Minute
	endTime := startTime.Add(duration)

	conflictTime, err := s.ShowtimeRepo.FindConflictByRoomId(ctx, room.ID, startTime.Add(-15*time.Minute), s.roundToNextHour(endTime))
	if len(conflictTime) > 0 {
		err := errors.New("conflict time")
		return nil, fault.Wrapf(err, "[%v] conflict time", caller).SetTag(fault.TagBadRequest).SetKey(fault.KeyShowtime)
	}

	showtime, err := s.ShowtimeRepo.Create(ctx, &domain.Showtime{
		MovieID:   movie.ID,
		RoomID:    room.ID,
		StartTime: startTime,
		EndTime:   endTime,
		Price:     req.Price,
	})

	if err != nil {
		return nil, err
	}

	var tickets []*domain.Ticket
	for _, seat := range room.Seats {
		tickets = append(tickets, &domain.Ticket{
			ShowtimeID: showtime.ID,
			SeatID:     seat.ID,
			Status:     "available",
			Showtime:   domain.Showtime{},
			Seat:       domain.Seat{},
		})
	}
	tickets, err = s.TicketRepo.Create(ctx, tickets)
	if err != nil {
		return nil, err
	}
	return showtime, nil
}

func (s *ShowtimeService) roundToNextHour(t time.Time) time.Time {
	m := t.Minute()
	sec := t.Second()
	nano := t.Nanosecond()

	// Reset về đầu phút để tính toán dễ hơn
	t = t.Add(-time.Duration(m)*time.Minute - time.Duration(sec)*time.Second - time.Duration(nano)*time.Nanosecond)

	switch {
	case m <= 15:
		return t.Add(30 * time.Minute)
	case m <= 30:
		return t.Add(45 * time.Minute)
	case m <= 45:
		return t.Add(1 * time.Hour)
	default:
		return t.Add(1*time.Hour + 15*time.Minute)
	}
}

func (s *ShowtimeService) GetByUserFilter(ctx context.Context, filter request.GetShowtimesByUserFilter) ([]*domain.Showtime, error) {
	caller := "ShowtimeService.GetByUserFilter"
	day, err := time.ParseInLocation("02-01-2006", filter.Day, time.FixedZone("UTC+7", 7*60*60))
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to parse day", caller).SetTag(fault.TagBadRequest).SetKey(fault.KeyShowtime)
	}

	movieId, err := uuid.Parse(filter.MovieId)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] invalid uuid", caller).SetTag(fault.TagBadRequest).SetKey(fault.KeyMovie)
	}
	cinemaId, err := uuid.Parse(filter.CinemaId)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] invalid uuid", caller).SetTag(fault.TagBadRequest).SetKey(fault.KeyCinema)
	}

	showtimes, err := s.ShowtimeRepo.GetListByFilter(ctx, movieId, cinemaId, day)
	if err != nil {
		return nil, err
	}
	return showtimes, nil
}

func (s *ShowtimeService) GetByCinemaFilter(ctx context.Context, filter request.GetShowtimesByCinemaIdFilter) ([]*domain.Showtime, error) {
	caller := "ShowtimeService.GetByCinemaFilter"
	day, err := time.ParseInLocation("02-01-2006", filter.Day, time.FixedZone("UTC+7", 7*60*60))
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to parse day", caller).SetTag(fault.TagBadRequest).SetKey(fault.KeyShowtime)
	}

	cinemaId, err := uuid.Parse(filter.CinemaId)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] invalid uuid", caller).SetTag(fault.TagBadRequest).SetKey(fault.KeyCinema)
	}

	showtimes, err := s.ShowtimeRepo.GetListByCinemaFilter(ctx, cinemaId, day)
	if err != nil {
		return nil, err
	}
	return showtimes, nil
}

func (s *ShowtimeService) GetByRoomFilter(ctx context.Context, filter request.GetShowtimesByRoomIdFilter) ([]*domain.Showtime, error) {
	caller := "ShowtimeService.GetByCinemaFilter"
	day, err := time.ParseInLocation("02-01-2006", filter.Day, time.FixedZone("UTC+7", 7*60*60))
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to parse day", caller).SetTag(fault.TagBadRequest).SetKey(fault.KeyShowtime)
	}

	showtimes, err := s.ShowtimeRepo.GetListByRoomFilter(ctx, uuid.MustParse(filter.RoomId), day)
	if err != nil {
		return nil, err
	}
	return showtimes, nil
}

func (s *ShowtimeService) GetById(ctx context.Context, id string) (*domain.Showtime, error) {
	caller := "ShowtimeService.GetById"
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] invalid uuid", caller).SetTag(fault.TagBadRequest).SetKey(fault.KeyShowtime)
	}
	showtime, err := s.ShowtimeRepo.GetById(ctx, uid)
	if err != nil {
		return nil, err
	}
	return showtime, nil
}

func (s *ShowtimeService) GetList(ctx context.Context, page request.GetListShowtime) ([]*domain.Showtime, int64, error) {
	caller := "ShowtimeService.GetList"
	showtimes, total, err := s.ShowtimeRepo.GetList(ctx, page)
	if err != nil {
		return nil, 0, fault.Wrapf(err, "[%v] failed to get showtimes", caller)
	}

	return showtimes, total, nil
}

func (s *ShowtimeService) Update(ctx context.Context, id string, req request.UpdateShowtime) (*domain.Showtime, error) {
	caller := "ShowtimeService.Update"

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] invalid uuid", caller).SetTag(fault.TagBadRequest).SetKey(fault.KeyShowtime)
	}

	// Get existing showtime
	existingShowtime, err := s.ShowtimeRepo.GetById(ctx, uid)
	if err != nil {
		return nil, err
	}

	// Validate movie and room
	movie, err := s.MovieRepo.GetById(ctx, req.MovieId)
	if err != nil {
		return nil, err
	}
	room, err := s.RoomRepo.GetById(ctx, req.RoomId)
	if err != nil {
		return nil, err
	}

	// Parse and validate start time
	startTime, err := time.ParseInLocation("02-01-2006 15:04", req.StartTime, time.FixedZone("UTC+7", 7*60*60))
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to parse start time", caller).SetTag(fault.TagBadRequest).SetKey(fault.KeyShowtime)
	}

	// Calculate end time
	duration := time.Duration(movie.Duration) * time.Minute
	endTime := startTime.Add(duration)

	// Check for time conflicts
	conflictTime, err := s.ShowtimeRepo.FindConflictByRoomId(ctx, room.ID, startTime.Add(-15*time.Minute), s.roundToNextHour(endTime))
	if len(conflictTime) > 0 {
		// Skip conflict check if it's the same showtime
		if conflictTime[0].ID != uid {
			err := errors.New("conflict time")
			return nil, fault.Wrapf(err, "[%v] conflict time", caller).SetTag(fault.TagBadRequest).SetKey(fault.KeyShowtime)
		}
	}

	// Update showtime
	existingShowtime.MovieID = movie.ID
	existingShowtime.RoomID = room.ID
	existingShowtime.StartTime = startTime
	existingShowtime.EndTime = endTime
	existingShowtime.Price = req.Price

	return s.ShowtimeRepo.Update(ctx, uid, existingShowtime)
}

func (s *ShowtimeService) Delete(ctx context.Context, id string) error {
	caller := "ShowtimeService.Delete"

	uid, err := uuid.Parse(id)
	if err != nil {
		return fault.Wrapf(err, "[%v] invalid uuid", caller).SetTag(fault.TagBadRequest).SetKey(fault.KeyShowtime)
	}

	// Check if showtime exists
	_, err = s.ShowtimeRepo.GetById(ctx, uid)
	if err != nil {
		return err
	}

	return s.ShowtimeRepo.Delete(ctx, uid)
}

func (s *ShowtimeService) CheckShowtimeAvailability(ctx context.Context, req request.CheckShowtimeAvailability) (*dto.ShowtimeAvailabilityResponse, error) {
	caller := "ShowtimeService.CheckShowtimeAvailability"

	movie, err := s.MovieRepo.GetById(ctx, req.MovieId)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] movie not found", caller).SetTag(fault.TagBadRequest).SetKey(fault.KeyMovie)
	}

	room, err := s.RoomRepo.GetById(ctx, req.RoomId)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] room not found", caller).SetTag(fault.TagBadRequest).SetKey(fault.KeyRoom)
	}

	startTime, err := time.ParseInLocation("02-01-2006 15:04", req.StartTime, time.FixedZone("UTC+7", 7*60*60))
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to parse start time", caller).SetTag(fault.TagBadRequest).SetKey(fault.KeyShowtime)
	}

	duration := time.Duration(movie.Duration) * time.Minute
	endTime := startTime.Add(duration)

	if req.Id != nil {
		conflictShowtimes, err := s.ShowtimeRepo.FindConflictToUpdate(ctx, room.ID, startTime.Add(-15*time.Minute), s.roundToNextHour(endTime), *req.Id)
		if err != nil {
			return nil, fault.Wrapf(err, "[%v] failed to check conflicts", caller)
		}
		response := &dto.ShowtimeAvailabilityResponse{
			IsAvailable: len(conflictShowtimes) == 0,
			Conflicts:   conflictShowtimes,
		}
		return response, nil
	}

	conflictShowtimes, err := s.ShowtimeRepo.FindConflictByRoomId(ctx, room.ID, startTime.Add(-15*time.Minute), s.roundToNextHour(endTime))
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to check conflicts", caller)
	}

	response := &dto.ShowtimeAvailabilityResponse{
		IsAvailable: len(conflictShowtimes) == 0,
		Conflicts:   conflictShowtimes,
	}

	return response, nil
}

func (s *ShowtimeService) CheckShowtimesAvailability(ctx context.Context, req request.CheckShowtimesAvailability) (*dto.ShowtimesAvailabilityResponse, error) {
	results := make([]dto.ShowtimeAvailabilityResult, 0, len(req.Showtimes))

	for _, showtime := range req.Showtimes {
		availability, err := s.CheckShowtimeAvailability(ctx, showtime)
		if err != nil {
			return nil, err
		}

		results = append(results, dto.ShowtimeAvailabilityResult{
			ShowtimeAvailabilityResponse: *availability,
			MovieId:                      showtime.MovieId.String(),
			RoomId:                       showtime.RoomId.String(),
			StartTime:                    showtime.StartTime,
		})
	}

	return &dto.ShowtimesAvailabilityResponse{
		Results: results,
	}, nil
}

func (s *ShowtimeService) CreateShowtimes(ctx context.Context, req request.CreateShowtimes) (*dto.CreateShowtimesResponse, error) {
	results := make([]dto.CreateShowtimeResult, 0, len(req.Showtimes))
	successCount := 0
	failedCount := 0

	for _, showtimeReq := range req.Showtimes {
		result := dto.CreateShowtimeResult{
			MovieId:   showtimeReq.MovieId.String(),
			RoomId:    showtimeReq.RoomId.String(),
			StartTime: showtimeReq.StartTime,
		}

		showtime, err := s.Create(ctx, showtimeReq)
		if err != nil {
			result.Success = false
			result.Error = err.Error()
			failedCount++
		} else {
			result.Success = true
			result.Showtime = showtime
			successCount++
		}

		results = append(results, result)
	}

	return &dto.CreateShowtimesResponse{
		Results: results,
		Summary: struct {
			Total   int `json:"total"`
			Success int `json:"success"`
			Failed  int `json:"failed"`
		}{
			Total:   len(req.Showtimes),
			Success: successCount,
			Failed:  failedCount,
		},
	}, nil
}
