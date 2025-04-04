package repo

import (
	"TTCS/src/core/domain"
	"context"
)

type TicketRepo struct {
	*BaseRepo
}

func (t TicketRepo) Create(ctx context.Context, ticket []*domain.Ticket) ([]*domain.Ticket, error) {
	if err := t.db.WithContext(ctx).Create(ticket).Error; err != nil {
		return nil, t.returnError(ctx, err)
	}
	return ticket, nil
}

func (t TicketRepo) Update(ctx context.Context, ticket *domain.Ticket) (*domain.Ticket, error) {
	if err := t.db.WithContext(ctx).Updates(ticket).Error; err != nil {
		return nil, t.returnError(ctx, err)
	}
	return ticket, nil
}

func NewTicketRepo(baseRepo *BaseRepo) domain.TicketRepo {
	return &TicketRepo{
		BaseRepo: baseRepo,
	}
}
