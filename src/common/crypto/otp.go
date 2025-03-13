package crypto

import (
	"fmt"
	"math/rand"
)

type OTPProvider interface {
	GenerateOTP() string
}

type totpProvider struct{}

func NewOTPProvider() OTPProvider {
	return &totpProvider{}
}

func (provider *totpProvider) GenerateOTP() string {
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	return otp
}
