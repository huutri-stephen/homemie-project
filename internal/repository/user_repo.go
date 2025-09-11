package repository

import (
	"homemie/models/dto"
	"strings"
	"time"

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


func (r *UserRepository) CreateUser(user *dto.User) (err error) {
	defer func(start time.Time) {
		r.logger.Info("Create user",
			zap.String("function", "CreateUser"),
			zap.Any("params", user),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	return r.db.Create(user).Error
}

func (r *UserRepository) GetUserByEmail(email string) (user *dto.User, err error) {
	defer func(start time.Time) {
		r.logger.Info("Get user by email",
			zap.String("function", "GetUserByEmail"),
			zap.String("params", email),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	user = &dto.User{}
	if err = r.db.Where("email = ?", strings.ToLower(email)).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(id int64) (user *dto.User, err error) {
	defer func(start time.Time) {
		r.logger.Info("Get user by ID",
			zap.String("function", "GetUserByID"),
			zap.Int64("params", id),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	user = &dto.User{}
	if err = r.db.Where("id = ?", id).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(user *dto.User) (err error) {
	defer func(start time.Time) {
		r.logger.Info("Update user",
			zap.String("function", "UpdateUser"),
			zap.Any("params", user),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	return r.db.Save(user).Error
}
