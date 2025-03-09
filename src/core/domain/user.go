package domain

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Phone     string    `gorm:"type:varchar(20);uniqueIndex;not null"`
	Role      string    `gorm:"type:varchar(50);default:customer;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Auth   Auth    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Orders []Order `gorm:"foreignKey:UserID"`
}

type UserRepository interface {
}

func (*User) TableName() string {
	return "users"
}
