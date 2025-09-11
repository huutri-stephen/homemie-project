package handler

import (
	"homemie/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MediaHandler struct {
	mediaService service.MediaService
	logger       *zap.Logger
}

func NewMediaHandler(mediaService service.MediaService, logger *zap.Logger) *MediaHandler {
	return &MediaHandler{
		mediaService: mediaService,
		logger:       logger,
	}
}

type GeneratePresignedURLRequest struct {
	ObjectName string `json:"objectName" binding:"required"`
	BucketName string `json:"bucketName" binding:"required"`
}

func (h *MediaHandler) GeneratePresignedUploadURL(c *gin.Context) {
	var req GeneratePresignedURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, err := h.mediaService.GeneratePresignedUploadURL(req.BucketName, req.ObjectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate presigned URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}
