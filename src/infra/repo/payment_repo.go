package repo

import (
	"TTCS/src/core/domain"
	"context"
	"time"

	"github.com/google/uuid"
)

type PaymentRepo struct {
	*BaseRepo
}

func NewPaymentRepo(baseRepo *BaseRepo) domain.PaymentRepo {
	return PaymentRepo{
		BaseRepo: baseRepo,
	}
}

func (p PaymentRepo) Create(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	if err := p.db.WithContext(ctx).Create(payment).Error; err != nil {
		return nil, p.returnError(ctx, err)
	}
	return payment, nil
}

func (p PaymentRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Payment, error) {
	var payments []domain.Payment
	if err := p.db.WithContext(ctx).Where("user_id = ?", userID).Find(&payments).Error; err != nil {
		return nil, p.returnError(ctx, err)
	}
	return payments, nil
}

func (p PaymentRepo) GetByCinemaID(ctx context.Context, cinemaID uuid.UUID) ([]domain.Payment, error) {
	var payments []domain.Payment
	if err := p.db.WithContext(ctx).
		Model(&domain.Payment{}).
		Joins("JOIN orders ON payment.order_id = orders.id").
		Joins("JOIN showtime ON orders.showtime_id = showtime.id").
		Joins("JOIN room ON showtime.room_id = room.id").
		Where("room.cinema_id = ?", cinemaID).
		Find(&payments).Error; err != nil {
		return nil, p.returnError(ctx, err)
	}
	return payments, nil
}

func (p PaymentRepo) GetByCinemaIDAndDateRange(ctx context.Context, cinemaID uuid.UUID, startDate, endDate time.Time) ([]domain.Payment, error) {
	var payments []domain.Payment
	if err := p.db.WithContext(ctx).
		Model(&domain.Payment{}).
		Joins("JOIN orders ON payment.order_id = orders.id").
		Joins("JOIN showtime ON orders.showtime_id = showtime.id").
		Joins("JOIN room ON showtime.room_id = room.id").
		Where("room.cinema_id = ? AND payment.payment_time >= ? AND payment.payment_time <= ?", cinemaID, startDate, endDate).
		Find(&payments).Error; err != nil {
		return nil, p.returnError(ctx, err)
	}
	return payments, nil
}
