package infra

import (
	"context"
	"database/sql"
	"homemie/internal/repo"
	"homemie/internal/service"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func StartCronJobs(db *sql.DB, logger *zap.Logger) {
	logger.Info("Initializing cron jobs")

	// Initialize booking service for cron job
	bookingRepo := repo.NewBookingRepo(db)
	listingRepo := repo.NewListingRepo(db)
	bookingService := service.NewBookingService(bookingRepo, listingRepo)

	// Start cron job
	c := cron.New()
	c.AddFunc("@hourly", func() { bookingService.AutoCompleteBookings(context.Background()) })
	go c.Start()
}
