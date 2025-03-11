package bootstrap

import (
	"TTCS/src/common/configs"
	"TTCS/src/common/log"
	"TTCS/src/core/domain"
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
		fx.Provide(repo.NewOtpRepo),
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

	err = db.AutoMigrate(
		&domain.Cinema{}, &domain.Discount{}, &domain.Movie{}, &domain.Genre{}, &domain.Room{}, &domain.Seat{},
		&domain.User{}, &domain.Order{}, &domain.Combo{}, &domain.OrderCombo{}, &domain.Ticket{},
		&domain.Showtime{}, &domain.Payment{}, &domain.Otp{},
	)
	if err != nil {
		log.Fatal("Failed to migrate Postgres")
		return nil
	}

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
