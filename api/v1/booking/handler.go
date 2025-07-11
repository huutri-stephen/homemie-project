package booking

import (
	"mihome/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateBookingRequest struct {
	ListingID uint      `json:"listing_id" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

func CreateBooking(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateBookingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userIDVal, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userID := userIDVal.(uint)

		booking := models.Booking{
			UserID:    userID,
			ListingID: req.ListingID,
			StartDate: req.StartDate,
			EndDate:   req.EndDate,
			Status:    "pending",
		}

		if err := db.Create(&booking).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
			return
		}

		c.JSON(http.StatusCreated, booking)
	}
}

func GetMyBookings(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userID := userIDVal.(uint)

		var bookings []models.Booking
		if err := db.Preload("Listing").Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
			return
		}

		c.JSON(http.StatusOK, bookings)
	}
}

func GetOwnerBookings(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userID := userIDVal.(uint)

		var bookings []models.Booking
		if err := db.
			Joins("JOIN listings ON listings.id = bookings.listing_id").
			Where("listings.owner_id = ?", userID).
			Preload("User").
			Preload("Listing").
			Find(&bookings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
			return
		}

		c.JSON(http.StatusOK, bookings)
	}
}

