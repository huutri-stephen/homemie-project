package router

import (
	"homemie/config"
	"homemie/internal/handler"
	"homemie/internal/repository"
	"homemie/internal/service"
	"homemie/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func InitUserRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg config.Config, logger *zap.Logger) {
	repo := repository.NewUserRepository(db, logger.Named("user_repo"))
	svc := service.NewUserService(repo, logger.Named("user_service"))
	h := handler.NewUserHandler(svc, logger.Named("user_handler"))

	user := rg.Group("/user")
	user.Use(utils.RequireAuth(logger))

	user.GET("/profile", h.GetUserProfile)
	user.PUT("/profile", h.UpdateUserProfile)
	user.PUT("/password", h.ChangePassword)
}
