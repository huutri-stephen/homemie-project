package router

import (
	"database/sql"
	"homemie/internal/handler"
	"homemie/internal/repo"
	"homemie/internal/service"
	"homemie/pkg/utils"

	"github.com/gin-gonic/gin"
)

func InitUserRoutes(rg *gin.RouterGroup, db *sql.DB) {
	repo := repo.NewUserRepository(db)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	user := rg.Group("/user")
	user.Use(utils.RequireAuth())

	user.GET("/profile", h.GetUserProfile)
	user.PUT("/profile", h.UpdateUserProfile)
	user.PUT("/password", h.ChangePassword)
}