package handler

import (
	"fmt"
	"homemie/internal/service"
	"homemie/models/response"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (h *MediaHandler) UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false,
			Error:   "Invalid form data",
		})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false,
			Error:   "Files are required",
		})
		return
	}

	bucketName := c.PostForm("bucketName")
	if bucketName == "" {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false,
			Error:   "Bucket name is required",
		})
		return
	}

	err = h.mediaService.CheckBucketName(bucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
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
			c.JSON(http.StatusInternalServerError, response.BaseResponse{
				Success: false,
				Error:   "Failed to open file",
			})
			return
		}
		defer fileContent.Close()

		url, err := h.mediaService.UploadFile(bucketName, objectName, fileContent, file.Size)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.BaseResponse{
				Success: false,
				Error:   "Failed to upload file",
			})
			return
		}
		urls = append(urls, url)
	}

	c.JSON(http.StatusOK, response.BaseResponse{
		Success: true,
		Data: response.MediaResponse{
			Urls: urls,
		},
	})
}
