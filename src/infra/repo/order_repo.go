package repo

import (
	"TTCS/src/core/domain"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (o OrderComboRepo) GetAll(ctx context.Context) ([]domain.OrderCombo, error) {
	var orderCombos []domain.OrderCombo
	if err := o.db.WithContext(ctx).Find(&orderCombos).Error; err != nil {
		return nil, o.returnError(ctx, err)
	}
	return orderCombos, nil
}

func (o OrderRepo) Update(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	if err := o.db.WithContext(ctx).Save(order).Error; err != nil {
		return nil, o.returnError(ctx, err)
	}
	return order, nil
}

func (o OrderRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get order first to check for discount
		var order domain.Order
		if err := tx.First(&order, id).Error; err != nil {
			return o.returnError(ctx, err)
		}

		// If order has discount, restore its usage limit
		if order.DiscountID != nil {
			var discount domain.Discount
			if err := tx.First(&discount, order.DiscountID).Error; err != nil {
				return o.returnError(ctx, err)
			}
			discount.UsageLimit += 1
			if err := tx.Save(&discount).Error; err != nil {
				return o.returnError(ctx, err)
			}
		}

		// Delete related order_combos
		if err := tx.Where("order_id = ?", id).Delete(&domain.OrderCombo{}).Error; err != nil {
			return o.returnError(ctx, err)
		}

		// Set order_id to null for related tickets
		if err := tx.Model(&domain.Ticket{}).
			Where("order_id = ?", id).
			Updates(map[string]interface{}{
				"order_id": nil,
				"status":   "available",
			}).Error; err != nil {
			return o.returnError(ctx, err)
		}

		// Finally delete the order
		if err := tx.Delete(&domain.Order{}, id).Error; err != nil {
			return o.returnError(ctx, err)
		}

		return nil
	})
}

func (o OrderRepo) GetOrdersByDateRange(ctx context.Context, start time.Time, end time.Time) ([]domain.Order, error) {
	var orders []domain.Order
	if err := o.db.WithContext(ctx).Where("created_at BETWEEN ? AND ?", start, end).Find(&orders).Error; err != nil {
		return nil, o.returnError(ctx, err)
	}
	return orders, nil
}

func (o OrderRepo) GetAllPendingOrders(ctx context.Context) ([]domain.Order, error) {
	var orders []domain.Order
	if err := o.db.WithContext(ctx).Where("status = ?", "pending").Find(&orders).Error; err != nil {
		return nil, o.returnError(ctx, err)
	}
	return orders, nil
}
