package bootstrap

import (
	"TTCS/src/present/httpui/controller"
	"TTCS/src/present/httpui/middleware"
	"TTCS/src/present/httpui/validator"
	"go.uber.org/fx"
)

func BuildControllers() fx.Option {
	return fx.Options(
		fx.Provide(controller.NewBaseController),
		fx.Provide(controller.NewAuthController),
		fx.Provide(controller.NewUserController),
		fx.Provide(controller.NewMovieController),
		fx.Provide(controller.NewShowtimeController),
	)
}

func BuildMiddlewares() fx.Option {
	return fx.Options(
		fx.Provide(middleware.NewAuthMiddleware),
	)
}

func BuildValidators() fx.Option {
	return fx.Provide(validator.NewValidator)
}
