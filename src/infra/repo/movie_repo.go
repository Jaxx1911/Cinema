package repo

import (
	"TTCS/src/core/domain"
	"TTCS/src/present/httpui/request"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type MovieRepo struct {
	*BaseRepo
}

func NewMovieRepo(baseRepo *BaseRepo) domain.MovieRepo {
	return &MovieRepo{
		BaseRepo: baseRepo,
	}
}

func (m MovieRepo) GetList(ctx context.Context, req request.GetListMovie) ([]*domain.Movie, int64, error) {
	var movies []*domain.Movie
	var totalCount int64

	// Lấy limit và offset từ page
	limit, offset := m.toLimitOffset(ctx, req.Page)
	query := m.db.Model(&domain.Movie{})
	if req.Tag != "all" {
		query = query.Where("tag = ?", req.Tag)
	}
	if req.Status != "all" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Term != "all" {
		query = query.Where("LOWER(title) LIKE ?", fmt.Sprintf("%%%s%%", strings.ToLower(req.Term)))
	}
	// Lấy tổng số phim
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, m.returnError(ctx, err)
	}

	// Lấy danh sách phim với limit và offset
	if err := query.Preload("Genres").Limit(limit).Offset(offset).Order("release_date").Find(&movies).Error; err != nil {
		return nil, 0, m.returnError(ctx, err)
	}

	return movies, totalCount, nil
}

func (m MovieRepo) GetListByStatus(ctx context.Context, page request.Page, showingStatus string) ([]*domain.Movie, error) {
	var movies []*domain.Movie
	limit, offset := m.toLimitOffset(ctx, page)
	if err := m.db.Preload("Genres").Where("status = ?", showingStatus).Limit(limit).Offset(offset).Order("release_date").Find(&movies).Error; err != nil {
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
	if err := m.db.Updates(movie).Error; err != nil {
		return nil, m.returnError(ctx, err)
	}
	if err := m.db.Model(movie).Association("Genres").Replace(movie.Genres); err != nil {
		return nil, m.returnError(ctx, err)
	}
	return movie, nil
}

func (m MovieRepo) GetDetail(ctx context.Context, id uuid.UUID) (*domain.Movie, error) {
	movie := &domain.Movie{}
	if err := m.db.Preload("Showtimes").Preload("Genres").Where("id = ?", id).First(movie).Error; err != nil {
		return nil, m.returnError(ctx, err)
	}
	return movie, nil
}

func (m MovieRepo) GetById(ctx context.Context, id uuid.UUID) (*domain.Movie, error) {
	movie := &domain.Movie{}
	if err := m.db.Preload("Genres").Where("id = ?", id).First(movie).Error; err != nil {
		return nil, m.returnError(ctx, err)
	}
	return movie, nil
}

func (m MovieRepo) GetListInDateRange(ctx context.Context, startDate time.Time, endDate time.Time) ([]*domain.Movie, error) {
	var movies []*domain.Movie
	if err := m.db.Preload("Genres").Where("status != ? AND release_date <= ?", "off", endDate).Order("title").Find(&movies).Error; err != nil {
		return nil, m.returnError(ctx, err)
	}
	return movies, nil
}

func (m MovieRepo) GetMoviesByReleaseDateAndStatus(ctx context.Context, releaseDate time.Time, status string) ([]*domain.Movie, error) {
	var movies []*domain.Movie
	// Query movies where release_date is today (same date) and status matches
	startOfDay := time.Date(releaseDate.Year(), releaseDate.Month(), releaseDate.Day(), 0, 0, 0, 0, releaseDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	if err := m.db.Preload("Genres").Where("status = ? AND release_date >= ? AND release_date < ?", status, startOfDay, endOfDay).Find(&movies).Error; err != nil {
		return nil, m.returnError(ctx, err)
	}
	return movies, nil
}
