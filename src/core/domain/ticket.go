package domain

type Ticket struct {
	ID         uint   `gorm:"primaryKey"`
	OrderID    uint   `gorm:"not null"`
	ShowtimeID uint   `gorm:"not null"`
	SeatID     uint   `gorm:"not null"`
	Status     string `gorm:"type:varchar(20);not null"` // booked, canceled, paid
}

func (*Ticket) TableName() string {
	return "ticket"
}
