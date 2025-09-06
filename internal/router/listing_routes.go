package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"homemie/internal/handler"
	"homemie/internal/repository"
	"homemie/internal/service"
)

func InitListingRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *zap.Logger) {
	listingRepo := repository.NewListingRepo(db, logger.Named("listing_repo"))
	addressRepo := repository.NewAddressRepository(db, logger.Named("address_repo"))
	listingImageRepo := repository.NewListingImageRepository(db, logger.Named("listing_image_repo"))
	svc := service.NewListingService(listingRepo, addressRepo, listingImageRepo, logger.Named("listing_service"))
	h := handler.NewListingHandler(svc, logger.Named("listing_handler"))

	listings := rg.Group("/listings")
	{
		listings.POST("", h.Create)
		listings.PUT("/:id", h.Update)
		listings.DELETE("/:id", h.Delete)
	}
}

func InitPublicListingRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *zap.Logger) {
	listingRepo := repository.NewListingRepo(db, logger.Named("listing_repo"))
	addressRepo := repository.NewAddressRepository(db, logger.Named("address_repo"))
	listingImageRepo := repository.NewListingImageRepository(db, logger.Named("listing_image_repo"))
	svc := service.NewListingService(listingRepo, addressRepo, listingImageRepo, logger.Named("listing_service"))
	h := handler.NewListingHandler(svc, logger.Named("listing_handler"))

	listings := rg.Group("/listings")
	{
		listings.GET("", h.SearchAndFilter)
		listings.GET("/:id", h.GetByID)
	}
}
