package response

import (
	"TTCS/src/core/domain"
	"time"

	"github.com/google/uuid"
)

type Otp struct {
	Otp string `json:"otp"`
}

type LoginResp struct {
	Token *Token `json:"token"`
}

type Token struct {
	AccessToken      string `json:"access_token"`
	AccessExpiredAt  int64  `json:"access_expired_at"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiredAt int64  `json:"refresh_expired_at"`
}

type User struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	Role      string     `json:"role"`
	AvatarURL string     `json:"avatar_url"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func UserFromDomain(d *domain.User) *User {
	var deletedAt *time.Time
	if d.DeletedAt.Valid {
		deletedAt = &d.DeletedAt.Time
	}
	return &User{
		ID:        d.ID,
		Name:      d.Name,
		Email:     d.Email,
		Phone:     d.Phone,
		Role:      d.Role,
		AvatarURL: d.AvatarUrl,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func UsersFromDomain(d []*domain.User) []*User {
	res := make([]*User, 0)
	for _, u := range d {
		res = append(res, UserFromDomain(u))
	}
	return res
}
