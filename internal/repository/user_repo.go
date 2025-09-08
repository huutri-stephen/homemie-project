package repository

import (
	"homemie/models/dto"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserRepository(db *gorm.DB, logger *zap.Logger) *UserRepository {
	return &UserRepository{db, logger}
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (user *dto.User, err error) {
	defer func(start time.Time) {
		r.logger.Info("Get user by id",
			zap.String("function", "GetUserByID"),
			zap.Any("id", id),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	if err = r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(user *dto.User) (err error) {
	defer func(start time.Time) {
		r.logger.Info("Update user",
			zap.String("function", "UpdateUser"),
			zap.Any("user", user),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	if err = r.db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}
