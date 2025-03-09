package bootstrap

import (
	"TTCS/src/common/configs"
	"TTCS/src/common/log"
	"TTCS/src/present/httpui/router"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
)

func BuildHTTPServerModule() fx.Option {
	return fx.Options(
		fx.Provide(gin.New),
		fx.Invoke(router.RegisterHandler),
		fx.Invoke(router.RegisterGinRouters),
		fx.Invoke(NewHttpServer),
	)
}

func NewHttpServer(lc fx.Lifecycle, engine *gin.Engine) {
	server := &http.Server{
		Addr:    configs.GetConfig().Server.Address,
		Handler: engine,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatal("Could not listen and serve: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info(ctx, "Shutting down server...")
			return server.Shutdown(ctx)
		},
	})
}
