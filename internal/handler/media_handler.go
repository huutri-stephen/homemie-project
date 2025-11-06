package handler

import (
	"fmt"
	"homemie/db/models/dto"
	"homemie/internal/service"
	"homemie/pkg/logger"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MediaHandler struct {
	mediaService service.IMediaService
}

func NewMediaHandler(mediaService service.IMediaService) *MediaHandler {
	return &MediaHandler{
		mediaService: mediaService,
	}
}

type GeneratePresignedURLRequest struct {
	ObjectName string `json:"objectName" binding:"required"`
	BucketName string `json:"bucketName" binding:"required"`
}

func (h *MediaHandler) GeneratePresignedUploadURL(c *gin.Context) {
	var req GeneratePresignedURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorw(c.Request.Context(), "Failed to bind generate presigned url request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, err := h.mediaService.GeneratePresignedUploadURL(c.Request.Context(), req.BucketName, req.ObjectName)
	if err != nil {
		logger.Errorw(c.Request.Context(), "Failed to generate presigned URL", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate presigned URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *MediaHandler) UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		logger.Errorw(c.Request.Context(), "Invalid form data", "error", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   "Invalid form data",
		})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		logger.Warnw(c.Request.Context(), "Files are required")
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   "Files are required",
		})
		return
	}

	bucketName := c.PostForm("bucketName")
	if bucketName == "" {
		logger.Warnw(c.Request.Context(), "Bucket name is required")
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   "Bucket name is required",
		})
		return
	}

	err = h.mediaService.CheckBucketName(c.Request.Context(), bucketName)
	if err != nil {
		logger.Errorw(c.Request.Context(), fmt.Sprintf("Bucket %s does not exist", bucketName), "error", err)
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Success: false,
			Error:   fmt.Sprintf("Bucket %s does not exist", bucketName),
		})
		return
	}

	var urls []string
	for _, file := range files {
		ext := filepath.Ext(file.Filename)
		objectName := uuid.New().String() + ext
		fileContent, err := file.Open()
		if err != nil {
			logger.Errorw(c.Request.Context(), "Failed to open file", "error", err)
			c.JSON(http.StatusInternalServerError, dto.BaseResponse{
				Success: false,
				Error:   "Failed to open file",
			})
			return
		}
		defer fileContent.Close()

		url, err := h.mediaService.UploadFile(c.Request.Context(), bucketName, objectName, fileContent, file.Size)
		if err != nil {
			logger.Errorw(c.Request.Context(), "Failed to upload file", "error", err)
			c.JSON(http.StatusInternalServerError, dto.BaseResponse{
				Success: false,
				Error:   "Failed to upload file",
			})
			return
		}
		urls = append(urls, url)
	}

	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Data: dto.MediaResponse{
			Urls: urls,
		},
	})
}