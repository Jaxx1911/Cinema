package service

import (
	"TTCS/src/common/log"
	"TTCS/src/core/domain"
	"TTCS/src/infra/upload"
	"TTCS/src/present/httpui/request"
	"context"
	"mime/multipart"

	"github.com/google/uuid"
)

type ComboService interface {
	GetList(ctx context.Context) ([]*domain.Combo, error)
	GetDetail(ctx context.Context, id uuid.UUID) (*domain.Combo, error)
	Create(ctx context.Context, req request.CreateComboRequest) (*domain.Combo, error)
	Update(ctx context.Context, id uuid.UUID, req request.UpdateComboRequest) (*domain.Combo, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type comboService struct {
	comboRepo domain.ComboRepository
	upload    *upload.UploadService
}

func NewComboService(comboRepo domain.ComboRepository, upload *upload.UploadService) ComboService {
	return &comboService{
		comboRepo: comboRepo,
		upload:    upload,
	}
}

func (s *comboService) GetList(ctx context.Context) ([]*domain.Combo, error) {
	combos, err := s.comboRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return combos, nil
}

func (s *comboService) GetDetail(ctx context.Context, id uuid.UUID) (*domain.Combo, error) {
	combo, err := s.comboRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return combo, nil
}

func (s *comboService) Create(ctx context.Context, req request.CreateComboRequest) (*domain.Combo, error) {
	url, err := s.uploadPoster(ctx, req.BannerUrl)
	if err != nil {
		return nil, err
	}

	combo := &domain.Combo{
		Name:        req.Name,
		Description: req.Description,
		BannerUrl:   url,
		Price:       req.Price,
	}
	err = s.comboRepo.Create(ctx, combo)
	if err != nil {
		return nil, err
	}
	return combo, nil
}

func (s *comboService) Update(ctx context.Context, id uuid.UUID, req request.UpdateComboRequest) (*domain.Combo, error) {
	combo, err := s.comboRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.BannerUrl != nil {
		url, err := s.uploadPoster(ctx, req.BannerUrl)
		if err != nil {
			return nil, err
		}
		combo.BannerUrl = url
	}

	combo.Name = req.Name
	combo.Description = req.Description

	combo.Price = req.Price

	err = s.comboRepo.Update(ctx, combo)
	if err != nil {
		return nil, err
	}
	return combo, nil
}

func (s *comboService) Delete(ctx context.Context, id uuid.UUID) error {
	caller := "ComboService.Delete"
	_, err := s.comboRepo.FindByID(ctx, id)
	if err != nil {
		log.Error(ctx, "[%v] failed to find combo: %v", caller, err)
		return err
	}
	err = s.comboRepo.Delete(ctx, id)
	if err != nil {
		log.Error(ctx, "[%v] failed to delete combo: %v", caller, err)
		return err
	}
	return nil
}

func (s *comboService) uploadPoster(ctx context.Context, file *multipart.FileHeader) (string, error) {
	caller := "ComboService.uploadPoster"
	url, err := s.upload.UploadFile(ctx, file)
	if err != nil {
		log.Error(ctx, "[%v] failed to upload poster: %v", caller, err)
		return "", err
	}
	return url, nil
}
