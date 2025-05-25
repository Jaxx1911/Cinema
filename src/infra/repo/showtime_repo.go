package repo

import (
	"TTCS/src/core/domain"
	"TTCS/src/present/httpui/request"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	if err := s.db.WithContext(ctx).Preload("Movie").Where("room_id = ? AND start_time < ? AND end_time > ?", roomId, endTime, startTime).
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
		Preload("Room").
		Preload("Movie").
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

func (r *ShowtimeRepo) GetList(ctx context.Context, page request.GetListShowtime) ([]*domain.Showtime, int64, error) {
	var showtimes []*domain.Showtime
	var total int64

	query := r.db.Model(&domain.Showtime{}).
		Preload("Movie").
		Preload("Room").
		Joins("JOIN room ON room.id = showtime.room_id")

	// Apply filters
	if page.MovieID != "" {
		query = query.Where("movie_id = ?", page.MovieID)
	}
	if page.RoomID != "" {
		query = query.Where("room_id = ?", page.RoomID)
	}
	if page.CinemaID != "" {
		query = query.Where("room.cinema_id = ?", page.CinemaID)
	}
	if page.FromDate != "" {
		fromDate, err := time.Parse("2006-01-02", page.FromDate)
		if err == nil {
			query = query.Where("start_time >= ?", fromDate)
		}
	}
	if page.ToDate != "" {
		toDate, err := time.Parse("2006-01-02", page.ToDate)
		if err == nil {
			// Set to end of day
			toDate = toDate.Add(24 * time.Hour).Add(-time.Second)
			query = query.Where("start_time <= ?", toDate)
		}
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, r.returnError(ctx, err)
	}

	// Apply pagination
	limit, offset := r.toLimitOffset(ctx, page.Page)
	if err := query.WithContext(ctx).Offset(offset).Limit(limit).Find(&showtimes).Error; err != nil {
		return nil, 0, r.returnError(ctx, err)
	}

	return showtimes, total, nil
}

func (r *ShowtimeRepo) Update(ctx context.Context, id uuid.UUID, showtime *domain.Showtime) (*domain.Showtime, error) {
	if err := r.db.WithContext(ctx).Model(&domain.Showtime{}).Where("id = ?", id).Updates(showtime).Error; err != nil {
		return nil, r.returnError(ctx, err)
	}
	return r.GetById(ctx, id)
}

func (r *ShowtimeRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Showtime{}, id).Error; err != nil {
		return r.returnError(ctx, err)
	}
	return nil
}
