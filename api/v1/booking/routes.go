package booking

import (
	"mihome/api/v1/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterBookingRoutes(r *gin.RouterGroup, db *gorm.DB) {
	booking := r.Group("/bookings")
	booking.Use(auth.RequireAuth())

	booking.POST("/", CreateBooking(db))
	booking.GET("/my-bookings", GetMyBookings(db))
	booking.GET("/owner-bookings", GetOwnerBookings(db))
}
