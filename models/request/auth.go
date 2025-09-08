package request

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
