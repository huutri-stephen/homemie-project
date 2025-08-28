package internal

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "homemie/internal/router"

    "homemie/config"
    "homemie/pkg/utils"
)

func NewRouter(db *gorm.DB, cfg config.Config) *gin.Engine {
    r := gin.Default()
    api := r.Group(cfg.Server.ApiVersion)

    router.InitAuthRoutes(api, db, cfg) 

    // Protected routes (require JWT)
    protected := api.Group("/")
    protected.Use(utils.RequireAuth())

    router.InitListingRoutes(protected, db)
    router.InitBookingRoutes(protected, db)

    return r
}

