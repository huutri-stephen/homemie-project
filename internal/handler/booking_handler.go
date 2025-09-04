package handler

import (
	"net/http"

	"homemie/internal/service"
	"homemie/models/dto"
	"homemie/models/request"
	"homemie/models/response"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	svc *service.BookingService
}

func NewBookingHandler(svc *service.BookingService) *BookingHandler {
	return &BookingHandler{svc}
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req request.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false, 
			Error: err.Error(),
		})
		return
	}
	userID := c.GetInt64("user_id")

	booking, err := h.svc.CreateBooking(dto.Booking{
		RenterID:          userID,
		ListingID:         req.ListingID,
		ScheduledTime:     req.ScheduledTime,
		MessageFromRenter: req.MessageFromRenter,
		Status:            dto.BookingStatusPending,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false, 
			Error: "Create booking failed",
		})
		return
	}

	c.JSON(http.StatusCreated, response.BaseResponse{Success: true, Data: booking})
}

func (h *BookingHandler) GetMyBookings(c *gin.Context) {
	userID := c.GetInt64("user_id")
	bookings, err := h.svc.GetMyBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false, 
			Error: "Can not fetch data",
		})
		return
	}
	c.JSON(http.StatusOK, response.BaseResponse{Success: true, Data: bookings})
}

func (h *BookingHandler) GetOwnerBookings(c *gin.Context) {
	userID := c.GetInt64("user_id")
	bookings, err := h.svc.GetOwnerBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false, 
			Error: "Can not fetch data",
		})
		return
	}
	c.JSON(http.StatusOK, response.BaseResponse{Success: true, Data: bookings})
}
