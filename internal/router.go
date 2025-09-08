package internal

import (
	"homemie/config"
	"homemie/internal/router"
	"homemie/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, cfg config.Config, logger *zap.Logger) *gin.Engine {
	r := gin.New()        // Use gin.New() instead of gin.Default() to have more control over middleware
	r.Use(gin.Recovery()) // Add recovery middleware
	r.Use(utils.StructuredLogger(logger))

	api := r.Group(cfg.Server.ApiVersion)

	router.InitAuthRoutes(api, db, cfg, logger)
	router.InitPublicListingRoutes(api, db, logger)

	// Protected routes (require JWT)
	protected := api.Group("/")
	protected.Use(utils.RequireAuth(logger))

	router.InitListingRoutes(protected, db, logger)
	router.InitBookingRoutes(protected, db, logger)
	router.InitUserRoutes(protected, db, cfg, logger)

	return r
}
