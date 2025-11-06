package router

import (
	"homemie/internal/handler"
	"homemie/internal/repository"
	"homemie/internal/service"
	"homemie/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func InitFavoriteRoutes(r *gin.RouterGroup, db *gorm.DB, logger *zap.Logger) {
	favoriteRepo := repository.NewFavoriteRepository(db)
	listingRepo := repository.NewListingRepo(db)
	favoriteService := service.NewFavoriteService(favoriteRepo, listingRepo)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService)

	favorite := r.Group("/favorites")
	favorite.Use(utils.AuthMiddleware())
	{
		favorite.POST("", favoriteHandler.Create)
		favorite.GET("", favoriteHandler.GetByUserID)
	}
}
