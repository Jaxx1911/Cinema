package bootstrap

import (
	"TTCS/src/present/httpui/controller"
	"go.uber.org/fx"
)

func BuildControllers() fx.Option {
	return fx.Options(
		fx.Provide(controller.NewBaseController),
		fx.Provide(controller.NewAuthController),
	)
}
