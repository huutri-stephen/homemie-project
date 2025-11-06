package service

import (
	"context"
	"errors"
	"homemie/db/models"
	"homemie/db/models/dto"
	"homemie/internal/repo"
	"homemie/pkg/logger"
	"time"

	"github.com/aarondl/null/v8"
	"golang.org/x/crypto/bcrypt"
)

//create interface
type IUserService interface {
	GetUserProfile(ctx context.Context, id int64) (*models.User, error)
	UpdateUserProfile(ctx context.Context, id int64, req dto.UpdateUserProfileRequest) error
	ChangePassword(ctx context.Context, id int64, req dto.ChangePasswordRequest) error
}

type userService struct {
	repo repo.IUserRepository
}

func NewUserService(repo repo.IUserRepository) IUserService {
	return &userService{repo: repo}
}

func (s *userService) GetUserProfile(ctx context.Context, id int64) (user *models.User, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Get user profile",
			"function", "GetUserProfile",
			"id", id,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())

	user, err = s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUserProfile(ctx context.Context, id int64, req dto.UpdateUserProfileRequest) (err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Update user profile",
			"function", "UpdateUserProfile",
			"id", id,
			"req", req,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())

	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	// Profile cơ bản
	if req.Name != "" {
		user.Name = null.StringFrom(req.Name)
	}
	if req.Phone != "" {
		user.Phone = null.StringFrom(req.Phone)
	}
	if req.DateOfBirth != nil {
		user.DateOfBirth = null.TimeFrom(*req.DateOfBirth)
	}
	if req.Gender != "" {
		user.Gender = null.StringFrom(req.Gender)
	}
	if req.AvatarURL != "" {
		user.AvatarURL = null.StringFrom(req.AvatarURL)
	}
	if req.Bio != "" {
		user.Bio = null.StringFrom(req.Bio)
	}

	// Dành cho loại account owner/agent/business
	if req.CompanyName != "" {
		user.CompanyName = null.StringFrom(req.CompanyName)
	}
	if req.BusinessLicenseNumber != "" {
		user.BusinessLicenseNumber = null.StringFrom(req.BusinessLicenseNumber)
	}
	if req.AgentLicenseNumber != "" {
		user.AgentLicenseNumber = null.StringFrom(req.AgentLicenseNumber)
	}
	if req.IdentityType != "" {
		user.IdentityType = null.StringFrom(req.IdentityType)
	}

	// Persist
	if err = s.repo.UpdateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *userService) ChangePassword(ctx context.Context, id int64, req dto.ChangePasswordRequest) (err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Change password",
			"function", "ChangePassword",
			"id", id,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())

	user, err := s.repo.GetUserByID(ctx, id)
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

	return s.repo.UpdateUser(ctx, user)
}