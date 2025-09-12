package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"homemie/internal/service"
	"homemie/models/request"
	"homemie/models/response"
	"go.uber.org/zap"
)

type ListingImageHandler struct {
	service service.ListingImageService
	logger  *zap.Logger
}

func NewListingImageHandler(service service.ListingImageService, logger *zap.Logger) *ListingImageHandler {
	return &ListingImageHandler{service, logger}
}

func (h *ListingImageHandler) AddListingImages(c *gin.Context) {
	var req request.AddListingImagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{Success: false, Error: "Invalid request body"})
		return
	}

	listingImages, err := h.service.AddListingImages(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.BaseResponse{Success: false, Error: "Failed to add listing images"})
		return
	}

	c.JSON(http.StatusCreated, response.BaseResponse{Success: true, Data: listingImages})
}
