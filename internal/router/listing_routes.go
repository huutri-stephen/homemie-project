package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"homemie/internal/handler"
	"homemie/internal/repo"
	"homemie/internal/service"
)

func InitListingRoutes(rg *gin.RouterGroup, db *sql.DB) {
	listingRepo := repo.NewListingRepo(db)
	addressRepo := repo.NewAddressRepository(db)
	listingImageRepo := repo.NewListingImageRepository(db)
	svc := service.NewListingService(listingRepo, addressRepo, listingImageRepo)
	h := handler.NewListingHandler(svc)

	listings := rg.Group("/listings")
	{
		listings.POST("", h.Create)
		listings.PUT("/:id", h.Update)
		listings.DELETE("/:id", h.Delete)
	}
}

func InitPublicListingRoutes(rg *gin.RouterGroup, db *sql.DB) {
	listingRepo := repo.NewListingRepo(db)
	addressRepo := repo.NewAddressRepository(db)
	listingImageRepo := repo.NewListingImageRepository(db)
	svc := service.NewListingService(listingRepo, addressRepo, listingImageRepo)
	h := handler.NewListingHandler(svc)

	listings := rg.Group("/listings")
	{
		listings.GET("", h.SearchAndFilter)
		listings.GET("/:id", h.GetByID)
	}
}