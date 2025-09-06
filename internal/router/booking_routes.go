package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"homemie/internal/handler"
	"homemie/internal/repository"
	"homemie/internal/service"
)

// InitBookingRoutes đăng ký các route liên quan đến booking
func InitBookingRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	bookingRepo := repository.NewBookingRepo(db)
	listingRepo := repository.NewListingRepo(db)
	svc := service.NewBookingService(bookingRepo, listingRepo)
	h := handler.NewBookingHandler(svc)

	bookings := rg.Group("/bookings")
	{
		bookings.POST("", h.CreateBooking)
		bookings.GET("/my-bookings", h.GetMyBookings)
		bookings.GET("/owner-bookings", h.GetOwnerBookings)

		bookings.PUT("/:id/respond", h.RespondToBooking)
		bookings.PUT("/:id/cancel", h.CancelBooking)
	}
}