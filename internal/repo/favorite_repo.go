package repo

import (
	"homemie/db/models"

	"gorm.io/gorm"
)

type IFavoriteRepository interface {
	Create(favorite *models.Favorite) error
	Delete(userID, listingID int64) error
	GetFavoriteListingsByUserID(userID int64) ([]*models.Listing, error)
	IsFavorite(userID, listingID int64) (bool, error)
}

type favoriteRepo struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) IFavoriteRepository {
	return &favoriteRepo{db}
}

func (r *favoriteRepo) Create(favorite *models.Favorite) error {
	return r.db.Create(favorite).Error
}

func (r *favoriteRepo) Delete(userID, listingID int64) error {
	return r.db.Where("user_id = ? AND listing_id = ?", userID, listingID).Delete(&models.Favorite{}).Error
}

func (r *favoriteRepo) GetFavoriteListingsByUserID(userID int64) ([]*models.Listing, error) {
	var listings []*models.Listing
	err := r.db.
		Table("listings").
		Joins("JOIN favorites ON favorites.listing_id = listings.id").
		Where("favorites.user_id = ?", userID).
		Find(&listings).Error
	return listings, err
}

func (r *favoriteRepo) IsFavorite(userID, listingID int64) (bool, error) {
	var count int64
	err := r.db.Model(&models.Favorite{}).Where("user_id = ? AND listing_id = ?", userID, listingID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
