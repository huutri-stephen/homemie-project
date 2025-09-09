package service

import (
	"errors"
	"homemie/internal/domain"
	"homemie/models/dto"

	"gorm.io/gorm"
)

var ErrAlreadyFavorited = errors.New("listing is already favorited")

type FavoriteService interface {
	AddToFavorites(userID, listingID int64) error
	RemoveFromFavorites(userID, listingID int64) error
	GetFavoriteListings(userID int64) ([]*dto.Listing, error)
}

type favoriteService struct {
	favoriteRepo domain.FavoriteRepository
	listingRepo  domain.ListingRepository
}

func NewFavoriteService(favoriteRepo domain.FavoriteRepository, listingRepo domain.ListingRepository) FavoriteService {
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

	favorite := &dto.Favorite{
		UserID:    userID,
		ListingID: listingID,
	}
	return s.favoriteRepo.Create(favorite)
}

func (s *favoriteService) RemoveFromFavorites(userID, listingID int64) error {
	return s.favoriteRepo.Delete(userID, listingID)
}

func (s *favoriteService) GetFavoriteListings(userID int64) ([]*dto.Listing, error) {
	return s.favoriteRepo.GetFavoriteListingsByUserID(userID)
}
