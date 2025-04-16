package service

import (
	"TTCS/src/core/domain"
	"TTCS/src/present/httpui/request"
	"context"
	"github.com/google/uuid"
)

type OrderService struct {
	orderRepo      domain.OrderRepo
	orderComboRepo domain.OrderComboRepository
	ticketRepo     domain.TicketRepo
}

func NewOrderService(orderRepo domain.OrderRepo, orderComboRepo domain.OrderComboRepository, ticketRepo domain.TicketRepo) *OrderService {
	return &OrderService{
		orderRepo:      orderRepo,
		orderComboRepo: orderComboRepo,
		ticketRepo:     ticketRepo,
	}
}

func (s *OrderService) Create(ctx context.Context, userId uuid.UUID, req request.CreateOrderRequest) (*domain.Order, error) {
	_ = "OrderService.Create"
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
