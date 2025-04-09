package service

import (
	"TTCS/src/core/domain"
	"context"
)

type ComboService interface {
	GetList(ctx context.Context) ([]*domain.Combo, error)
}

type comboService struct {
	comboRepo domain.ComboRepository
}

func NewComboService(comboRepo domain.ComboRepository) ComboService {
	return &comboService{
		comboRepo: comboRepo,
	}
}

func (s *comboService) GetList(ctx context.Context) ([]*domain.Combo, error) {
	combos, err := s.comboRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return combos, nil
}
