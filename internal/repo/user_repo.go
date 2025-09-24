package repo

import (
	"homemie/db/models"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int64) (*models.User, error)
	UpdateUser(user *models.User) error
}

type userRepo struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserRepository(db *gorm.DB, logger *zap.Logger) IUserRepository {
	return &userRepo{db, logger}
}

func (r *userRepo) CreateUser(user *models.User) (err error) {
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

func (r *userRepo) GetUserByEmail(email string) (user *models.User, err error) {
	defer func(start time.Time) {
		r.logger.Info("Get user by email",
			zap.String("function", "GetUserByEmail"),
			zap.String("params", email),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	user = &models.User{}
	if err = r.db.Where("email = ?", strings.ToLower(email)).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) GetUserByID(id int64) (user *models.User, err error) {
	defer func(start time.Time) {
		r.logger.Info("Get user by ID",
			zap.String("function", "GetUserByID"),
			zap.Int64("params", id),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	user = &models.User{}
	if err = r.db.Where("id = ?", id).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) UpdateUser(user *models.User) (err error) {
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
