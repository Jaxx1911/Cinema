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

func (s *ShowtimeRepo) GetListByMovieId(ctx context.Context, movieId string) ([]*domain.Showtime, error) {
	//uuid, err := uuid.Parse(movieId)
	//if err != nil {
	//	return nil, s.returnError(ctx, err)
	//}
	//var showtimes []*domain.Showtime
	//if err := s.db.WithContext(ctx).Where("movie_id = ?", uuid).Find(&showtimes).Error; err != nil {
	//}
	return nil, nil
}

func (s *ShowtimeRepo) FindConflictByRoomId(ctx context.Context, roomId uuid.UUID, startTime, endTime time.Time) ([]domain.Showtime, error) {
	var conflictingShowtimes []domain.Showtime
	if err := s.db.Where("room_id = ? AND start_time < ? AND end_time > ?", roomId, endTime, startTime).
		Find(&conflictingShowtimes).Error; err != nil {
		return nil, s.returnError(ctx, err)
	}
	return conflictingShowtimes, nil
}
