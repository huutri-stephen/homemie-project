package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"homemie/internal/handler"
	"homemie/internal/repository"
	"homemie/internal/service"
)

func InitListingRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	repo := repository.NewListingRepo(db)
	svc := service.NewListingService(repo)
	h := handler.NewListingHandler(svc)

	listings := rg.Group("/listings")
	{
		listings.POST("", h.Create)
		listings.GET("", h.GetAll)
		listings.GET("/:id", h.GetByID)
		listings.PUT("/:id", h.Update)
		listings.DELETE("/:id", h.Delete)
	}
}
