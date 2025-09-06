package infra

import (
	"homemie/internal/repository"
	"homemie/internal/service"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func StartCronJobs(db *gorm.DB, logger *zap.Logger) {
	logger.Info("Initializing cron jobs")

	// Initialize booking service for cron job
	bookingRepo := repository.NewBookingRepo(db, logger.Named("booking_repo"))
	listingRepo := repository.NewListingRepo(db, logger.Named("listing_repo"))
	bookingService := service.NewBookingService(bookingRepo, listingRepo, logger.Named("booking_service"))

	// Start cron job
	c := cron.New()
	c.AddFunc("@hourly", bookingService.AutoCompleteBookings)
	go c.Start()
}
