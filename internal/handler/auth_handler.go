package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"homemie/internal/service"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type SignUpRequest struct {
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	Name                  string `json:"name" binding:"required"`
	Email                 string `json:"email" binding:"required,email"`
	Password              string `json:"password" binding:"required,min=6"`
	Phone                 string `json:"phone"`
	DateOfBirth           string `json:"date_of_birth"`
	Gender                string `json:"gender"`
	AvatarURL             string `json:"avatar_url"`
	Bio                   string `json:"bio"`
	UserType              string `json:"user_type" binding:"required,oneof=renter owner"`
	IdentityType          string `json:"identity_type"`
	CompanyName           string `json:"company_name"`
	BusinessLicenseNumber string `json:"business_license_number"`
	AgentLicenseNumber    string `json:"agent_license_number"`
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.SignUp(service.SignUpInput{
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call method SendVerificationEmail after creating user
	if err := h.svc.SendVerificationEmail(req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful. Please check your email to verify your account."})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, user, err := h.svc.Login(service.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"user": gin.H{
			"id":        user.ID,
			"name":      user.Name,
			"email":     user.Email,
			"role":      user.Role,
			"user_type": user.UserType,
			"status":    user.Status,
		},
	})
}

type SendVerificationEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *AuthHandler) SendVerificationEmail(c *gin.Context) {
	var req SendVerificationEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.SendVerificationEmail(req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent"})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	email := c.Query("email")

	if token == "" || email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token and email are required"})
		return
	}

	if err := h.svc.VerifyEmail(token, email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}
