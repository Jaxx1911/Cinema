package bootstrap

import (
	"TTCS/src/common/jwt"
	"go.uber.org/fx"
)

func BuildCrypto() fx.Option {
	return fx.Options(
		fx.Provide(jwt.NewHashProvider),
	)
}
