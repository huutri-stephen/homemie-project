package service

import (
	"errors"
	"homemie/db/models"
	"homemie/db/models/dto"
	"homemie/internal/repo"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

//create interface
type IUserService interface {
	GetUserProfile(id int64) (*models.User, error)
	UpdateUserProfile(id int64, req dto.UpdateUserProfileRequest) error
	ChangePassword(id int64, req dto.ChangePasswordRequest) error
}

type userService struct {
	repo   repo.IUserRepository
	logger *zap.Logger
}

func NewUserService(repo repo.IUserRepository, logger *zap.Logger) IUserService {
	return &userService{repo: repo, logger: logger}
}

func (s *userService) GetUserProfile(id int64) (user *models.User, err error) {
	defer func(start time.Time) {
		s.logger.Info("Get user profile",
			zap.String("function", "GetUserProfile"),
			zap.Any("id", id),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	user, err = s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUserProfile(id int64, req dto.UpdateUserProfileRequest) (err error) {
	defer func(start time.Time) {
		s.logger.Info("Update user profile",
			zap.String("function", "UpdateUserProfile"),
			zap.Any("id", id),
			zap.Any("req", req),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return err
	}

	// Profile cơ bản
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.DateOfBirth != nil {
		user.DateOfBirth = req.DateOfBirth
	}
	if req.Gender != "" {
		user.Gender = req.Gender
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}

	// Dành cho loại account owner/agent/business
	if req.CompanyName != "" {
		user.CompanyName = req.CompanyName
	}
	if req.BusinessLicenseNumber != "" {
		user.BusinessLicenseNumber = req.BusinessLicenseNumber
	}
	if req.AgentLicenseNumber != "" {
		user.AgentLicenseNumber = req.AgentLicenseNumber
	}
	if req.IdentityType != "" {
		user.IdentityType = req.IdentityType
	}

	// Persist
	if err = s.repo.UpdateUser(user); err != nil {
		return err
	}

	return nil
}

func (s *userService) ChangePassword(id int64, req dto.ChangePasswordRequest) (err error) {
	defer func(start time.Time) {
		s.logger.Info("Change password",
			zap.String("function", "ChangePassword"),
			zap.Any("id", id),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return errors.New("invalid credentials")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)

	return s.repo.UpdateUser(user)
}
