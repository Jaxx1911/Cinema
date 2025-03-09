package domain

type Auth struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	UserID       uint   `gorm:"uniqueIndex;not null" json:"user_id"`
	PasswordHash string `gorm:"type:text;not null" json:"password_hash"`
}

type AuthRepository interface{}

func (*Auth) TableName() string {
	return "auth"
}
