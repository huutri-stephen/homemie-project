package handler

import (
	"homemie/internal/service"
	"homemie/models/request"
	"homemie/models/response"
	"net/http"
	"strconv"

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
			Error:   err.Error(),
		})
		return
	}
	userID := c.GetInt64("user_id")

	booking, err := h.svc.CreateBooking(userID, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false,
			Error:   "Create booking failed",
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
			Error:   "Can not fetch data",
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
			Error:   "Can not fetch data",
		})
		return
	}
	c.JSON(http.StatusOK, response.BaseResponse{Success: true, Data: bookings})
}

func (h *BookingHandler) RespondToBooking(c *gin.Context) {
	bookingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{Success: false, Error: "Invalid booking ID"})
		return
	}

	var req request.RespondBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{Success: false, Error: err.Error()})
		return
	}

	ownerID := c.GetInt64("user_id")
	booking, err := h.svc.RespondToBooking(bookingID, ownerID, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "unauthorized" {
			status = http.StatusForbidden
		} else if err.Error() == "booking not found" {
			status = http.StatusNotFound
		} else if err.Error() == "booking cannot be responded to" {
			status = http.StatusBadRequest
		}
		c.JSON(status, response.BaseResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse{Success: true, Data: booking})
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	bookingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{Success: false, Error: "Invalid booking ID"})
		return
	}

	userID := c.GetInt64("user_id")
	userRole := c.GetString("user_role")

	booking, err := h.svc.CancelBooking(bookingID, userID, userRole)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "unauthorized" {
			status = http.StatusForbidden
		} else if err.Error() == "booking not found" {
			status = http.StatusNotFound
		} else if err.Error() == "booking cannot be cancelled" {
			status = http.StatusBadRequest
		}
		c.JSON(status, response.BaseResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse{Success: true, Data: booking})
}