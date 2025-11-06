package handler

import (
	"homemie/db/models/dto"
	"homemie/internal/service"
	"homemie/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ListingHandler struct {
	svc service.IListingService
}

func NewListingHandler(svc service.IListingService) *ListingHandler {
	return &ListingHandler{svc}
}

func (h *ListingHandler) Create(c *gin.Context) {
	var req dto.CreateListingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorw(c.Request.Context(), "Failed to bind create listing request", "error", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID := c.GetInt64("user_id")
	req.OwnerID = userID

	logger.Info(c.Request.Context(), "Processing create listing request")
	listing, err := h.svc.Create(c.Request.Context(), req)

	if err != nil {
		logger.Errorw(c.Request.Context(), "Failed to create listing", "error", err)
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Success: false,
			Error:   "Create listing failed",
		})
		return
	}

	logger.Infow(c.Request.Context(), "Successfully created listing", "listing_id", listing.ID)
	c.JSON(http.StatusCreated, dto.BaseResponse{Success: true, Data: listing})
}

func (h *ListingHandler) SearchAndFilter(c *gin.Context) {
	var filter dto.SearchFilterListing
	if err := c.ShouldBindQuery(&filter); err != nil {
		logger.Errorw(c.Request.Context(), "Failed to bind search and filter query", "error", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Info(c.Request.Context(), "Processing search and filter request")
	listings, pagination, err := h.svc.SearchAndFilter(c.Request.Context(), &filter)
	if err != nil {
		logger.Errorw(c.Request.Context(), "Failed to search and filter listings", "error", err)
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Success: false,
			Error:   "Get list failed",
		})
		return
	}

	logger.Info(c.Request.Context(), "Successfully retrieved listings")
	c.JSON(http.StatusOK, dto.BaseResponse{Success: true, Data: gin.H{
		"listings":   listings,
		"pagination": pagination,
	}})
}

func (h *ListingHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Errorw(c.Request.Context(), "Invalid listing ID format", "error", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Success: false, Error: "Invalid ID format"})
		return
	}

	logger.Infow(c.Request.Context(), "Processing get listing by ID request", "id", id)
	listing, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		logger.Errorw(c.Request.Context(), "Failed to get listing by ID", "error", err, "id", id)
		c.JSON(http.StatusNotFound, dto.BaseResponse{
			Success: false,
			Error:   "Not found",
		})
		return
	}

	logger.Infow(c.Request.Context(), "Successfully retrieved listing by ID", "id", id)
	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Data:    listing,
	})
}

func (h *ListingHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Errorw(c.Request.Context(), "Invalid listing ID format", "error", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Success: false, Error: "Invalid ID format"})
		return
	}

	userID := c.GetInt64("user_id")

	var req dto.CreateListingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorw(c.Request.Context(), "Failed to bind update listing request", "error", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Infow(c.Request.Context(), "Processing update listing request", "id", id)
	listing, err := h.svc.Update(c.Request.Context(), id, userID, req)
	if err != nil {
		logger.Errorw(c.Request.Context(), "Failed to update listing", "error", err, "id", id)
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, dto.BaseResponse{
				Success: false,
				Error:   "Unauthorized to edit this listing",
			})
		} else {
			c.JSON(http.StatusInternalServerError, dto.BaseResponse{
				Success: false,
				Error:   "Update failed",
			})
		}
		return
	}

	logger.Infow(c.Request.Context(), "Successfully updated listing", "id", id)
	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Data:    listing,
	})
}

func (h *ListingHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Errorw(c.Request.Context(), "Invalid listing ID format", "error", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Success: false, Error: "Invalid ID format"})
		return
	}

	userID := c.GetInt64("user_id")

	logger.Infow(c.Request.Context(), "Processing delete listing request", "id", id)
	err = h.svc.Delete(c.Request.Context(), int64(id), userID)
	if err != nil {
		logger.Errorw(c.Request.Context(), "Failed to delete listing", "error", err, "id", id)
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, dto.BaseResponse{
				Success: false,
				Error:   "Unauthorized to delete this listing",
			})
		} else {
			c.JSON(http.StatusInternalServerError, dto.BaseResponse{
				Success: false,
				Error:   "Delete failed",
			})
		}
		return
	}

	logger.Infow(c.Request.Context(), "Successfully deleted listing", "id", id)
	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "Delete listing successfully",
	})
}