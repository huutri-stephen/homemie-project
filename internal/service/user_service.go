package service

import (
	"homemie/internal/domain"
	"homemie/models/dto"
	"homemie/models/request"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService struct {
	repo   domain.UserRepository
	logger *zap.Logger
}

func NewUserService(repo domain.UserRepository, logger *zap.Logger) *UserService {
	return &UserService{repo: repo, logger: logger}
}

func (s *UserService) GetUserProfile(id uuid.UUID) (user *dto.User, err error) {
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

func (s *UserService) UpdateUserProfile(id uuid.UUID, req request.UpdateUserProfileRequest) (err error) {
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

