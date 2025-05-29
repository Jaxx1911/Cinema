package main

import (
	"TTCS/src/bootstrap"
	"TTCS/src/common/genqr"
	"TTCS/src/common/log"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	"go.uber.org/fx"
)

const (
	defaultGracefulTimeout = 15 * time.Second
)

func init() {
	//var pathConfig string
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load the env vars: %v", err)
	}

	//flag.StringVar(&pathConfig, "config", ".env", "path to config file")
	//flag.Parse()

	//err := configs.InitConfig(pathConfig)
	//if err != nil {
	//	panic(err)
	//}
	log.NewLogger()
	genqr.InitQrService()
	//Json float to float (default float to string)
	decimal.MarshalJSONWithoutQuotes = true
}

func main() {
	log.Debug(context.Background(), "App %s is running", os.Getenv("SERVER_NAME"))
	app := fx.New(
		bootstrap.BuildCrypto(),
		bootstrap.BuildMailService(),
		bootstrap.BuildDatabasesModule(),
		bootstrap.BuildHTTPServerModule(),
		bootstrap.BuildServices(),
		bootstrap.BuildValidators(),
		bootstrap.BuildMiddlewares(),
		bootstrap.BuildControllers(),
		bootstrap.BuildCronjobModule(),
	)
	startContext, cancel := context.WithTimeout(context.Background(), defaultGracefulTimeout)
	defer cancel()
	if err := app.Start(startContext); err != nil {
		log.Fatal(err.Error())
	}
	interruptHandle(app)
}

func interruptHandle(app *fx.App) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	ctx := context.Background()
	log.Debug(ctx, "Listening Signal...")
	s := <-c
	log.Info(ctx, "Received signal: %s. Shutting down Server ...", s)

	stopCtx, cancel := context.WithTimeout(context.Background(), defaultGracefulTimeout)
	defer cancel()

	if err := app.Stop(stopCtx); err != nil {
		log.Fatal(err.Error())
	}
}
