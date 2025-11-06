package internal

import (
	"database/sql"
	"homemie/config"
	"homemie/internal/router"
	"homemie/pkg/utils"

	// "github.com/aws/aws-sdk-go-v2/service/s3"
	// "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(db *sql.DB, cfg *config.Config, externalEndpoint string) *gin.Engine {
	// func NewRouter(db *sql.DB, cfg config.Config, logger *zap.Logger, s3Client *s3.Client, externalEndpoint string) *gin.Engine {
	r := gin.New()        // Use gin.New() instead of gin.Default() to have more control over middleware
	r.Use(gin.Recovery()) // Add recovery middleware
	r.Use(utils.StructuredLogger())

	// r.Use(cors.New(cors.Config{
	//     AllowOrigins:     []string{"http://localhost:5173"}, // Vite mặc định port 5173
	//     AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	//     AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	//     ExposeHeaders:    []string{"Content-Length"},
	//     AllowCredentials: true,
	// }))

	api := r.Group(cfg.Server.ApiVersion)

	router.InitAuthRoutes(api, db, cfg)
	router.InitPublicListingRoutes(api, db)

	// Protected routes (require JWT)
	protected := api.Group("/")
	protected.Use(utils.RequireAuth())

	router.InitListingRoutes(protected, db)
	router.InitBookingRoutes(protected, db)
	router.InitUserRoutes(protected, db)
	router.InitFavoriteRoutes(protected, db)
	// router.InitMediaRoutes(protected, db, logger, s3Client, externalEndpoint)
	router.InitListingImageRoutes(protected, db)

	return r
}