package repo

import (
	"context"
	"database/sql"
	"fmt"
	"homemie/db/models"
	"homemie/pkg/logger"
	"time"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/gin-gonic/gin"
)

type IAuthRepository interface {
	CreateToken(ctx context.Context, token *models.Token) error
	GetToken(ctx context.Context, token string, userID int64, tokenType string) (*models.Token, error)
	DeleteToken(ctx context.Context, token *models.Token) error
}

type authRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) IAuthRepository {
	return &authRepo{db: db}
}

func (r *authRepo) CreateToken(ctx context.Context, token *models.Token) (err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Create token",
			"function", "CreateToken",
			"params", token,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	err = token.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return fmt.Errorf("failed to create token: %w", err)
	}
	return nil
}

func (r *authRepo) GetToken(ctx context.Context, token string, userID int64, tokenType string) (t *models.Token, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Get token",
			"function", "GetToken",
			"params", gin.H{"token": token, "user_id": userID, "token_type": tokenType},
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	t, err = models.Tokens(
		models.TokenWhere.Token.EQ(token),
		models.TokenWhere.UserID.EQ(userID),
		models.TokenWhere.TokenType.EQ(tokenType),
	).One(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	return t, nil
}

func (r *authRepo) DeleteToken(ctx context.Context, token *models.Token) (err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Delete token",
			"function", "DeleteToken",
			"params", token,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	_, err = token.Delete(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}
	return nil
}