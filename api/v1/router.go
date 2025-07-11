package v1

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "mihome/api/v1/auth"
    "mihome/api/v1/listing"
    "mihome/api/v1/booking"

    "mihome/config"
)

func NewRouter(db *gorm.DB, cfg config.Config) *gin.Engine {
    r := gin.Default()
    api := r.Group("/api/v1")

    // Auth routes (public)
    auth.RegisterAuthRoutes(api.Group("/auth"), db)

    // Protected routes (require JWT)
    protected := api.Group("/")
    protected.Use(auth.RequireAuth())

    listing.RegisterListingRoutes(protected, db)
    booking.RegisterBookingRoutes(protected, db)

    return r
}
