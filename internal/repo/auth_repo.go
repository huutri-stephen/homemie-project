package repo

import (
	"homemie/db/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IAuthRepository interface {
	CreateToken(token *models.Token) error
	GetToken(token string, userID int64, tokenType string) (*models.Token, error)
	DeleteToken(token *models.Token) error
}

type authRepo struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewAuthRepo(db *gorm.DB, logger *zap.Logger) IAuthRepository {
	return &authRepo{db: db, logger: logger}
}

func (r *authRepo) CreateToken(token *models.Token) (err error) {
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

func (r *authRepo) GetToken(token string, userID int64, tokenType string) (t *models.Token, err error) {
	defer func(start time.Time) {
		r.logger.Info("Get token",
			zap.String("function", "GetToken"),
			zap.Any("params", gin.H{"token": token, "user_id": userID, "token_type": tokenType}),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	t = &models.Token{}
	if err = r.db.Where("token = ? AND user_id = ? AND token_type = ?", token, userID, tokenType).First(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (r *authRepo) DeleteToken(token *models.Token) (err error) {
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
