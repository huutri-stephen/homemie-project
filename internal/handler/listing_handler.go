package handler

import (
	"homemie/internal/service"
	"homemie/models/dto"
	"homemie/models/request"
	"homemie/models/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ListingHandler struct {
	svc    *service.ListingService
	logger *zap.Logger
}

func NewListingHandler(svc *service.ListingService, logger *zap.Logger) *ListingHandler {
	return &ListingHandler{svc, logger}
}

func (h *ListingHandler) getLogger(c *gin.Context) *zap.Logger {
	if logger, exists := c.Get("logger"); exists {
		if zapLogger, ok := logger.(*zap.Logger); ok {
			return zapLogger
		}
	}
	return h.logger
}

func (h *ListingHandler) Create(c *gin.Context) {
	logger := h.getLogger(c)
	var req request.CreateListingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind create listing request", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID := c.GetInt64("user_id")
	req.OwnerID = userID

	logger.Info("Processing create listing request")
	listing, err := h.svc.Create(req)

	if err != nil {
		logger.Error("Failed to create listing", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false,
			Error:   "Create listing failed",
		})
		return
	}

	logger.Info("Successfully created listing", zap.Int64("listing_id", listing.ID))
	c.JSON(http.StatusCreated, response.BaseResponse{Success: true, Data: listing})
}

func (h *ListingHandler) SearchAndFilter(c *gin.Context) {
	logger := h.getLogger(c)
	var filter dto.SearchFilterListing
	if err := c.ShouldBindQuery(&filter); err != nil {
		logger.Error("Failed to bind search and filter query", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Info("Processing search and filter request")
	listings, pagination, err := h.svc.SearchAndFilter(&filter)
	if err != nil {
		logger.Error("Failed to search and filter listings", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false,
			Error:   "Get list failed",
		})
		return
	}

	logger.Info("Successfully retrieved listings")
	c.JSON(http.StatusOK, response.BaseResponse{Success: true, Data: gin.H{
		"listings":   listings,
		"pagination": pagination,
	}})
}

func (h *ListingHandler) GetByID(c *gin.Context) {
	logger := h.getLogger(c)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("Invalid listing ID format", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.BaseResponse{Success: false, Error: "Invalid ID format"})
		return
	}

	logger.Info("Processing get listing by ID request", zap.Int64("id", id))
	listing, err := h.svc.GetByID(id)
	if err != nil {
		logger.Error("Failed to get listing by ID", zap.Error(err), zap.Int64("id", id))
		c.JSON(http.StatusNotFound, response.BaseResponse{
			Success: false,
			Error:   "Not found",
		})
		return
	}

	logger.Info("Successfully retrieved listing by ID", zap.Int64("id", id))
	c.JSON(http.StatusOK, response.BaseResponse{
		Success: true,
		Data:    listing,
	})
}

func (h *ListingHandler) Update(c *gin.Context) {
	logger := h.getLogger(c)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("Invalid listing ID format", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.BaseResponse{Success: false, Error: "Invalid ID format"})
		return
	}

	userID := c.GetInt64("user_id")

	var req request.CreateListingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind update listing request", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Info("Processing update listing request", zap.Int64("id", id))
	listing, err := h.svc.Update(id, userID, req)
	if err != nil {
		logger.Error("Failed to update listing", zap.Error(err), zap.Int64("id", id))
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, response.BaseResponse{
				Success: false,
				Error:   "Unauthorized to edit this listing",
			})
		} else {
			c.JSON(http.StatusInternalServerError, response.BaseResponse{
				Success: false,
				Error:   "Update failed",
			})
		}
		return
	}

	logger.Info("Successfully updated listing", zap.Int64("id", id))
	c.JSON(http.StatusOK, response.BaseResponse{
		Success: true,
		Data:    listing,
	})
}

func (h *ListingHandler) Delete(c *gin.Context) {
	logger := h.getLogger(c)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("Invalid listing ID format", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.BaseResponse{Success: false, Error: "Invalid ID format"})
		return
	}

	userID := c.GetInt64("user_id")

	logger.Info("Processing delete listing request", zap.Int64("id", id))
	err = h.svc.Delete(int64(id), userID)
	if err != nil {
		logger.Error("Failed to delete listing", zap.Error(err), zap.Int64("id", id))
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, response.BaseResponse{
				Success: false,
				Error:   "Unauthorized to delete this listing",
			})
		} else {
			c.JSON(http.StatusInternalServerError, response.BaseResponse{
				Success: false,
				Error:   "Delete failed",
			})
		}
		return
	}

	logger.Info("Successfully deleted listing", zap.Int64("id", id))
	c.JSON(http.StatusOK, response.BaseResponse{
		Success: true,
		Message: "Delete listing successfully",
	})
}
