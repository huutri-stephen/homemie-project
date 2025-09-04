package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"homemie/internal/service"
	"homemie/models/request"
	"homemie/models/response"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var req request.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false, 
			Error: err.Error(),
		})
		return
	}

	err := h.svc.SignUp(request.SignUpRequest{
		FirstName:             req.FirstName,
		LastName:              req.LastName,
		Name:                  req.Name,
		Email:                 req.Email,
		Password:              req.Password,
		Phone:                 req.Phone,
		DateOfBirth:           req.DateOfBirth,
		Gender:                req.Gender,
		AvatarURL:             req.AvatarURL,
		Bio:                   req.Bio,
		UserType:              req.UserType,
		IdentityType:          req.IdentityType,
		CompanyName:           req.CompanyName,
		BusinessLicenseNumber: req.BusinessLicenseNumber,
		AgentLicenseNumber:    req.AgentLicenseNumber,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false, 
			Error: err.Error(),
		})
		return
	}

	// Call method SendVerificationEmail after creating user
	if err := h.svc.SendVerificationEmail(req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false, 
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response.BaseResponse{
		Success: true, 
		Message: "Registration successful. Please check your email to verify your account.",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false, 
			Error: err.Error(),
		})
		return
	}

	accessToken, refreshToken, user, err := h.svc.Login(request.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, response.BaseResponse{
			Success: false, 
			Error: err.Error(),
		})
		return
	}

	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "", false, true)

	c.JSON(http.StatusOK, response.BaseResponse{
		Success: true,
		Message: "Login successful",
		Data: response.LoginResponse{
			AccessToken: accessToken,
			User: response.UserPayload{
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
	var req request.SendVerificationEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false, 
			Error: err.Error(),
		})
		return
	}

	if err := h.svc.SendVerificationEmail(req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false, 
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse{
		Success: true, 
		Message: "Verification email sent",
	})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	email := c.Query("email")

	if token == "" || email == "" {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false, 
			Error: "Token and email are required",
		})
		return
	}

	if err := h.svc.VerifyEmail(token, email); err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false, 
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse{
		Success: true, 
		Message: "Email verified successfully",
	})
}
