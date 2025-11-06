package router

import (
	"database/sql"
	"homemie/config"
	"homemie/internal/handler"
	"homemie/internal/repo"
	"homemie/internal/service"

	"homemie/pkg/utils"

	"github.com/gin-gonic/gin"
)

// InitAuthRoutes khởi tạo các route cho Auth
func InitAuthRoutes(rg *gin.RouterGroup, db *sql.DB, cfg *config.Config) {
	authRepo := repo.NewAuthRepo(db)
	userRepo := repo.NewUserRepository(db)
	emailTempl := utils.NewEmailTemplates(cfg, db)
	svc := service.NewAuthService(authRepo, userRepo, emailTempl)
	h := handler.NewAuthHandler(svc)

	auth := rg.Group("/auth")

	auth.POST("/signup", h.SignUp)
	auth.POST("/login", h.Login)
	auth.POST("/send-verification-email", h.SendVerificationEmail)
	auth.GET("/verify-email", h.VerifyEmail)
	auth.POST("/forgot-password", h.ForgotPassword)
	auth.POST("/reset-password", h.ResetPassword)
}