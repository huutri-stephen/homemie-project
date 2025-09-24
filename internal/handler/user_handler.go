package handler

import (
	"homemie/db/models/dto"
	"homemie/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	svc    service.IUserService
	logger *zap.Logger
}

func NewUserHandler(svc service.IUserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{svc: svc, logger: logger}
}

func (h *UserHandler) GetUserProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.BaseResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	user, err := h.svc.GetUserProfile(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Data:    user,
	})
}

func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.BaseResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	var req dto.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := h.svc.UpdateUserProfile(userID.(int64), req); err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "User profile updated successfully",
	})
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.BaseResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := h.svc.ChangePassword(userID.(int64), req); err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}
