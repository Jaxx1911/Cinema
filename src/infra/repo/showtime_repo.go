package repo

import (
	"TTCS/src/core/domain"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ShowtimeRepo struct {
	*BaseRepo
	db *gorm.DB
}

func NewShowtimeRepo(baseRepo *BaseRepo, db *gorm.DB) domain.ShowtimeRepo {
	return &ShowtimeRepo{
		BaseRepo: baseRepo,
		db:       db,
	}
}

func (s *ShowtimeRepo) Create(ctx context.Context, showtime *domain.Showtime) (*domain.Showtime, error) {
	if err := s.db.WithContext(ctx).Create(showtime).Error; err != nil {
		return nil, s.returnError(ctx, err)
	}
	return showtime, nil
}

func (s *ShowtimeRepo) GetListByFilter(ctx context.Context, movieId uuid.UUID, cinemaId uuid.UUID, day time.Time) ([]*domain.Showtime, error) {
	var showtimes []*domain.Showtime

	if err := s.db.WithContext(ctx).
		Preload("Room").
		Joins("JOIN room ON room.id = showtime.room_id").
		Where("showtime.movie_id = ? AND room.cinema_id = ? AND DATE(showtime.start_time) >= ? AND DATE(showtime.start_time) < ?", movieId, cinemaId, day, day.Add(24*time.Hour)).
		Order("start_time asc").
		Find(&showtimes).Error; err != nil {
		return nil, s.returnError(ctx, err)
	}
	return showtimes, nil
}

func (s *ShowtimeRepo) FindConflictByRoomId(ctx context.Context, roomId uuid.UUID, startTime, endTime time.Time) ([]domain.Showtime, error) {
	var conflictingShowtimes []domain.Showtime
	if err := s.db.WithContext(ctx).Where("room_id = ? AND start_time < ? AND end_time > ?", roomId, endTime, startTime).
		Find(&conflictingShowtimes).Error; err != nil {
		return nil, s.returnError(ctx, err)
	}
	return conflictingShowtimes, nil
}

func (s *ShowtimeRepo) GetListByCinemaFilter(ctx context.Context, id uuid.UUID, day time.Time) ([]*domain.Showtime, error) {
	var showtimes []*domain.Showtime

	if err := s.db.WithContext(ctx).
		Preload("Room").
		Joins("JOIN room ON room.id = showtime.room_id").
		Where("room.cinema_id = ? AND DATE(showtime.start_time) >= ? AND DATE(showtime.start_time) < ?", id, day, day.Add(24*time.Hour)).
		Order("start_time asc").Find(&showtimes).Error; err != nil {
		return nil, s.returnError(ctx, err)
	}
	return showtimes, nil
}

func (s *ShowtimeRepo) GetListByRoomFilter(ctx context.Context, id uuid.UUID, day time.Time) ([]*domain.Showtime, error) {
	var showtimes []*domain.Showtime

	if err := s.db.WithContext(ctx).
		Where("room_id = ? AND DATE(showtime.start_time) >= ? AND DATE(showtime.start_time) < ?", id, day, day.Add(24*time.Hour)).
		Order("start_time asc").Find(&showtimes).Error; err != nil {
		return nil, s.returnError(ctx, err)
	}
	return showtimes, nil
}

func (s *ShowtimeRepo) GetById(ctx context.Context, id uuid.UUID) (*domain.Showtime, error) {
	var showtime domain.Showtime
	if err := s.db.WithContext(ctx).Preload("Tickets").Preload("Room").Where("id = ?", id).First(&showtime).Error; err != nil {
		return nil, s.returnError(ctx, err)
	}
	return &showtime, nil
}
