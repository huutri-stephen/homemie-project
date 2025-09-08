package request

import "time"

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

