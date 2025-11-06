package router

import (
	"database/sql"
	"homemie/internal/handler"
	"homemie/internal/repo"
	"homemie/internal/service"

	"github.com/gin-gonic/gin"
)

func InitListingImageRoutes(rg *gin.RouterGroup, db *sql.DB) {
	repo := repo.NewListingImageRepository(db)
	svc := service.NewListingImageService(repo)
	h := handler.NewListingImageHandler(svc)

	listingImages := rg.Group("/listing-images")
	{
		listingImages.POST("", h.AddListingImages)
	}
}