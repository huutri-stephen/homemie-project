package router

import (
	"homemie/internal/handler"
	"homemie/internal/repo"
	"homemie/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func InitListingImageRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *zap.Logger) {
	repo := repo.NewListingImageRepository(db, logger.Named("listing_image_repo"))
	svc := service.NewListingImageService(repo, logger.Named("listing_image_service"))
	h := handler.NewListingImageHandler(svc, logger.Named("listing_image_handler"))

	listingImages := rg.Group("/listing-images")
	{
		listingImages.POST("", h.AddListingImages)
	}
}
