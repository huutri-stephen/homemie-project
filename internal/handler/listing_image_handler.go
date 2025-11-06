package handler

import (
	"homemie/db/models/dto"
	"homemie/internal/service"
	"homemie/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ListingImageHandler struct {
	service service.IListingImageService
}

func NewListingImageHandler(service service.IListingImageService) *ListingImageHandler {
	return &ListingImageHandler{service}
}

func (h *ListingImageHandler) AddListingImages(c *gin.Context) {
	var req dto.AddListingImagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorw(c.Request.Context(), "Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Success: false, Error: "Invalid request body"})
		return
	}

	listingImages, err := h.service.AddListingImages(c.Request.Context(), req)
	if err != nil {
		logger.Errorw(c.Request.Context(), "Failed to add listing images", "error", err)
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Success: false, Error: "Failed to add listing images"})
		return
	}

	c.JSON(http.StatusCreated, dto.BaseResponse{Success: true, Data: listingImages})
}