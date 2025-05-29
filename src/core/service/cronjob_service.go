package service

import (
	"TTCS/src/common/log"
	"context"
	"time"

	"github.com/robfig/cron/v3"
)

type CronjobService struct {
	movieService *MovieService
	cron         *cron.Cron
}

func NewCronjobService(movieService *MovieService) *CronjobService {
	// Create cron with timezone (UTC+7 for Vietnam)
	location, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	c := cron.New(cron.WithLocation(location))

	return &CronjobService{
		movieService: movieService,
		cron:         c,
	}
}

// Start initializes and starts all cronjobs
func (c *CronjobService) Start() error {
	ctx := context.Background()

	// Add cronjob to run at 00:00 every day to update movie status
	_, err := c.cron.AddFunc("0 0 * * *", func() {
		log.Info(ctx, "Starting daily movie status update cronjob")

		if err := c.movieService.UpdateMovieStatusOnReleaseDate(ctx); err != nil {
			log.Error(ctx, "Failed to update movie status: %v", err)
		} else {
			log.Info(ctx, "Successfully completed daily movie status update")
		}
	})

	if err != nil {
		return err
	}

	// Start the cron scheduler
	c.cron.Start()
	log.Info(ctx, "Cronjob service started successfully")

	return nil
}

// Stop gracefully stops the cronjob service
func (c *CronjobService) Stop() {
	ctx := context.Background()
	log.Info(ctx, "Stopping cronjob service...")

	stopCtx := c.cron.Stop()
	<-stopCtx.Done()

	log.Info(ctx, "Cronjob service stopped successfully")
}

// GetRunningJobs returns information about currently scheduled jobs
func (c *CronjobService) GetRunningJobs() []cron.Entry {
	return c.cron.Entries()
}
