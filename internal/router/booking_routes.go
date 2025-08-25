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
	repo := repository.NewBookingRepo(db)
	svc := service.NewBookingService(repo)
	h := handler.NewBookingHandler(svc)

	bookings := rg.Group("/bookings")
	{
		bookings.POST("", h.CreateBooking)
		bookings.GET("/my-bookings", h.GetMyBookings)
		bookings.GET("/owner-bookings", h.GetOwnerBookings)

		// Optional: các API mở rộng về sau
		// bookings.POST("/:id/approve", h.ApproveBooking)
		// bookings.POST("/:id/reject", h.RejectBooking)
		// bookings.POST("/:id/cancel", h.CancelBooking)
	}
}
