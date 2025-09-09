package domain

import (
	"homemie/models/dto"
)

type FavoriteRepository interface {
	Create(favorite *dto.Favorite) error
	Delete(userID, listingID int64) error
	GetFavoriteListingsByUserID(userID int64) ([]*dto.Listing, error)
	IsFavorite(userID, listingID int64) (bool, error)
}