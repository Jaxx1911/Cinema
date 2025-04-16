package bootstrap

import (
	"TTCS/src/common/configs"
	"TTCS/src/common/log"
	"TTCS/src/core/domain"
	"TTCS/src/infra/cache"
	"TTCS/src/infra/repo"
	"TTCS/src/infra/upload"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func BuildDatabasesModule() fx.Option {
	return fx.Options(
		fx.Provide(NewPostgresDB),
		fx.Provide(NewMinioClient),
		fx.Provide(upload.NewUploadService),
		fx.Provide(cache.NewRedisClient),
		fx.Provide(repo.NewBaseRepo),
		fx.Provide(repo.NewOtpRepo),
		fx.Provide(repo.NewUserRepo),
		fx.Provide(repo.NewGenreRepo),
		fx.Provide(repo.NewMovieRepo),
		fx.Provide(repo.NewShowtimeRepo),
		fx.Provide(repo.NewRoomRepo),
		fx.Provide(repo.NewTicketRepo),
		fx.Provide(repo.NewCinemaRepo),
		fx.Provide(repo.NewSeatRepo),
		fx.Provide(repo.NewComboRepo),
		fx.Provide(repo.NewDiscountRepo),
		fx.Provide(repo.NewOrderComboRepo),
		fx.Provide(repo.NewOrderRepo),
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

func NewMinioClient() *minio.Client {
	ctx := context.Background()

	minioClient, err := minio.New(configs.GetConfig().Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(configs.GetConfig().Minio.AccessKey, configs.GetConfig().Minio.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal("Failed to create minio client")
		return nil
	}

	err = minioClient.MakeBucket(ctx, configs.GetConfig().Minio.Bucket, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, configs.GetConfig().Minio.Bucket)
		if errBucketExists == nil && exists {
			log.Info(ctx, "Bucket already exists")
		} else {
			log.Fatal("Failed to create minio bucket")
		}
	} else {
		log.Info(ctx, fmt.Sprintf("Successfully created minio bucket %s", configs.GetConfig().Minio.Bucket))
	}
	return minioClient
}
