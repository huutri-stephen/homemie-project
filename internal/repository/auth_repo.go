package repository

import (
	"homemie/internal/domain"
	"homemie/models/dto"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type authRepo struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewAuthRepo(db *gorm.DB, logger *zap.Logger) domain.AuthRepository {
	return &authRepo{db: db, logger: logger}
}

func (r *authRepo) CreateUser(user *dto.User) (err error) {
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

func (r *authRepo) GetUserByEmail(email string) (user *dto.User, err error) {
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

func (r *authRepo) GetUserByID(id int64) (user *dto.User, err error) {
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

func (r *authRepo) CreateToken(token *dto.Token) (err error) {
	defer func(start time.Time) {
		r.logger.Info("Create token",
			zap.String("function", "CreateToken"),
			zap.Any("params", token),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	return r.db.Create(token).Error
}

func (r *authRepo) GetToken(token string, userID int64, tokenType dto.TokenType) (t *dto.Token, err error) {
	defer func(start time.Time) {
		r.logger.Info("Get token",
			zap.String("function", "GetToken"),
			zap.Any("params", gin.H{"token": token, "user_id": userID, "token_type": tokenType}),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	t = &dto.Token{}
	if err = r.db.Where("token = ? AND user_id = ? AND token_type = ?", token, userID, tokenType).First(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (r *authRepo) UpdateUser(user *dto.User) (err error) {
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

func (r *authRepo) DeleteToken(token *dto.Token) (err error) {
	defer func(start time.Time) {
		r.logger.Info("Delete token",
			zap.String("function", "DeleteToken"),
			zap.Any("params", token),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	return r.db.Delete(token).Error
}
