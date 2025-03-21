package domain

import (
	"context"
	"github.com/google/uuid"
)

type Ticket struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID    *uuid.UUID `gorm:"type:uuid"`
	ShowtimeID uuid.UUID  `gorm:"type:uuid;not null"`
	SeatID     uuid.UUID  `gorm:"type:uuid;not null"`
	Status     string     `gorm:"type:varchar(20);not null"`

	Order    Order    `gorm:"foreignKey:OrderID"`
	Showtime Showtime `gorm:"foreignKey:ShowtimeID"`
	Seat     Seat     `gorm:"foreignKey:SeatID"`
}

type TicketRepo interface {
	Create(ctx context.Context, ticket []*Ticket) ([]*Ticket, error)
	Update(ctx context.Context, ticket *Ticket) (*Ticket, error)
}

func (*Ticket) TableName() string {
	return "ticket"
}
