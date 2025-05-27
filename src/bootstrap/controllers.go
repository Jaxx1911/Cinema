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
		fx.Provide(controller.NewCinemaController),
		fx.Provide(controller.NewSeatController),
		fx.Provide(controller.NewRoomController),
		fx.Provide(controller.NewComboController),
		fx.Provide(controller.NewDiscountController),
		fx.Provide(controller.NewOrderController),
		fx.Provide(controller.NewPaymentController),
		fx.Provide(controller.NewGenreController),
		fx.Provide(controller.NewWebSocketController),
		fx.Provide(controller.NewStatisticController),
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
