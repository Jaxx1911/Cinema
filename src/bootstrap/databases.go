package bootstrap

import (
	"TTCS/src/common/configs"
	"TTCS/src/common/log"
	"TTCS/src/infra/cache"
	"TTCS/src/infra/repo"
	"context"
	"fmt"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func BuildDatabasesModule() fx.Option {
	return fx.Options(
		fx.Provide(NewPostgresDB),
		fx.Provide(cache.NewRedisClient),
		fx.Provide(repo.NewBaseRepo),
		fx.Provide(repo.NewAuthRepo),
		fx.Provide(repo.NewUserRepo),
	)
}

func NewPostgresDB(lc fx.Lifecycle) *gorm.DB {
	cf := configs.GetConfig().Postgres
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		cf.Host, cf.Port, cf.User, cf.DB, cf.SslMode, cf.Password)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database")
		return nil
	}
	log.Info(context.Background(), "Successfully connected to Postgres")

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Info(ctx, "Closing Postgres connection")
			sqlDB, err := db.DB()
			if err != nil {
				return err
			}
			return sqlDB.Close()
		},
	})
	return db
}
