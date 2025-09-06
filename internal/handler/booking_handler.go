package handler

import (
	"homemie/internal/service"
	"homemie/models/request"
	"homemie/models/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BookingHandler struct {
	svc    *service.BookingService
	logger *zap.Logger
}

func NewBookingHandler(svc *service.BookingService, logger *zap.Logger) *BookingHandler {
	return &BookingHandler{svc, logger}
}

func (h *BookingHandler) getLogger(c *gin.Context) *zap.Logger {
	if logger, exists := c.Get("logger"); exists {
		if zapLogger, ok := logger.(*zap.Logger); ok {
			return zapLogger
		}
	}
	return h.logger
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	logger := h.getLogger(c)
	var req request.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind create booking request", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	userID := c.GetInt64("user_id")

	logger.Info("Processing create booking request")
	booking, err := h.svc.CreateBooking(userID, req)

	if err != nil {
		logger.Error("Failed to create booking", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false,
			Error:   "Create booking failed",
		})
		return
	}

	logger.Info("Successfully created booking", zap.Int64("booking_id", booking.ID))
	c.JSON(http.StatusCreated, response.BaseResponse{Success: true, Data: booking})
}

func (h *BookingHandler) GetMyBookings(c *gin.Context) {
	logger := h.getLogger(c)
	userID := c.GetInt64("user_id")

	logger.Info("Processing get my bookings request")
	bookings, err := h.svc.GetMyBookings(userID)
	if err != nil {
		logger.Error("Failed to get my bookings", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false,
			Error:   "Can not fetch data",
		})
		return
	}

	logger.Info("Successfully retrieved my bookings")
	c.JSON(http.StatusOK, response.BaseResponse{Success: true, Data: bookings})
}

func (h *BookingHandler) GetOwnerBookings(c *gin.Context) {
	logger := h.getLogger(c)
	userID := c.GetInt64("user_id")

	logger.Info("Processing get owner bookings request")
	bookings, err := h.svc.GetOwnerBookings(userID)
	if err != nil {
		logger.Error("Failed to get owner bookings", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false,
			Error:   "Can not fetch data",
		})
		return
	}

	logger.Info("Successfully retrieved owner bookings")
	c.JSON(http.StatusOK, response.BaseResponse{Success: true, Data: bookings})
}

func (h *BookingHandler) RespondToBooking(c *gin.Context) {
	logger := h.getLogger(c)
	bookingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logger.Error("Invalid booking ID format", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.BaseResponse{Success: false, Error: "Invalid booking ID"})
		return
	}

	var req request.RespondBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind respond booking request", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.BaseResponse{Success: false, Error: err.Error()})
		return
	}

	ownerID := c.GetInt64("user_id")
	logger.Info("Processing respond to booking request", zap.Int64("booking_id", bookingID), zap.String("status", req.Status))
	booking, err := h.svc.RespondToBooking(bookingID, ownerID, req)
	if err != nil {
		logger.Error("Failed to respond to booking", zap.Error(err), zap.Int64("booking_id", bookingID))
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

	logger.Info("Successfully responded to booking", zap.Int64("booking_id", bookingID))
	c.JSON(http.StatusOK, response.BaseResponse{Success: true, Data: booking})
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	logger := h.getLogger(c)
	bookingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logger.Error("Invalid booking ID format", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.BaseResponse{Success: false, Error: "Invalid booking ID"})
		return
	}

	userID := c.GetInt64("user_id")
	userRole := c.GetString("user_role")

	logger.Info("Processing cancel booking request", zap.Int64("booking_id", bookingID))
	booking, err := h.svc.CancelBooking(bookingID, userID, userRole)
	if err != nil {
		logger.Error("Failed to cancel booking", zap.Error(err), zap.Int64("booking_id", bookingID))
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

	logger.Info("Successfully cancelled booking", zap.Int64("booking_id", bookingID))
	c.JSON(http.StatusOK, response.BaseResponse{Success: true, Data: booking})
}
