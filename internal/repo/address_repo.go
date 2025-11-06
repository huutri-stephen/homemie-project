package repo

import (
	"context"
	"database/sql"
	"fmt"
	"homemie/db/models"
	"homemie/pkg/logger"
	"time"

	"github.com/aarondl/sqlboiler/v4/boil"
)

type IAddressRepository interface {
	Create(ctx context.Context, address *models.Address) (*models.Address, error)
}

type addressRepo struct {
	db *sql.DB
}

func NewAddressRepository(db *sql.DB) IAddressRepository {
	return &addressRepo{db}
}

func (r *addressRepo) Create(ctx context.Context, address *models.Address) (createdAddress *models.Address, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Create address",
			"function", "Create",
			"params", address,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())

	if err = address.Insert(ctx, r.db, boil.Infer()); err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}
	return address, nil
}