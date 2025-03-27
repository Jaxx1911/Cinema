package service

import (
	"TTCS/src/common/fault"
	"TTCS/src/core/domain"
	"TTCS/src/present/httpui/request"
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
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

	conflictTime, err := s.ShowtimeRepo.FindConflictByRoomId(ctx, room.ID, startTime.Add(-30*time.Minute), s.roundToNextHour(endTime))
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
	var roundedTime time.Time
	if t.Minute() <= 30 {
		roundedTime = t.Add(time.Hour/2 - time.Duration(t.Minute())*time.Minute)
	} else {
		roundedTime = t.Add(time.Hour - time.Duration(t.Minute())*time.Minute)
	}
	return roundedTime.Add(30 * time.Minute)
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
