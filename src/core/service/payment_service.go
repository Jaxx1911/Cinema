package service

import (
	"TTCS/src/common/log"
	"TTCS/src/core/domain"
	"TTCS/src/present/httpui/request"
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

type PaymentService struct {
	paymentRepo domain.PaymentRepo
	orderRepo   domain.OrderRepo
}

func NewPaymentService(paymentRepo domain.PaymentRepo, orderRepo domain.OrderRepo) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
	}
}

func (p *PaymentService) HandleCallback(ctx context.Context, callback request.PaymentCallback) (*domain.Payment, error) {
	//log.Info(ctx, "receive callback test %v", callback)
	//return &domain.Payment{}, nil
	caller := "PaymentService.HandleCallback"
	oid, err := uuid.Parse(callback.Payment.Content)
	if err != nil {
		log.Error(ctx, "[%v] invalid content %+v", caller, err)
		return nil, err
	}
	order, err := p.orderRepo.GetByID(ctx, oid)
	if err != nil {
		return nil, err
	}
	if int64(order.TotalPrice) != callback.Payment.Amount {
		err := errors.New("invalid amount")
		log.Error(ctx, "[%v] invalid amount %+v", caller, err)
		return nil, err
	}

	payment, err := p.paymentRepo.Create(ctx, &domain.Payment{
		UserID:        order.UserID,
		OrderID:       order.ID,
		TransactionID: callback.Payment.TransactionId,
		Status:        "success",
		Amount:        callback.Payment.Amount,
		PaymentTime:   time.Now(),
	})
	order.Status = "success"

	order, err = p.orderRepo.Update(ctx, order)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (p *PaymentService) GetPaymentsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Payment, error) {
	payments, err := p.paymentRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return payments, nil
}
