package domain

import "github.com/google/uuid"

type Ticket struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID    uuid.UUID `gorm:"type:uuid;not null"`
	ShowtimeID uuid.UUID `gorm:"type:uuid;not null"`
	SeatID     uuid.UUID `gorm:"type:uuid;not null"`
	Status     string    `gorm:"type:varchar(20);not null"`

	Order    Order    `gorm:"foreignKey:OrderID"`
	Showtime Showtime `gorm:"foreignKey:ShowtimeID"`
	Seat     Seat     `gorm:"foreignKey:SeatID"`
}

func (*Ticket) TableName() string {
	return "ticket"
}
