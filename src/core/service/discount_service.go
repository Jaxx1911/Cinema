package service

import (
	"TTCS/src/common/fault"
	"TTCS/src/core/domain"
	"context"
	"github.com/google/uuid"
)

type DiscountService struct {
	discountRepo domain.DiscountRepository
}

func NewDiscountService(discountRepo domain.DiscountRepository) *DiscountService {
	return &DiscountService{discountRepo: discountRepo}
}

func (d DiscountService) GetDiscount(ctx context.Context, id string) (*domain.Discount, error) {
	caller := "DiscountService.GetDiscount"
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] invalid id", caller)
	}
	return d.discountRepo.GetDiscount(ctx, uid)
}

func (d DiscountService) GetListDiscount(ctx context.Context) ([]domain.Discount, error) {
	return d.discountRepo.GetListDiscount(ctx)
}

func (d DiscountService) GetDiscountByCode(ctx context.Context, code string) (domain.Discount, error) {
	return d.GetDiscountByCode(ctx, code)
}
