package bootstrap

import (
	"TTCS/src/core/service"
	"context"

	"go.uber.org/fx"
)

func BuildCronjobModule() fx.Option {
	return fx.Options(
		fx.Invoke(func(lc fx.Lifecycle, cronjobService *service.CronjobService) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return cronjobService.Start()
				},
				OnStop: func(ctx context.Context) error {
					cronjobService.Stop()
					return nil
				},
			})
		}),
	)
}
