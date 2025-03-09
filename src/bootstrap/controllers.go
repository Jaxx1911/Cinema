package bootstrap

import (
	"TTCS/src/present/httpui/controller"
	"TTCS/src/present/httpui/validator"
	"go.uber.org/fx"
)

func BuildControllers() fx.Option {
	return fx.Options(
		fx.Provide(controller.NewBaseController),
		fx.Provide(controller.NewAuthController),
	)
}

func BuildValidators() fx.Option {
	return fx.Provide(validator.NewValidator)
}
