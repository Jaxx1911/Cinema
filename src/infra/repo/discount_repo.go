package repo

import (
	"TTCS/src/core/domain"
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiscountRepo struct {
	*BaseRepo
}

func (d DiscountRepo) GetDiscount(ctx context.Context, id uuid.UUID) (*domain.Discount, error) {
	var discount domain.Discount
	if err := d.db.WithContext(ctx).First(&discount, id).Error; err != nil {
		return nil, d.returnError(ctx, err)
	}
	return &discount, nil
}

func (d DiscountRepo) GetListDiscount(ctx context.Context) ([]domain.Discount, error) {
	var discount []domain.Discount
	if err := d.db.WithContext(ctx).Find(&discount).Error; err != nil {
		return nil, d.returnError(ctx, err)
	}
	return discount, nil
}

func (d DiscountRepo) GetDiscountByCode(ctx context.Context, code string) (*domain.Discount, error) {
	var discount domain.Discount
	if err := d.db.WithContext(ctx).First(&discount, "code = ?", code).Error; err != nil {
		return nil, d.returnError(ctx, err)
	}
	return &discount, nil
}

func (d DiscountRepo) CreateDiscount(ctx context.Context, discount domain.Discount) (*domain.Discount, error) {
	dc, err := d.GetDiscountByCode(ctx, discount.Code)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, d.returnError(ctx, err)
	}
	if dc != nil {
		return nil, d.returnError(ctx, errors.New("discount already exists"))
	}
	if err := d.db.WithContext(ctx).Create(&discount).Error; err != nil {
		return nil, d.returnError(ctx, err)
	}
	return &discount, nil
}

func (d DiscountRepo) UpdateDiscount(ctx context.Context, discount domain.Discount) (*domain.Discount, error) {
	if err := d.db.WithContext(ctx).Save(&discount).Error; err != nil {
		return nil, d.returnError(ctx, err)
	}
	return &discount, nil
}

func NewDiscountRepo(basRepo *BaseRepo) domain.DiscountRepository {
	return DiscountRepo{
		BaseRepo: basRepo,
	}
}
