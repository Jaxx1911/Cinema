package repo

import (
	"TTCS/src/core/domain"
	"context"
	"github.com/google/uuid"
)

type OrderRepo struct {
	*BaseRepo
}

func NewOrderRepo(baseRepo *BaseRepo) domain.OrderRepo {
	return &OrderRepo{BaseRepo: baseRepo}
}

func (o OrderRepo) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	if err := o.db.WithContext(ctx).Create(order).Error; err != nil {
		return nil, o.returnError(ctx, err)
	}
	return order, nil
}

type OrderComboRepo struct {
	*BaseRepo
}

func NewOrderComboRepo(baseRepo *BaseRepo) domain.OrderComboRepository {
	return &OrderComboRepo{baseRepo}
}

func (o OrderComboRepo) Create(ctx context.Context, orderCombo *domain.OrderCombo) (*domain.OrderCombo, error) {
	if err := o.db.WithContext(ctx).Create(orderCombo).Error; err != nil {
		return nil, o.returnError(ctx, err)
	}
	return orderCombo, nil
}

func (o OrderComboRepo) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]domain.OrderCombo, error) {
	var order []domain.OrderCombo
	if err := o.db.WithContext(ctx).Find(&order, "order_id = ?", orderID).Error; err != nil {
		return nil, o.returnError(ctx, err)
	}
	return order, nil
}
