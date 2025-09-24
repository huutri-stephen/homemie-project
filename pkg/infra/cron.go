package infra

import (
	"homemie/internal/repo"
	"homemie/internal/service"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func StartCronJobs(db *gorm.DB, logger *zap.Logger) {
	logger.Info("Initializing cron jobs")

	// Initialize booking service for cron job
	bookingRepo := repo.NewBookingRepo(db, logger.Named("booking_repo"))
	listingRepo := repo.NewListingRepo(db, logger.Named("listing_repo"))
	bookingService := service.NewBookingService(bookingRepo, listingRepo, logger.Named("booking_service"))

	// Start cron job
	c := cron.New()
	c.AddFunc("@hourly", bookingService.AutoCompleteBookings)
	go c.Start()
}
