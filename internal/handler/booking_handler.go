package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"mihome/internal/service"
)

type BookingHandler struct {
	svc *service.BookingService
}

func NewBookingHandler(svc *service.BookingService) *BookingHandler {
	return &BookingHandler{svc}
}

type CreateBookingRequest struct {
	ListingID uint      `json:"listing_id" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetUint("user_id")

	booking, err := h.svc.CreateBooking(service.CreateBookingInput{
		UserID:    userID,
		ListingID: req.ListingID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tạo booking thất bại"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": booking})
}

func (h *BookingHandler) GetMyBookings(c *gin.Context) {
	userID := c.GetUint("user_id")
	bookings, err := h.svc.GetMyBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy dữ liệu"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": bookings})
}

func (h *BookingHandler) GetOwnerBookings(c *gin.Context) {
	userID := c.GetUint("user_id")
	bookings, err := h.svc.GetOwnerBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy dữ liệu"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": bookings})
}
