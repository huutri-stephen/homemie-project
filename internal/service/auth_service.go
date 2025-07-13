package service

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"mihome/internal/domain"
	"mihome/models"
	"mihome/pkg/utils"
)

type AuthService struct {
	repo domain.AuthRepository
}

func NewAuthService(repo domain.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

type SignUpInput struct {
	Name     string
	Email    string
	Password string
	Phone    string
}

type LoginInput struct {
	Email    string
	Password string
}

func (s *AuthService) SignUp(input SignUpInput) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Name:     input.Name,
		Email:    strings.ToLower(input.Email),
		Password: string(hashedPassword),
		Phone:    input.Phone,
		Role:     "renter",
	}

	return s.repo.CreateUser(user)
}

func (s *AuthService) Login(input LoginInput) (string, *models.User, error) {
	user, err := s.repo.GetUserByEmail(strings.ToLower(input.Email))
	if err != nil {
		return "", nil, errors.New("tài khoản không tồn tại")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", nil, errors.New("mật khẩu không đúng")
	}

	token, err := utils.GenerateJWT(*user)
	if err != nil {
		return "", nil, errors.New("không thể tạo token")
	}

	return token, user, nil
}
