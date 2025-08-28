package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"homemie/config"
	"homemie/internal/handler"
	"homemie/internal/repository"
	"homemie/internal/service"
)

// InitAuthRoutes khởi tạo các route cho Auth
func InitAuthRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg config.Config) {
	repo := repository.NewAuthRepo(db)
	svc := service.NewAuthService(repo, cfg, db)
	h := handler.NewAuthHandler(svc)

	auth := rg.Group("/auth")

	auth.POST("/signup", h.SignUp)
	auth.POST("/login", h.Login)
	auth.POST("/send-verification-email", h.SendVerificationEmail)
	auth.GET("/verify-email", h.VerifyEmail)
}
