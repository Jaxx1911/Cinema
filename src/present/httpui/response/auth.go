package response

import (
	"TTCS/src/core/domain"
	"github.com/google/uuid"
	"time"
)

type Otp struct {
	Otp string `json:"otp"`
}

type LoginResp struct {
	Token *Token `json:"token"`
	User  *User  `json:"user"`
}

type Token struct {
	AccessToken      string `json:"access_token"`
	AccessExpiredAt  int64  `json:"access_expired_at"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiredAt int64  `json:"refresh_expired_at"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func UserFromDomain(d *domain.User) *User {
	return &User{
		ID:        d.ID,
		Name:      d.Name,
		Email:     d.Email,
		Phone:     d.Phone,
		Role:      d.Role,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
