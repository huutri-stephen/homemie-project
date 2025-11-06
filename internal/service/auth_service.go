package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"homemie/db/models"
	"homemie/db/models/dto"
	"homemie/internal/repo"
	"homemie/pkg/logger"
	"homemie/pkg/utils"
	"strings"
	"time"

	"github.com/aarondl/null/v8"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	SignUp(ctx context.Context, input dto.SignUpRequest) error
	Login(ctx context.Context, input dto.LoginRequest) (accessToken string, refreshToken string, user *models.User, err error)
	SendVerificationEmail(ctx context.Context, email string) error
	VerifyEmail(ctx context.Context, token string, email string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, input dto.ResetPasswordRequest) error
}

type authService struct {
	authRepo repo.IAuthRepository
	userRepo repo.IUserRepository
	emailTempl *utils.EmailTemplates
}

func NewAuthService(authRepo repo.IAuthRepository, userRepo repo.IUserRepository, emailTempl *utils.EmailTemplates) IAuthService {
	return &authService{authRepo: authRepo, userRepo: userRepo, emailTempl: emailTempl}
}

func (s *authService) SignUp(ctx context.Context, input dto.SignUpRequest) (err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Sign up",
			"function", "SignUp",
			"params", input,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var dateOfBirth null.Time
	if input.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", input.DateOfBirth)
		if err == nil {
			dateOfBirth = null.TimeFrom(dob)
		}
	}

	user := &models.User{
		Name:                  null.StringFrom(input.Name),
		Email:                 strings.ToLower(input.Email),
		PasswordHash:          string(hashedPassword),
		Phone:                 null.StringFrom(input.Phone),
		DateOfBirth:           dateOfBirth,
		Gender:                null.StringFrom(input.Gender),
		AvatarURL:             null.StringFrom(input.AvatarURL),
		Bio:                   null.StringFrom(input.Bio),
		UserType:              null.StringFrom(input.UserType),
		IdentityType:          null.StringFrom(input.IdentityType),
		CompanyName:           null.StringFrom(input.CompanyName),
		BusinessLicenseNumber: null.StringFrom(input.BusinessLicenseNumber),
		AgentLicenseNumber:    null.StringFrom(input.AgentLicenseNumber),
		Status:                null.StringFrom("inactive"), // Default status
		Role:                  null.StringFrom("user"),
	}

	return s.userRepo.CreateUser(ctx, user)
}

func (s *authService) Login(ctx context.Context, input dto.LoginRequest) (accessToken string, refreshToken string, user *models.User, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Login",
			"function", "Login",
			"params", input,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())

	user, err = s.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return "", "", nil, errors.New("user not found")
	}

	if user.Status.String != "active" {
		return "", "", nil, errors.New("account is not active, please verify your email")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return "", "", nil, errors.New("invalid credentials")
	}

	accessToken, refreshToken, err = utils.GenerateTokens(user)
	if err != nil {
		return "", "", nil, errors.New("could not generate tokens")
	}

	return
}

func (s *authService) SendVerificationEmail(ctx context.Context, email string) (err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Send verification email",
			"function", "SendVerificationEmail",
			"params", email,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	token, err := generateRandomToken(32)
	if err != nil {
		return err
	}

	t := &models.Token{
		UserID:    user.ID,
		TokenType: models.TokenTypeEnumEmailVerification,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err = s.authRepo.CreateToken(ctx, t); err != nil {
		return err
	}

	return s.emailTempl.SendVerificationEmail(ctx, user.Email, user.Name.String, token)
}

func (s *authService) VerifyEmail(ctx context.Context, token string, email string) (err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Verify email",
			"function", "VerifyEmail",
			"params", gin.H{"token": token, "email": email},
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	t, err := s.authRepo.GetToken(ctx, token, user.ID, models.TokenTypeEnumEmailVerification)
	if err != nil {
		return errors.New("invalid token")
	}

	if t.ExpiresAt.Before(time.Now()) {
		return errors.New("token expired")
	}

	// delete the token after verification
	if err = s.authRepo.DeleteToken(ctx, t); err != nil {
		return err
	}

	now := time.Now()
	user.Status = null.StringFrom("active")
	user.EmailVerifiedAt = null.TimeFrom(now)

	return s.userRepo.UpdateUser(ctx, user)
}

func (s *authService) ForgotPassword(ctx context.Context, email string) (err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Forgot password",
			"function", "ForgotPassword",
			"params", email,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	token, err := generateRandomToken(32)
	if err != nil {
		return err
	}

	now := time.Now()
	user.ResetPasswordToken = null.StringFrom(token)
	expiresAt := now.Add(15 * time.Minute)
	user.ResetPasswordExpiresAt = null.TimeFrom(expiresAt)

	if err = s.userRepo.UpdateUser(ctx, user); err != nil {
		return err
	}

	return s.emailTempl.SendPasswordResetEmail(ctx, user.Email, user.Name.String, token)
}

func (s *authService) ResetPassword(ctx context.Context, input dto.ResetPasswordRequest) (err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Reset password",
			"function", "ResetPassword",
			"params", input,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())

	user, err := s.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return errors.New("user not found")
	}

	if !user.ResetPasswordToken.Valid || user.ResetPasswordToken.String != input.Token {
		return errors.New("invalid token")
	}

	if user.ResetPasswordExpiresAt.Time.Before(time.Now()) {
		return errors.New("token expired")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	user.ResetPasswordToken = null.StringFromPtr(nil)
	user.ResetPasswordExpiresAt = null.TimeFromPtr(nil)

	return s.userRepo.UpdateUser(ctx, user)
}

func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}