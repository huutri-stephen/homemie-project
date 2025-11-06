package repo

import (
	"context"
	"database/sql"
	"fmt"
	"homemie/db/models"
	"homemie/pkg/logger"

	"github.com/aarondl/sqlboiler/v4/boil"
)

type IListingImageRepository interface {
	AddListingImages(ctx context.Context, listingImages []models.ListingImage) ([]models.ListingImage, error)
}

type listingImageRepo struct {
	db *sql.DB
}

func NewListingImageRepository(db *sql.DB) IListingImageRepository {
	return &listingImageRepo{db}
}

func (r *listingImageRepo) AddListingImages(ctx context.Context, listingImages []models.ListingImage) ([]models.ListingImage, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.FromContext(ctx).Errorw("starting transaction failed", "error", err)
		return nil, fmt.Errorf("starting transaction failed: %w", err)
	}
	defer tx.Rollback()

	for i := range listingImages {
		err = (&listingImages[i]).Insert(ctx, tx, boil.Infer())
		if err != nil {
			logger.FromContext(ctx).Errorw("inserting listing image failed", "error", err)
			return nil, fmt.Errorf("inserting listing image failed: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		logger.FromContext(ctx).Errorw("committing transaction failed", "error", err)
		return nil, fmt.Errorf("committing transaction failed: %w", err)
	}

	return listingImages, nil
}