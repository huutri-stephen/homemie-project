package handler

import (
	"errors"
	"homemie/db/models/dto"
	"homemie/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	service service.IFavoriteService
}

func NewFavoriteHandler(service service.IFavoriteService) *FavoriteHandler {
	return &FavoriteHandler{service}
}

func (h *FavoriteHandler) AddToFavorites(c *gin.Context) {
	var req dto.AddFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   "Invalid request payload",
		})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.BaseResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	if err := h.service.AddToFavorites(userID.(int64), req.ListingID); err != nil {
		if errors.Is(err, service.ErrAlreadyFavorited) {
			c.JSON(http.StatusConflict, dto.BaseResponse{
				Success: false,
				Error:   "Listing is already in favorites",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Success: false,
			Error:   "Failed to add to favorites",
		})
		return
	}

	c.JSON(http.StatusCreated, dto.BaseResponse{
		Success: true,
		Data:    "Added to favorites",
	})
}

func (h *FavoriteHandler) RemoveFromFavorites(c *gin.Context) {
	listingID, err := strconv.ParseInt(c.Param("listing_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   "Invalid listing ID",
		})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.BaseResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	if err := h.service.RemoveFromFavorites(userID.(int64), listingID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Success: false,
			Error:   "Failed to remove from favorites",
		})
		return
	}

	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Data:    "Removed from favorites",
	})
}

func (h *FavoriteHandler) GetFavoriteListings(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.BaseResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	listings, err := h.service.GetFavoriteListings(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Success: false,
			Error:   "Failed to retrieve favorite listings",
		})
		return
	}

	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Data:    listings,
	})
}
