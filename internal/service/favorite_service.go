package service

import (
	"errors"
	"homemie/db/models"
	"homemie/internal/repo"

	"gorm.io/gorm"
)

var ErrAlreadyFavorited = errors.New("listing is already favorited")

type IFavoriteService interface {
	AddToFavorites(userID, listingID int64) error
	RemoveFromFavorites(userID, listingID int64) error
	GetFavoriteListings(userID int64) ([]*models.Listing, error)
}

type favoriteService struct {
	favoriteRepo repo.IFavoriteRepository
	listingRepo  repo.IListingRepository
}

func NewFavoriteService(favoriteRepo repo.IFavoriteRepository, listingRepo repo.IListingRepository) IFavoriteService {
	return &favoriteService{favoriteRepo, listingRepo}
}

func (s *favoriteService) AddToFavorites(userID, listingID int64) error {
	_, err := s.listingRepo.FindByID(listingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("listing not found")
		}
		return err
	}

	isFav, err := s.favoriteRepo.IsFavorite(userID, listingID)
	if err != nil {
		return err
	}
	if isFav {
		return ErrAlreadyFavorited
	}

	favorite := &models.Favorite{
		UserID:    userID,
		ListingID: listingID,
	}
	return s.favoriteRepo.Create(favorite)
}

func (s *favoriteService) RemoveFromFavorites(userID, listingID int64) error {
	return s.favoriteRepo.Delete(userID, listingID)
}

func (s *favoriteService) GetFavoriteListings(userID int64) ([]*models.Listing, error) {
	return s.favoriteRepo.GetFavoriteListingsByUserID(userID)
}
