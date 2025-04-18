package service

import (
	"TTCS/src/common/fault"
	"TTCS/src/core/domain"
	"TTCS/src/infra/upload"
	"TTCS/src/present/httpui/request"
	"context"
	"mime/multipart"
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

func (m *MovieService) GetList(ctx context.Context, page request.Page, status string) ([]*domain.Movie, error) {
	_ = "MovieService.GetList"
	movies, err := m.movieRepo.GetList(ctx, page, status)
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (m *MovieService) GetDetail(ctx context.Context, id string) (*domain.Movie, error) {
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

	listGenre, err := m.genreRepo.GetByIDs(ctx, req.Genres)
	if err != nil {
		return nil, err
	}

	movie, err := m.movieRepo.Create(ctx, &domain.Movie{
		Title:          req.Title,
		Duration:       req.Duration,
		PosterURL:      pUrl,
		LargePosterURL: lUrl,
		Director:       req.Director,
		Caster:         req.Caster,
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

	listGenre, err := m.genreRepo.GetByIDs(ctx, req.Genres)
	if err != nil {
		return nil, err
	}

	movie.Title = req.Title
	movie.Duration = req.Duration
	movie.Description = req.Description
	movie.ReleaseDate = releaseDate
	movie.TrailerURL = req.TrailerURL
	movie.Genres = listGenre
	movie.Director = req.Director
	movie.Caster = req.Caster
	movie.Status = req.Status

	movie, err = m.movieRepo.Update(ctx, movie)
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (m *MovieService) UpdatePoster(ctx context.Context, id string, posterImage *multipart.FileHeader) (*domain.Movie, error) {
	caller := "MovieService.UpdatePoster"
	movie, err := m.movieRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
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

func (m *MovieService) GetListInDateRange(ctx context.Context) ([]*domain.Movie, error) {
	dayStart := time.Now()
	dayEnd := time.Now().AddDate(0, 0, 5)
	movies, err := m.movieRepo.GetListInDateRange(ctx, dayStart, dayEnd)
	if err != nil {
		return nil, err
	}
	return movies, nil
}
