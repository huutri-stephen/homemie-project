package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"homemie/config"
	"homemie/db/models"
	"homemie/db/models/dto"
	"homemie/internal/repo"
	"homemie/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IAuthService interface {
	SignUp(input dto.SignUpRequest) error
	Login(input dto.LoginRequest) (accessToken string, refreshToken string, user *models.User, err error)
	SendVerificationEmail(email string) error
	VerifyEmail(token string, email string) error
	ForgotPassword(email string) error
	ResetPassword(input dto.ResetPasswordRequest) error
}

type authService struct {
	authRepo repo.IAuthRepository
	userRepo repo.IUserRepository
	Cfg      config.Config
	DB       *gorm.DB
	logger   *zap.Logger
}

func NewAuthService(authRepo repo.IAuthRepository, userRepo repo.IUserRepository, cfg config.Config, db *gorm.DB, logger *zap.Logger) IAuthService {
	return &authService{authRepo: authRepo, userRepo: userRepo, Cfg: cfg, DB: db, logger: logger}
}

func (s *authService) SignUp(input dto.SignUpRequest) (err error) {
	defer func(start time.Time) {
		s.logger.Info("Sign up",
			zap.String("function", "SignUp"),
			zap.Any("params", input),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var dateOfBirth *time.Time
	if input.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", input.DateOfBirth)
		if err == nil {
			dateOfBirth = &dob
		}
	}

	user := &models.User{
		FirstName:             input.FirstName,
		LastName:              input.LastName,
		Name:                  input.Name,
		Email:                 strings.ToLower(input.Email),
		PasswordHash:          string(hashedPassword),
		Phone:                 input.Phone,
		DateOfBirth:           dateOfBirth,
		Gender:                input.Gender,
		AvatarURL:             input.AvatarURL,
		Bio:                   input.Bio,
		UserType:              input.UserType,
		IdentityType:          input.IdentityType,
		CompanyName:           input.CompanyName,
		BusinessLicenseNumber: input.BusinessLicenseNumber,
		AgentLicenseNumber:    input.AgentLicenseNumber,
		Status:                "inactive", // Default status
		Role:                  "user",
	}

	return s.userRepo.CreateUser(user)
}

func (s *authService) Login(input dto.LoginRequest) (accessToken string, refreshToken string, user *models.User, err error) {
	defer func(start time.Time) {
		s.logger.Info("Login",
			zap.String("function", "Login"),
			zap.Any("params", input),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	user, err = s.userRepo.GetUserByEmail(input.Email)
	if err != nil {
		return "", "", nil, errors.New("user not found")
	}

	if user.Status != "active" {
		return "", "", nil, errors.New("account is not active, please verify your email")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return "", "", nil, errors.New("invalid credentials")
	}

	accessToken, refreshToken, err = utils.GenerateTokens(*user)
	if err != nil {
		return "", "", nil, errors.New("could not generate tokens")
	}

	return
}

func (s *authService) SendVerificationEmail(email string) (err error) {
	defer func(start time.Time) {
		s.logger.Info("Send verification email",
			zap.String("function", "SendVerificationEmail"),
			zap.String("params", email),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	user, err := s.userRepo.GetUserByEmail(email)
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

	if err = s.authRepo.CreateToken(t); err != nil {
		return err
	}

	return utils.SendVerificationEmail(s.Cfg, s.DB, user.Email, user.Name, token)
}

func (s *authService) VerifyEmail(token string, email string) (err error) {
	defer func(start time.Time) {
		s.logger.Info("Verify email",
			zap.String("function", "VerifyEmail"),
			zap.Any("params", gin.H{"token": token, "email": email}),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	t, err := s.authRepo.GetToken(token, user.ID, models.TokenTypeEnumEmailVerification)
	if err != nil {
		return errors.New("invalid token")
	}

	if t.ExpiresAt.Before(time.Now()) {
		return errors.New("token expired")
	}

	// delete the token after verification
	if err = s.authRepo.DeleteToken(t); err != nil {
		return err
	}

	now := time.Now()
	user.Status = "active"
	user.EmailVerifiedAt = &now

	return s.userRepo.UpdateUser(user)
}

func (s *authService) ForgotPassword(email string) (err error) {
	defer func(start time.Time) {
		s.logger.Info("Forgot password",
			zap.String("function", "ForgotPassword"),
			zap.String("params", email),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	token, err := generateRandomToken(32)
	if err != nil {
		return err
	}

	now := time.Now()
	user.ResetPasswordToken = token
	expiresAt := now.Add(15 * time.Minute)
	user.ResetPasswordExpiresAt = &expiresAt

	if err = s.userRepo.UpdateUser(user); err != nil {
		return err
	}

	return utils.SendPasswordResetEmail(s.Cfg, s.DB, user.Email, user.Name, token)
}

func (s *authService) ResetPassword(input dto.ResetPasswordRequest) (err error) {
	defer func(start time.Time) {
		s.logger.Info("Reset password",
			zap.String("function", "ResetPassword"),
			zap.Any("params", input),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	user, err := s.userRepo.GetUserByEmail(input.Email)
	if err != nil {
		return errors.New("user not found")
	}

	if user.ResetPasswordToken.String == "" || user.ResetPasswordToken.String != input.Token {
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
	user.ResetPasswordToken = ""
	user.ResetPasswordExpiresAt = nil

	return s.userRepo.UpdateUser(user)
}

func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
