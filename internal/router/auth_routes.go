package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mihome/internal/handler"
	"mihome/internal/repository"
	"mihome/internal/service"
)

// InitAuthRoutes khởi tạo các route cho Auth
func InitAuthRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	repo := repository.NewAuthRepo(db)
	svc := service.NewAuthService(repo)
	h := handler.NewAuthHandler(svc)

	auth := rg.Group("/auth")

	auth.POST("/signup", h.SignUp)
	auth.POST("/login", h.Login)
}
