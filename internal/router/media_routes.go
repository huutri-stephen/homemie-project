package router

import (
	"homemie/internal/handler"
	"homemie/internal/service"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func InitMediaRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *zap.Logger, s3Client *s3.Client) {
	mediaService := service.NewMediaService(s3Client, logger)
	mediaHandler := handler.NewMediaHandler(mediaService, logger)

	media := rg.Group("/media")
	{
		media.POST("/presigned-url", mediaHandler.GeneratePresignedUploadURL)
	}
}
