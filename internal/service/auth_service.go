package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"homemie/config"
	"homemie/internal/domain"
	"homemie/models/dto"
	"homemie/models/request"
	"homemie/pkg/utils"
)

type AuthService struct {
	repo   domain.AuthRepository
	Cfg    config.Config
	DB     *gorm.DB
	logger *zap.Logger
}

func NewAuthService(repo domain.AuthRepository, cfg config.Config, db *gorm.DB, logger *zap.Logger) *AuthService {
	return &AuthService{repo: repo, Cfg: cfg, DB: db, logger: logger}
}

func (s *AuthService) SignUp(input request.SignUpRequest) (err error) {
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

	user := &dto.User{
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

	return s.repo.CreateUser(user)
}

func (s *AuthService) Login(input request.LoginRequest) (accessToken string, refreshToken string, user *dto.User, err error) {
	defer func(start time.Time) {
		s.logger.Info("Login",
			zap.String("function", "Login"),
			zap.Any("params", input),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	user, err = s.repo.GetUserByEmail(input.Email)
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

func (s *AuthService) SendVerificationEmail(email string) (err error) {
	defer func(start time.Time) {
		s.logger.Info("Send verification email",
			zap.String("function", "SendVerificationEmail"),
			zap.String("params", email),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	token, err := generateRandomToken(32)
	if err != nil {
		return err
	}

	t := &dto.Token{
		UserID:    user.ID,
		TokenType: dto.EmailVerification,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err = s.repo.CreateToken(t); err != nil {
		return err
	}

	return utils.SendVerificationEmail(s.Cfg, s.DB, user.Email, user.Name, token)
}

func (s *AuthService) VerifyEmail(token string, email string) (err error) {
	defer func(start time.Time) {
		s.logger.Info("Verify email",
			zap.String("function", "VerifyEmail"),
			zap.Any("params", gin.H{"token": token, "email": email}),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	t, err := s.repo.GetToken(token, user.ID, dto.EmailVerification)
	if err != nil {
		return errors.New("invalid token")
	}

	if t.ExpiresAt.Before(time.Now()) {
		return errors.New("token expired")
	}

	// delete the token after verification
	if err = s.repo.DeleteToken(t); err != nil {
		return err
	}

	now := time.Now()
	user.Status = "active"
	user.EmailVerifiedAt = &now

	return s.repo.UpdateUser(user)
}

func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
