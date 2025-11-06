package router

import (
	"database/sql"
	"homemie/internal/handler"
	"homemie/internal/repo"
	"homemie/internal/service"
	"homemie/pkg/utils"

	"github.com/gin-gonic/gin"
)

func InitFavoriteRoutes(r *gin.RouterGroup, db *sql.DB) {
	favoriteRepo := repo.NewFavoriteRepository(db)
	listingRepo := repo.NewListingRepo(db)
	favoriteService := service.NewFavoriteService(favoriteRepo, listingRepo)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService)

	favorite := r.Group("/favorites")
	favorite.Use(utils.RequireAuth())
	{
		favorite.POST("", favoriteHandler.AddToFavorites)
		favorite.DELETE("/:listing_id", favoriteHandler.RemoveFromFavorites)
		favorite.GET("", favoriteHandler.GetFavoriteListings)
	}
}