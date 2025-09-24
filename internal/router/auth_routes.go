package router

import (
	"homemie/config"
	"homemie/internal/handler"
	"homemie/internal/repo"
	"homemie/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InitAuthRoutes khởi tạo các route cho Auth
func InitAuthRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg config.Config, logger *zap.Logger) {
	authRepo := repo.NewAuthRepo(db, logger.Named("auth_repo"))
	userRepo := repo.NewUserRepository(db, logger.Named("user_repo"))
	svc := service.NewAuthService(authRepo, userRepo, cfg, db, logger.Named("auth_service"))
	h := handler.NewAuthHandler(svc, logger.Named("auth_handler"))

	auth := rg.Group("/auth")

	auth.POST("/signup", h.SignUp)
	auth.POST("/login", h.Login)
	auth.POST("/send-verification-email", h.SendVerificationEmail)
	auth.GET("/verify-email", h.VerifyEmail)
	auth.POST("/forgot-password", h.ForgotPassword)
	auth.POST("/reset-password", h.ResetPassword)
}
