package repo

import (
	"TTCS/src/core/domain"
	"context"
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
