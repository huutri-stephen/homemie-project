package infra

import (
    "log"
    "gorm.io/gorm"

	"homemie/internal/repository"
	"homemie/internal/service"

	"github.com/robfig/cron/v3"
)

func SeedData(db *gorm.DB) {
    log.Println("Seeding mock data...")
    log.Println("Seeded successfully.")
}

func CronJobs(db *gorm.DB) {
    // Initialize booking service for cron job
    bookingRepo := repository.NewBookingRepo(db)
    listingRepo := repository.NewListingRepo(db)
    bookingService := service.NewBookingService(bookingRepo, listingRepo)

    // Start cron job
    c := cron.New()
    c.AddFunc("0 0 0 * * *", bookingService.AutoCompleteBookings) // Runs daily at 00:00:00
    go c.Start()
}