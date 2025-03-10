package bootstrap

import (
	"TTCS/src/common/crypto"
	"go.uber.org/fx"
)

func BuildCrypto() fx.Option {
	return fx.Options(
		fx.Provide(crypto.NewHashProvider),
	)
}
