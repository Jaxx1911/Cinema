package repo

import (
	"TTCS/src/common/fault"
	"TTCS/src/core/domain"
	"TTCS/src/infra/cache"
	"TTCS/src/present/httpui/request"
	"context"
	uuid2 "github.com/google/uuid"
	"gorm.io/gorm"
)

type MovieRepo struct {
	*BaseRepo
	db    *gorm.DB
	redis *cache.RedisCache
}

func NewMovieRepo(baseRepo *BaseRepo, db *gorm.DB, redis *cache.RedisCache) domain.MovieRepo {
	return &MovieRepo{
		BaseRepo: baseRepo,
		db:       db,
		redis:    redis,
	}
}

func (m MovieRepo) GetList(ctx context.Context, page request.Page, showingStatus string) ([]*domain.Movie, error) {
	var movies []*domain.Movie
	limit, offset := m.toLimitOffset(ctx, page)
	if err := m.db.Where("status = ?", showingStatus).Limit(limit).Offset(offset).Order("release_date").Find(movies).Error; err != nil {
		return nil, m.returnError(ctx, err)
	}
	return movies, nil
}

func (m MovieRepo) Create(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
	if err := m.db.Create(movie).Error; err != nil {
		return nil, m.returnError(ctx, err)
	}
	return movie, nil
}

func (m MovieRepo) Update(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
	if err := m.db.Save(movie).Error; err != nil {
		return nil, m.returnError(ctx, err)
	}
	return movie, nil
}

func (m MovieRepo) GetDetail(ctx context.Context, id string) (*domain.Movie, error) {
	var movie domain.Movie
	uuid, err := uuid2.Parse(id)
	if err != nil {
		return nil, fault.Wrapf(err, "invalid uuid").SetTag(fault.TagBadRequest)
	}
	if err = m.db.Preload("Showtimes").Preload("Genres").Where("id = ?", uuid).First(&movie).Error; err != nil {
		return nil, m.returnError(ctx, err)
	}
	return &movie, nil
}

func (m MovieRepo) GetById(ctx context.Context, id string) (*domain.Movie, error) {
	var movie domain.Movie
	uuid, err := uuid2.Parse(id)
	if err != nil {
		return nil, fault.Wrapf(err, "invalid uuid").SetTag(fault.TagBadRequest)
	}
	if err = m.db.Where("id = ?", uuid).First(&movie).Error; err != nil {
		return nil, m.returnError(ctx, err)
	}
	return &movie, nil
}
