package request

type SendVerificationEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}