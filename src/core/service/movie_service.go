package service

import (
	"TTCS/src/common/fault"
	"TTCS/src/core/domain"
	"TTCS/src/infra/upload"
	"TTCS/src/present/httpui/request"
	"context"
	"github.com/google/uuid"
	"mime/multipart"
	"strings"
	"time"
)

type MovieService struct {
	movieRepo domain.MovieRepo
	genreRepo domain.GenreRepo
	upload    *upload.UploadService
}

func NewMovieService(movieRepo domain.MovieRepo, genreRepo domain.GenreRepo, upload *upload.UploadService) *MovieService {
	return &MovieService{
		movieRepo: movieRepo,
		genreRepo: genreRepo,
		upload:    upload,
	}
}

func (m *MovieService) GetList(ctx context.Context, page request.Page) ([]*domain.Movie, int64, error) {
	_ = "MovieService.GetList"
	movies, total, err := m.movieRepo.GetList(ctx, page)
	if err != nil {
		return nil, 0, err
	}
	return movies, total, nil
}

func (m *MovieService) GetListByStatus(ctx context.Context, page request.Page, status string) ([]*domain.Movie, error) {
	_ = "MovieService.GetListByStatus"
	movies, err := m.movieRepo.GetListByStatus(ctx, page, status)
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (m *MovieService) GetDetail(ctx context.Context, id uuid.UUID) (*domain.Movie, error) {
	_ = "MovieService.GetDetail"
	movie, err := m.movieRepo.GetDetail(ctx, id)
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (m *MovieService) Create(ctx context.Context, req request.CreateMovieRequest) (*domain.Movie, error) {
	caller := "MovieService.Create"

	pUrl, err := m.upload.UploadFile(ctx, req.PosterImage)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to upload movie poster", caller)
	}

	lUrl, err := m.upload.UploadFile(ctx, req.LargePosterImage)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to upload movie large poster", caller)
	}

	releaseDate, err := time.Parse("02-01-2006", req.ReleaseDate)

	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to parse release date", caller).SetTag(fault.TagBadRequest).SetKey(fault.KeyMovie)
	}
	ids := make([]uuid.UUID, 0)
	for _, g := range req.Genres {
		ids = append(ids, uuid.MustParse(g))
	}
	listGenre, err := m.genreRepo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	movie, err := m.movieRepo.Create(ctx, &domain.Movie{
		Title:          req.Title,
		Duration:       req.Duration,
		PosterURL:      pUrl,
		LargePosterURL: lUrl,
		Director:       strings.Join(req.Director, ","),
		Caster:         strings.Join(req.Caster, ","),
		Description:    req.Description,
		ReleaseDate:    releaseDate,
		TrailerURL:     req.TrailerURL,
		Genres:         listGenre,
		Tag:            req.Tag,
	})
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (m *MovieService) Update(ctx context.Context, req request.UpdateMovieRequest) (*domain.Movie, error) {
	caller := "MovieService.Update"

	movie, err := m.movieRepo.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	releaseDate, err := time.Parse("02-01-2006", req.ReleaseDate)

	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to parse release date", caller).SetTag(fault.TagBadRequest).SetKey(fault.KeyMovie)
	}
	ids := make([]uuid.UUID, 0)
	for _, g := range req.Genres {
		ids = append(ids, uuid.MustParse(g))
	}
	listGenre, err := m.genreRepo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	movie.Title = req.Title
	movie.Duration = req.Duration
	movie.Description = req.Description
	movie.ReleaseDate = releaseDate
	movie.TrailerURL = req.TrailerURL
	movie.Genres = listGenre
	movie.Director = strings.Join(req.Director, ",")
	movie.Caster = strings.Join(req.Caster, ",")
	movie.Status = req.Status

	movie, err = m.movieRepo.Update(ctx, movie)
	if req.PosterImage != nil {
		movie, err = m.UpdatePoster(ctx, movie, req.PosterImage)
	}
	if err != nil {
		return nil, err
	}
	if req.LargePosterImage != nil {
		movie, err = m.UpdateLargePoster(ctx, movie, req.LargePosterImage)
	}

	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (m *MovieService) UpdatePoster(ctx context.Context, movie *domain.Movie, posterImage *multipart.FileHeader) (*domain.Movie, error) {
	caller := "MovieService.UpdatePoster"
	url, err := m.upload.UploadFile(ctx, posterImage)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to upload movie poster", caller)
	}
	movie.PosterURL = url
	movie, err = m.movieRepo.Update(ctx, movie)
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (m *MovieService) UpdateLargePoster(ctx context.Context, movie *domain.Movie, posterImage *multipart.FileHeader) (*domain.Movie, error) {
	caller := "MovieService.UpdatePoster"
	url, err := m.upload.UploadFile(ctx, posterImage)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to upload movie poster", caller)
	}
	movie.LargePosterURL = url
	movie, err = m.movieRepo.Update(ctx, movie)
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (m *MovieService) GetListInDateRange(ctx context.Context) ([]*domain.Movie, error) {
	dayStart := time.Now()
	dayEnd := time.Now().AddDate(0, 0, 5)
	movies, err := m.movieRepo.GetListInDateRange(ctx, dayStart, dayEnd)
	if err != nil {
		return nil, err
	}
	return movies, nil
}
