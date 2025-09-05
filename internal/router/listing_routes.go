package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"homemie/internal/handler"
	"homemie/internal/repository"
	"homemie/internal/service"
)

func InitListingRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	listingRepo := repository.NewListingRepo(db)
	addressRepo := repository.NewAddressRepository(db)
	listingImageRepo := repository.NewListingImageRepository(db)
	svc := service.NewListingService(listingRepo, addressRepo, listingImageRepo)
	h := handler.NewListingHandler(svc)

	listings := rg.Group("/listings")
	{
		listings.POST("", h.Create)
		listings.PUT("/:id", h.Update)
		listings.DELETE("/:id", h.Delete)
	}
}

func InitPublicListingRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	listingRepo := repository.NewListingRepo(db)
	addressRepo := repository.NewAddressRepository(db)
	listingImageRepo := repository.NewListingImageRepository(db)
	svc := service.NewListingService(listingRepo, addressRepo, listingImageRepo)
	h := handler.NewListingHandler(svc)

	listings := rg.Group("/listings")
	{
		listings.GET("", h.SearchAndFilter)
		listings.GET("/:id", h.GetByID)
	}
}