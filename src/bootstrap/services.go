package bootstrap

import (
	"TTCS/src/core/service"
	"go.uber.org/fx"
)

func BuildServices() fx.Option {
	return fx.Options(
		fx.Provide(service.NewAuthService),
		fx.Provide(service.NewUserService),
		fx.Provide(service.NewMovieService),
		fx.Provide(service.NewShowtimeService),
	)
}
