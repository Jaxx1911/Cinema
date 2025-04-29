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

func (o OrderRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	var order domain.Order
	if err := o.db.WithContext(ctx).Preload("Tickets.Seat").First(&order, id).Error; err != nil {
		return nil, o.returnError(ctx, err)
	}
	return &order, nil
}

func (o OrderRepo) GetDetailByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	var order domain.Order
	if err := o.db.WithContext(ctx).
		Preload("Showtime").
		Preload("Showtime.Room").
		Preload("Showtime.Movie").
		Preload("Showtime.Movie.Genres").
		Preload("OrderCombos").
		Preload("OrderCombos.Combo").
		Preload("Tickets.Seat").First(&order, id).Error; err != nil {
		return nil, o.returnError(ctx, err)
	}
	return &order, nil
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

func (o OrderRepo) Update(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	if err := o.db.WithContext(ctx).Save(order).Error; err != nil {
		return nil, o.returnError(ctx, err)
	}
	return order, nil
}
