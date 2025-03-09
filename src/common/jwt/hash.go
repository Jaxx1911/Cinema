package jwt

import "golang.org/x/crypto/bcrypt"

type HashProvider interface {
	Hash(plain string) (string, error)
	ComparePassword(plain, hash string) error
}

type bcryptProvider struct {
}

func NewHashProvider() HashProvider {
	return &bcryptProvider{}
}

func (h *bcryptProvider) Hash(plain string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(hashBytes), err
}

func (h *bcryptProvider) ComparePassword(plain string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err
}
