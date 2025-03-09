package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Token struct {
	Token  string `json:"token"`
	Expire int64  `json:"expire"`
}

type Payload struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type myClaims struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Generate(secret string, payload Payload, expire int64) (*Token, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &myClaims{
		payload.Id,
		payload.Username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(expire))),
		},
	})

	myToken, err := t.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	return &Token{myToken, expire}, nil
}

var ErrTokenExpired = errors.New("Token is expired")

func Verify(secret string, token string) (*Payload, error) {
	res, err := jwt.ParseWithClaims(token, &myClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if errors.Is(err, jwt.ErrTokenExpired) {
		return nil, ErrTokenExpired
	}
	if err != nil {
		return nil, err
	}
	if res.Valid {
		return nil, fmt.Errorf("invalid token %s", token)
	}
	claims, ok := res.Claims.(*myClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims with token %s", token)
	}
	return &Payload{claims.Id, claims.Username}, nil
}
