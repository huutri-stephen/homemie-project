package router

import (
	"homemie/internal/handler"
	"homemie/internal/repo"
	"homemie/internal/service"
	"homemie/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func InitFavoriteRoutes(r *gin.RouterGroup, db *gorm.DB, logger *zap.Logger) {
	favoriteRepo := repo.NewFavoriteRepository(db)
	listingRepo := repo.NewListingRepo(db, logger)
	favoriteService := service.NewFavoriteService(favoriteRepo, listingRepo)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService)

	favorite := r.Group("/favorites")
	favorite.Use(utils.RequireAuth(logger)) // Assuming you have an auth middleware
	{
		favorite.POST("", favoriteHandler.AddToFavorites)
		favorite.DELETE("/:listing_id", favoriteHandler.RemoveFromFavorites)
		favorite.GET("", favoriteHandler.GetFavoriteListings)
	}
}
