package repository

import (
	"homemie/internal/domain"
	"homemie/models/dto"

	"gorm.io/gorm"
)

type favoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) domain.FavoriteRepository {
	return &favoriteRepository{db}
}

func (r *favoriteRepository) Create(favorite *dto.Favorite) error {
	return r.db.Create(favorite).Error
}

func (r *favoriteRepository) Delete(userID, listingID int64) error {
	return r.db.Where("user_id = ? AND listing_id = ?", userID, listingID).Delete(&dto.Favorite{}).Error
}

func (r *favoriteRepository) GetFavoriteListingsByUserID(userID int64) ([]*dto.Listing, error) {
	var listings []*dto.Listing
	err := r.db.
		Table("listings").
		Joins("JOIN favorites ON favorites.listing_id = listings.id").
		Where("favorites.user_id = ?", userID).
		Find(&listings).Error
	return listings, err
}

func (r *favoriteRepository) IsFavorite(userID, listingID int64) (bool, error) {
	var count int64
	err := r.db.Model(&dto.Favorite{}).Where("user_id = ? AND listing_id = ?", userID, listingID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
