package handler

import (
	"net/http"

	"homemie/db/models/dto"
	"homemie/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ListingImageHandler struct {
	service service.IListingImageService
	logger  *zap.Logger
}

func NewListingImageHandler(service service.IListingImageService, logger *zap.Logger) *ListingImageHandler {
	return &ListingImageHandler{service, logger}
}

func (h *ListingImageHandler) AddListingImages(c *gin.Context) {
	var req dto.AddListingImagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Success: false, Error: "Invalid request body"})
		return
	}

	listingImages, err := h.service.AddListingImages(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Success: false, Error: "Failed to add listing images"})
		return
	}

	c.JSON(http.StatusCreated, dto.BaseResponse{Success: true, Data: listingImages})
}
