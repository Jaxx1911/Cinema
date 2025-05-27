package bootstrap

import (
	"TTCS/src/common/log"
	"TTCS/src/common/ws"
	"TTCS/src/present/httpui/router"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
	"os"
)

func BuildHTTPServerModule() fx.Option {
	return fx.Options(
		fx.Provide(gin.New),
		fx.Provide(ws.NewHub),
		fx.Invoke(router.RegisterHandler),
		fx.Invoke(router.RegisterGinRouters),
		fx.Invoke(newWsHub),
		fx.Invoke(newHttpServer),
	)
}

func newWsHub(lc fx.Lifecycle, hub *ws.Hub) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go hub.Run()
			return nil
		},
	})
}

func newHttpServer(lc fx.Lifecycle, engine *gin.Engine) {
	server := &http.Server{
		Addr:    os.Getenv("SERVER_ADDRESS"),
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
