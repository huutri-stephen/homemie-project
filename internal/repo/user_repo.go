package repo

import (
	"context"
	"database/sql"
	"fmt"
	"homemie/db/models"
	"homemie/pkg/logger"
	"strings"
	"time"

	"github.com/aarondl/sqlboiler/v4/boil"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return &userRepo{db}
}

func (r *userRepo) CreateUser(ctx context.Context, user *models.User) (err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Create user",
			"function", "CreateUser",
			"params", user,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	err = user.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (user *models.User, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Get user by email",
			"function", "GetUserByEmail",
			"params", email,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	user, err = models.Users(models.UserWhere.Email.EQ(strings.ToLower(email))).One(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

func (r *userRepo) GetUserByID(ctx context.Context, id int64) (user *models.User, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Get user by ID",
			"function", "GetUserByID",
			"params", id,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	user, err = models.Users(models.UserWhere.ID.EQ(id)).One(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return user, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, user *models.User) (err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Update user",
			"function", "UpdateUser",
			"params", user,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	_, err = user.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}