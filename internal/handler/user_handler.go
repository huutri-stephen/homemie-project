package handler

import (
	"homemie/internal/service"
	"homemie/models/request"
	"homemie/models/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserHandler struct {
	svc    *service.UserService
	logger *zap.Logger
}

func NewUserHandler(svc *service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{svc: svc, logger: logger}
}

func (h *UserHandler) GetUserProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.BaseResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	user, err := h.svc.GetUserProfile(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse{
		Success: true,
		Data:    user,
	})
}

func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.BaseResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	var req request.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := h.svc.UpdateUserProfile(userID.(uuid.UUID), req); err != nil {
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse{
		Success: true,
		Message: "User profile updated successfully",
	})
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.BaseResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := h.svc.ChangePassword(userID.(uuid.UUID), req); err != nil {
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}
