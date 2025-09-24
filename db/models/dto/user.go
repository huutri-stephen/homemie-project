package dto

import (
	"homemie/db/models"
	"time"
)

// Request DTOs
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

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type UpdateUserProfileRequest struct {
	FirstName             string     `json:"first_name,omitempty"`
	LastName              string     `json:"last_name,omitempty"`
	Name                  string     `json:"name,omitempty"`
	Phone                 string     `json:"phone,omitempty"`
	DateOfBirth           *time.Time `json:"date_of_birth,omitempty"`
	Gender                string     `json:"gender,omitempty"` // "male", "female", "other"
	AvatarURL             string     `json:"avatar_url,omitempty"`
	Bio                   string     `json:"bio,omitempty"`
	CompanyName           string     `json:"company_name,omitempty"`
	BusinessLicenseNumber string     `json:"business_license_number,omitempty"`
	AgentLicenseNumber    string     `json:"agent_license_number,omitempty"`
	IdentityType          string     `json:"identity_type,omitempty"` // "personal", "company"
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type SendVerificationEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// Response DTOs
type LoginResponse struct {
	AccessToken string      `json:"access_token"`
	User        UserPayload `json:"user"`
}

type UserPayload struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	UserType string `json:"user_type"`
	Status   string `json:"status"`
}

func NewUserPayload(user *models.User) UserPayload {
	return UserPayload{
		ID:       user.ID,
		Name:     user.Name.String,
		Email:    user.Email,
		Role:     user.Role.String,
		UserType: user.UserType.String,
		Status:   user.Status.String,
	}
}