package service

import (
	"TTCS/src/common/log"
	"TTCS/src/core/domain"
	"TTCS/src/present/httpui/request"
	"context"

	"github.com/google/uuid"
)

type OrderService struct {
	orderRepo      domain.OrderRepo
	orderComboRepo domain.OrderComboRepository
	ticketRepo     domain.TicketRepo
	discountRepo   domain.DiscountRepository
}

func NewOrderService(orderRepo domain.OrderRepo, orderComboRepo domain.OrderComboRepository, ticketRepo domain.TicketRepo, discountRepo domain.DiscountRepository) *OrderService {
	return &OrderService{
		orderRepo:      orderRepo,
		orderComboRepo: orderComboRepo,
		ticketRepo:     ticketRepo,
		discountRepo:   discountRepo,
	}
}

func (s *OrderService) Create(ctx context.Context, userId uuid.UUID, req request.CreateOrderRequest) (*domain.Order, error) {
	_ = "OrderService.Create"

	// If discount is used, decrease its usage limit
	if req.DiscountId != nil {
		discount, err := s.discountRepo.GetDiscount(ctx, *req.DiscountId)
		if err != nil {
			return nil, err
		}
		if discount.UsageLimit > 0 {
			discount.UsageLimit -= 1
			_, err = s.discountRepo.UpdateDiscount(ctx, *discount)
			if err != nil {
				return nil, err
			}
		}
	}

	order, err := s.orderRepo.Create(ctx, &domain.Order{
		ShowtimeID: req.ShowtimeID,
		UserID:     userId,
		DiscountID: req.DiscountId,
		Status:     "pending",
		TotalPrice: req.TotalPrice,
	})
	if err != nil {
		return nil, err
	}
	tickets, err := s.ticketRepo.FindByBatch(ctx, req.Tickets)
	if err != nil {
		return nil, err
	}
	for i := range tickets {
		tickets[i].OrderID = &order.ID
		tickets[i].Status = "pending"
	}
	tickets, err = s.ticketRepo.UpdateBatch(ctx, tickets)
	if err != nil {
		return nil, err
	}
	var orderCombos []domain.OrderCombo
	for _, combo := range req.Combos {
		orderCombo, err := s.orderComboRepo.Create(ctx, &domain.OrderCombo{
			OrderID:    order.ID,
			ComboID:    combo.Id,
			Quantity:   combo.Quantity,
			TotalPrice: combo.TotalPrice,
		})
		if err != nil {
			return nil, err
		}
		orderCombos = append(orderCombos, *orderCombo)
	}
	order.OrderCombos = orderCombos
	order.Tickets = tickets
	return order, nil
}

func (s *OrderService) GetById(ctx context.Context, orderId string) (*domain.Order, error) {
	caller := "OrderService.GetById"

	uid, err := uuid.Parse(orderId)
	if err != nil {
		log.Error(ctx, "[%v] failed to get order detail %+v", caller, err)
		return nil, err
	}

	order, err := s.orderRepo.GetByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) GetDetailById(ctx context.Context, orderId string) (*domain.Order, error) {
	caller := "OrderService.GetDetailById"

	uid, err := uuid.Parse(orderId)
	if err != nil {
		log.Error(ctx, "[%v] failed to get order detail %+v", caller, err)
		return nil, err
	}

	order, err := s.orderRepo.GetDetailByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) Delete(ctx context.Context, orderId string) error {
	err := s.orderRepo.Delete(ctx, uuid.MustParse(orderId))
	if err != nil {
		return err
	}
	return nil
}
