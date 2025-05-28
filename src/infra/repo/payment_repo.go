package repo

import (
	"TTCS/src/core/domain"
	"TTCS/src/present/httpui/request"
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
		Preload("User").
		Preload("Order").
		Preload("Order.Discount").
		Joins("JOIN orders ON payment.order_id = orders.id").
		Joins("JOIN showtime ON orders.showtime_id = showtime.id").
		Joins("JOIN room ON showtime.room_id = room.id").
		Where("room.cinema_id = ? AND payment.payment_time >= ? AND payment.payment_time <= ?", cinemaID, startDate, endDate).
		Find(&payments).Error; err != nil {
		return nil, p.returnError(ctx, err)
	}
	return payments, nil
}

func (p PaymentRepo) GetList(ctx context.Context, req request.GetListPaymentRequest) ([]domain.Payment, int64, error) {
	var payments []domain.Payment
	var total int64

	query := p.db.WithContext(ctx).Model(&domain.Payment{}).
		Preload("User").
		Preload("Order").
		Preload("Order.Discount")

	// Apply status filter
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, p.returnError(ctx, err)
	}

	// Apply sorting
	sortBy := req.SortBy
	if sortBy == "" {
		sortBy = "date_desc" // Default sort by date descending (newest first)
	}

	switch sortBy {
	case "date_desc":
		query = query.Order("payment_time DESC")
	case "date_asc":
		query = query.Order("payment_time ASC")
	case "amount_desc":
		query = query.Order("amount DESC")
	case "amount_asc":
		query = query.Order("amount ASC")
	default:
		query = query.Order("payment_time DESC")
	}

	// Apply pagination
	limit, offset := p.toLimitOffset(ctx, req.Page)
	if err := query.Offset(offset).Limit(limit).Find(&payments).Error; err != nil {
		return nil, 0, p.returnError(ctx, err)
	}

	return payments, total, nil
}
