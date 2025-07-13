package internal

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "mihome/internal/router"

    "mihome/config"
    "mihome/pkg/utils"
)

func NewRouter(db *gorm.DB, cfg config.Config) *gin.Engine {
    r := gin.Default()
    api := r.Group("/api/v1")

    router.InitAuthRoutes(api, db) 

    // Protected routes (require JWT)
    protected := api.Group("/")
    protected.Use(utils.RequireAuth())

    router.InitListingRoutes(protected, db)
    router.InitBookingRoutes(protected, db)

    return r
}

