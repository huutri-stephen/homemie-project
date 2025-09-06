package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"homemie/internal/handler"
	"homemie/internal/repository"
	"homemie/internal/service"
)

// InitBookingRoutes đăng ký các route liên quan đến booking
func InitBookingRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *zap.Logger) {
	bookingRepo := repository.NewBookingRepo(db, logger.Named("booking_repo"))
	listingRepo := repository.NewListingRepo(db, logger.Named("listing_repo"))
	svc := service.NewBookingService(bookingRepo, listingRepo, logger.Named("booking_service"))
	h := handler.NewBookingHandler(svc, logger.Named("booking_handler"))

	bookings := rg.Group("/bookings")
	{
		bookings.POST("", h.CreateBooking)
		bookings.GET("/my-bookings", h.GetMyBookings)
		bookings.GET("/owner-bookings", h.GetOwnerBookings)

		bookings.PUT("/:id/respond", h.RespondToBooking)
		bookings.PUT("/:id/cancel", h.CancelBooking)
	}
}
