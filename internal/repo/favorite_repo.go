package repo

import (
	"context"
	"database/sql"
	"fmt"
	"homemie/db/models"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

type IFavoriteRepository interface {
	Create(ctx context.Context, favorite *models.Favorite) error
	Delete(ctx context.Context, userID, listingID int64) error
	GetFavoriteListingsByUserID(ctx context.Context, userID int64) (models.ListingSlice, error)
	IsFavorite(ctx context.Context, userID, listingID int64) (bool, error)
}

type favoriteRepo struct {
	db *sql.DB
}

func NewFavoriteRepository(db *sql.DB) IFavoriteRepository {
	return &favoriteRepo{db}
}

func (r *favoriteRepo) Create(ctx context.Context, favorite *models.Favorite) error {
	if err := favorite.Insert(ctx, r.db, boil.Infer()); err != nil {
		return fmt.Errorf("failed to create favorite: %w", err)
	}
	return nil
}

func (r *favoriteRepo) Delete(ctx context.Context, userID, listingID int64) error {
	_, err := models.Favorites(
		models.FavoriteWhere.UserID.EQ(userID),
		models.FavoriteWhere.ListingID.EQ(listingID),
	).DeleteAll(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to delete favorite: %w", err)
	}
	return nil
}

func (r *favoriteRepo) GetFavoriteListingsByUserID(ctx context.Context, userID int64) (models.ListingSlice, error) {
	listings, err := models.Listings(
		qm.InnerJoin("favorites f ON f.listing_id = listings.id"),
		qm.Where("f.user_id = ?", userID),
	).All(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to get favorite listings: %w", err)
	}
	return listings, nil
}

func (r *favoriteRepo) IsFavorite(ctx context.Context, userID, listingID int64) (bool, error) {
	exists, err := models.Favorites(
		models.FavoriteWhere.UserID.EQ(userID),
		models.FavoriteWhere.ListingID.EQ(listingID),
	).Exists(ctx, r.db)
	if err != nil {
		return false, fmt.Errorf("failed to check if is favorite: %w", err)
	}
	return exists, nil
}
