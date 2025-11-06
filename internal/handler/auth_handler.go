package handler

import (
	"homemie/db/models/dto"
	"homemie/internal/service"
	"homemie/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc service.IAuthService
}

func NewAuthHandler(svc service.IAuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var req dto.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorw(c.Request.Context(), "Failed to bind sign up request", "error", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Infow(c.Request.Context(), "Processing sign up request", "email", req.Email)
	err := h.svc.SignUp(c.Request.Context(), req)

	if err != nil {
		logger.Errorw(c.Request.Context(), "Sign up failed", "error", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Call method SendVerificationEmail after creating user
	if err := h.svc.SendVerificationEmail(c.Request.Context(), req.Email); err != nil {
		logger.Errorw(c.Request.Context(), "Failed to send verification email after sign up", "error", err)
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Infow(c.Request.Context(), "Sign up successful", "email", req.Email)
	c.JSON(http.StatusCreated, dto.BaseResponse{
		Success: true,
		Message: "Registration successful. Please check your email to verify your account.",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorw(c.Request.Context(), "Failed to bind login request", "error", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Infow(c.Request.Context(), "Processing login request", "email", req.Email)
	accessToken, refreshToken, user, err := h.svc.Login(c.Request.Context(), req)

	if err != nil {
		logger.Errorw(c.Request.Context(), "Login failed", "error", err, "email", req.Email)
		c.JSON(http.StatusUnauthorized, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "", false, true)

	logger.Infow(c.Request.Context(), "Login successful", "email", req.Email)
	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "Login successful",
		Data: dto.LoginResponse{
			AccessToken: accessToken,
			User:        dto.NewUserPayload(user),
		},
	})
}

func (h *AuthHandler) SendVerificationEmail(c *gin.Context) {
	var req dto.SendVerificationEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorw(c.Request.Context(), "Failed to bind send verification email request", "error", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Infow(c.Request.Context(), "Processing send verification email request", "email", req.Email)
	if err := h.svc.SendVerificationEmail(c.Request.Context(), req.Email); err != nil {
		logger.Errorw(c.Request.Context(), "Failed to send verification email", "error", err, "email", req.Email)
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Infow(c.Request.Context(), "Verification email sent successfully", "email", req.Email)
	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "Verification email sent",
	})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	email := c.Query("email")

	if token == "" || email == "" {
		logger.Warnw(c.Request.Context(), "Missing token or email in verification request")
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   "Token and email are required",
		})
		return
	}

	logger.Infow(c.Request.Context(), "Processing email verification request", "email", email)
	if err := h.svc.VerifyEmail(c.Request.Context(), token, email); err != nil {
		logger.Errorw(c.Request.Context(), "Email verification failed", "error", err, "email", email)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Infow(c.Request.Context(), "Email verified successfully", "email", email)
	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "Email verified successfully",
	})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorw(c.Request.Context(), "Failed to bind forgot password request", "error", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Infow(c.Request.Context(), "Processing forgot password request", "email", req.Email)
	if err := h.svc.ForgotPassword(c.Request.Context(), req.Email); err != nil {
		logger.Errorw(c.Request.Context(), "Failed to send password reset email", "error", err, "email", req.Email)
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Infow(c.Request.Context(), "Password reset email sent successfully", "email", req.Email)
	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "Password reset email sent",
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorw(c.Request.Context(), "Failed to bind reset password request", "error", err)
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Infow(c.Request.Context(), "Processing reset password request", "email", req.Email)
	if err := h.svc.ResetPassword(c.Request.Context(), req); err != nil {
		logger.Errorw(c.Request.Context(), "Failed to reset password", "error", err, "email", req.Email)
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logger.Infow(c.Request.Context(), "Password reset successfully", "email", req.Email)
	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "Password reset successfully",
	})
}