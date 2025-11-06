package handler

import (
	"homemie/db/models/dto"
	"homemie/internal/service"
	"homemie/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	svc service.IBookingService
}

func NewBookingHandler(svc service.IBookingService) *BookingHandler {
	return &BookingHandler{svc}
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req dto.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorw(c.Request.Context(), "Failed to bind create booking request", "error", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	userID := c.GetInt64("user_id")

	logger.Info(c.Request.Context(), "Processing create booking request")
	booking, err := h.svc.CreateBooking(c.Request.Context(), userID, req)

	if err != nil {
		logger.Errorw(c.Request.Context(), "Failed to create booking", "error", err)
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Success: false,
			Error:   "Create booking failed",
		})
		return
	}

	logger.Infow(c.Request.Context(), "Successfully created booking", "booking_id", booking.ID)
	c.JSON(http.StatusCreated, dto.BaseResponse{Success: true, Data: booking})
}

func (h *BookingHandler) GetMyBookings(c *gin.Context) {
	userID := c.GetInt64("user_id")

	logger.Info(c.Request.Context(), "Processing get my bookings request")
	bookings, err := h.svc.GetMyBookings(c.Request.Context(), userID)
	if err != nil {
		logger.Errorw(c.Request.Context(), "Failed to get my bookings", "error", err)
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Success: false,
			Error:   "Can not fetch data",
		})
		return
	}

	logger.Info(c.Request.Context(), "Successfully retrieved my bookings")
	c.JSON(http.StatusOK, dto.BaseResponse{Success: true, Data: bookings})
}

func (h *BookingHandler) GetOwnerBookings(c *gin.Context) {
	userID := c.GetInt64("user_id")

	logger.Info(c.Request.Context(), "Processing get owner bookings request")
	bookings, err := h.svc.GetOwnerBookings(c.Request.Context(), userID)
	if err != nil {
		logger.Errorw(c.Request.Context(), "Failed to get owner bookings", "error", err)
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Success: false,
			Error:   "Can not fetch data",
		})
		return
	}

	logger.Info(c.Request.Context(), "Successfully retrieved owner bookings")
	c.JSON(http.StatusOK, dto.BaseResponse{Success: true, Data: bookings})
}

func (h *BookingHandler) RespondToBooking(c *gin.Context) {
	bookingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logger.Errorw(c.Request.Context(), "Invalid booking ID format", "error", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Success: false, Error: "Invalid booking ID"})
		return
	}

	var req dto.RespondBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorw(c.Request.Context(), "Failed to bind respond booking request", "error", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Success: false, Error: err.Error()})
		return
	}

	ownerID := c.GetInt64("user_id")
	logger.Infow(c.Request.Context(), "Processing respond to booking request", "booking_id", bookingID, "status", req.Status)
	booking, err := h.svc.RespondToBooking(c.Request.Context(), bookingID, ownerID, req)
	if err != nil {
		logger.Errorw(c.Request.Context(), "Failed to respond to booking", "error", err, "booking_id", bookingID)
		status := http.StatusInternalServerError
		if err.Error() == "unauthorized" {
			status = http.StatusForbidden
		} else if err.Error() == "booking not found" {
			status = http.StatusNotFound
		} else if err.Error() == "booking cannot be responded to" {
			status = http.StatusBadRequest
		}
		c.JSON(status, dto.BaseResponse{Success: false, Error: err.Error()})
		return
	}

	logger.Infow(c.Request.Context(), "Successfully responded to booking", "booking_id", bookingID)
	c.JSON(http.StatusOK, dto.BaseResponse{Success: true, Data: booking})
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	bookingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logger.Errorw(c.Request.Context(), "Invalid booking ID format", "error", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Success: false, Error: "Invalid booking ID"})
		return
	}

	userID := c.GetInt64("user_id")
	userRole := c.GetString("user_role")

	logger.Infow(c.Request.Context(), "Processing cancel booking request", "booking_id", bookingID)
	booking, err := h.svc.CancelBooking(c.Request.Context(), bookingID, userID, userRole)
	if err != nil {
		logger.Errorw(c.Request.Context(), "Failed to cancel booking", "error", err, "booking_id", bookingID)
		status := http.StatusInternalServerError
		if err.Error() == "unauthorized" {
			status = http.StatusForbidden
		} else if err.Error() == "booking not found" {
			status = http.StatusNotFound
		} else if err.Error() == "booking cannot be cancelled" {
			status = http.StatusBadRequest
		}
		c.JSON(status, dto.BaseResponse{Success: false, Error: err.Error()})
		return
	}

	logger.Infow(c.Request.Context(), "Successfully cancelled booking", "booking_id", bookingID)
	c.JSON(http.StatusOK, dto.BaseResponse{Success: true, Data: booking})
}