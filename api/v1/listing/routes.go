package listing

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func RegisterListingRoutes(rg *gin.RouterGroup, db *gorm.DB) {
    handler := NewHandler(db)

    rg.GET("/listings", handler.GetAllListings)
    rg.GET("/listings/:id", handler.GetListingByID)
    rg.POST("/listings", handler.CreateListing)
    rg.PUT("/listings/:id", handler.UpdateListing)
    rg.DELETE("/listings/:id", handler.DeleteListing)
}
