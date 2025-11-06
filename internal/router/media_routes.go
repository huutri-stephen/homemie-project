package router

import (
	"database/sql"
	"homemie/internal/handler"
	"homemie/internal/service"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func InitMediaRoutes(rg *gin.RouterGroup, db *sql.DB, s3Client *s3.Client, externalEndpoint string) {
	mediaService := service.NewMediaService(s3Client, externalEndpoint)
	mediaHandler := handler.NewMediaHandler(mediaService)

	media := rg.Group("/media")
	{
		media.POST("/presigned-url", mediaHandler.GeneratePresignedUploadURL)
		media.POST("/upload", mediaHandler.UploadFiles)
	}
}