package v1

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "mihome/config"
)

func NewRouter(db *gorm.DB, cfg config.Config) *gin.Engine {
    r := gin.Default()

    // Test route
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })

    return r
}
