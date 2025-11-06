package handler

import (
	"homemie/db/models/dto"
	"homemie/internal/service"
	"homemie/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc service.IUserService
}

func NewUserHandler(svc service.IUserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetUserProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		logger.Warnw(c.Request.Context(), "Unauthorized")
		c.JSON(http.StatusUnauthorized, dto.BaseResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	user, err := h.svc.GetUserProfile(c.Request.Context(), userID.(int64))
	if err != nil {
		logger.Errorw(c.Request.Context(), "Failed to get user profile", "error", err)
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
		logger.Warnw(c.Request.Context(), "Unauthorized")
		c.JSON(http.StatusUnauthorized, dto.BaseResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	var req dto.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorw(c.Request.Context(), "Failed to bind update user profile request", "error", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := h.svc.UpdateUserProfile(c.Request.Context(), userID.(int64), req); err != nil {
		logger.Errorw(c.Request.Context(), "Failed to update user profile", "error", err)
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
		logger.Warnw(c.Request.Context(), "Unauthorized")
		c.JSON(http.StatusUnauthorized, dto.BaseResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorw(c.Request.Context(), "Failed to bind change password request", "error", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := h.svc.ChangePassword(c.Request.Context(), userID.(int64), req); err != nil {
		logger.Errorw(c.Request.Context(), "Failed to change password", "error", err)
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