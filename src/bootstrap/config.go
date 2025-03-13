package bootstrap

import (
	"TTCS/src/common/crypto"
	"TTCS/src/common/mail"
	"go.uber.org/fx"
)

func BuildCrypto() fx.Option {
	return fx.Options(
		fx.Provide(crypto.NewHashProvider),
		fx.Provide(crypto.NewJwtProvider),
		fx.Provide(crypto.NewOTPProvider),
	)
}

func BuildMailService() fx.Option {
	return fx.Options(
		fx.Provide(mail.NewGmailService))
}
