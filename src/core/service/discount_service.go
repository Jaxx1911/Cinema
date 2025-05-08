package service

import (
	"TTCS/src/common/log"
	"TTCS/src/core/domain"
	"TTCS/src/present/httpui/request"
	"context"
	"github.com/google/uuid"
	"time"
)

type DiscountService struct {
	discountRepo domain.DiscountRepository
}

func NewDiscountService(discountRepo domain.DiscountRepository) *DiscountService {
	return &DiscountService{discountRepo: discountRepo}
}

func (d DiscountService) GetDiscount(ctx context.Context, id uuid.UUID) (*domain.Discount, error) {
	return d.discountRepo.GetDiscount(ctx, id)
}

func (d DiscountService) GetListDiscount(ctx context.Context) ([]domain.Discount, error) {
	return d.discountRepo.GetListDiscount(ctx)
}

func (d DiscountService) GetDiscountByCode(ctx context.Context, code string) (domain.Discount, error) {
	return d.GetDiscountByCode(ctx, code)
}

func (d DiscountService) Create(ctx context.Context, req request.Discount) (*domain.Discount, error) {
	caller := "DiscountService.Create"
	location := time.Now().Location()
	start, err := time.ParseInLocation("2006-01-02", req.StartDate, location)
	if err != nil {
		log.Error(ctx, "[%v] invalid start date format +%v", caller, err)
		return nil, err
	}
	end, err := time.ParseInLocation("2006-01-02", req.EndDate, location)
	if err != nil {
		log.Error(ctx, "[%v] invalid end date format +%v", caller, err)
		return nil, err
	}

	return d.discountRepo.CreateDiscount(ctx, domain.Discount{
		Code:       req.Code,
		Percentage: req.Percentage,
		StartDate:  start,
		EndDate:    end,
	})
}

func (d DiscountService) Update(ctx context.Context, id uuid.UUID, req request.Discount) (*domain.Discount, error) {
	caller := "DiscountService.Update"
	location := time.Now().Location()
	start, err := time.ParseInLocation("2006-01-02", req.StartDate, location)
	if err != nil {
		log.Error(ctx, "[%v] invalid start date format +%v", caller, err)
		return nil, err
	}
	end, err := time.ParseInLocation("2006-01-02", req.EndDate, location)
	if err != nil {
		log.Error(ctx, "[%v] invalid end date format +%v", caller, err)
		return nil, err
	}

	discount, err := d.discountRepo.GetDiscount(ctx, id)
	if err != nil {
		return nil, err
	}
	discount.Code = req.Code
	discount.Percentage = req.Percentage
	discount.StartDate = start
	discount.EndDate = end
	return d.discountRepo.UpdateDiscount(ctx, *discount)
}
