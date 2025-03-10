package crypto

import (
	"github.com/pquerna/otp/totp"
	"time"
)

type OTPProvider interface {
	GenerateCode(secret string) (string, error)
	Validate(passcode, secret string) bool
}

type totpProvider struct{}

func NewOTPProvider() OTPProvider {
	return &totpProvider{}
}

func (provider *totpProvider) GenerateCode(secret string) (string, error) {
	passcode, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		return "", err
	}
	return passcode, nil
}

func (provider *totpProvider) Validate(passcode, secret string) bool {
	valid := totp.Validate(passcode, secret)
	return valid
}
