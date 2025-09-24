package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"homemie/internal/service"
	"homemie/db/models/dto"
)

type AuthHandler struct {
	svc    service.IAuthService
	logger *zap.Logger
}

func NewAuthHandler(svc service.IAuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{svc: svc, logger: logger}
}

func (h *AuthHandler) getLogger(c *gin.Context) *zap.Logger {
	if logger, exists := c.Get("logger"); exists {
		if zapLogger, ok := logger.(*zap.Logger); ok {
			return zapLogger
		}
	}
	return h.logger
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	logger := h.getLogger(c)
	var req dto.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind sign up request", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Info("Processing sign up request", zap.String("email", req.Email))
	err := h.svc.SignUp(req)

	if err != nil {
		logger.Error("Sign up failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Call method SendVerificationEmail after creating user
	if err := h.svc.SendVerificationEmail(req.Email); err != nil {
		logger.Error("Failed to send verification email after sign up", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Info("Sign up successful", zap.String("email", req.Email))
	c.JSON(http.StatusCreated, dto.BaseResponse{
		Success: true,
		Message: "Registration successful. Please check your email to verify your account.",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	logger := h.getLogger(c)
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Info("Processing login request", zap.String("email", req.Email))
	accessToken, refreshToken, user, err := h.svc.Login(req)

	if err != nil {
		logger.Error("Login failed", zap.Error(err), zap.String("email", req.Email))
		c.JSON(http.StatusUnauthorized, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "", false, true)

	logger.Info("Login successful", zap.String("email", req.Email))
	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "Login successful",
		Data: dto.LoginResponse{
			AccessToken: accessToken,
			User: dto.UserPayload{
				ID:       user.ID,
				Name:     user.Name,
				Email:    user.Email,
				Role:     user.Role,
				UserType: user.UserType,
				Status:   user.Status,
			},
		},
	})
}

func (h *AuthHandler) SendVerificationEmail(c *gin.Context) {
	logger := h.getLogger(c)
	var req dto.SendVerificationEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind send verification email request", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Info("Processing send verification email request", zap.String("email", req.Email))
	if err := h.svc.SendVerificationEmail(req.Email); err != nil {
		logger.Error("Failed to send verification email", zap.Error(err), zap.String("email", req.Email))
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Info("Verification email sent successfully", zap.String("email", req.Email))
	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "Verification email sent",
	})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	logger := h.getLogger(c)
	token := c.Query("token")
	email := c.Query("email")

	if token == "" || email == "" {
		logger.Warn("Missing token or email in verification request")
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   "Token and email are required",
		})
		return
	}

	logger.Info("Processing email verification request", zap.String("email", email))
	if err := h.svc.VerifyEmail(token, email); err != nil {
		logger.Error("Email verification failed", zap.Error(err), zap.String("email", email))
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Info("Email verified successfully", zap.String("email", email))
	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "Email verified successfully",
	})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	logger := h.getLogger(c)
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind forgot password request", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Info("Processing forgot password request", zap.String("email", req.Email))
	if err := h.svc.ForgotPassword(req.Email); err != nil {
		logger.Error("Failed to send password reset email", zap.Error(err), zap.String("email", req.Email))
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Info("Password reset email sent successfully", zap.String("email", req.Email))
	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "Password reset email sent",
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	logger := h.getLogger(c)
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind reset password request", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Info("Processing reset password request", zap.String("email", req.Email))
	if err := h.svc.ResetPassword(req); err != nil {
		logger.Error("Failed to reset password", zap.Error(err), zap.String("email", req.Email))
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Info("Password reset successfully", zap.String("email", req.Email))
	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "Password reset successfully",
	})
}
