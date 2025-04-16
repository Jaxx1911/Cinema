package repo

import (
	"TTCS/src/core/domain"
	"context"
	"github.com/google/uuid"
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

func (t TicketRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Ticket, error) {
	var ticket domain.Ticket
	if err := t.db.WithContext(ctx).First(&ticket, id).Error; err != nil {
		return nil, t.returnError(ctx, err)
	}
	return &ticket, nil
}

func (t TicketRepo) FindByBatch(ctx context.Context, ids []uuid.UUID) ([]domain.Ticket, error) {
	var tickets []domain.Ticket
	if err := t.db.WithContext(ctx).Where("id IN ?", ids).Find(&tickets).Error; err != nil {
		return nil, t.returnError(ctx, err)
	}
	return tickets, nil
}

func (t TicketRepo) UpdateBatch(ctx context.Context, tickets []domain.Ticket) ([]domain.Ticket, error) {
	tx := t.db.WithContext(ctx).Begin()

	for _, ticket := range tickets {
		if err := tx.Updates(ticket).Error; err != nil {
			tx.Rollback()
			return nil, t.returnError(ctx, err)
		}
	}
	tx.Commit()
	return tickets, nil
}

func NewTicketRepo(baseRepo *BaseRepo) domain.TicketRepo {
	return &TicketRepo{
		BaseRepo: baseRepo,
	}
}
