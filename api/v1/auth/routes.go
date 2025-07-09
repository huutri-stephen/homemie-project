package auth

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, db *gorm.DB) {
    handler := NewAuthHandler(db)

    rg.POST("/signup", handler.SignUp)
    rg.POST("/login", handler.Login)
}
