package service

import (
	"context"
	"database/sql"
	"errors"
	"homemie/db/models"
	"homemie/internal/repo"
)

var ErrAlreadyFavorited = errors.New("listing is already favorited")

type IFavoriteService interface {
	AddToFavorites(ctx context.Context, userID, listingID int64) error
	RemoveFromFavorites(ctx context.Context, userID, listingID int64) error
	GetFavoriteListings(ctx context.Context, userID int64) (models.ListingSlice, error)
}

type favoriteService struct {
	favoriteRepo repo.IFavoriteRepository
	listingRepo  repo.IListingRepository
}

func NewFavoriteService(favoriteRepo repo.IFavoriteRepository, listingRepo repo.IListingRepository) IFavoriteService {
	return &favoriteService{favoriteRepo, listingRepo}
}

func (s *favoriteService) AddToFavorites(ctx context.Context, userID, listingID int64) error {
	_, err := s.listingRepo.FindByID(ctx, listingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("listing not found")
		}
		return err
	}

	isFav, err := s.favoriteRepo.IsFavorite(ctx, userID, listingID)
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
	return s.favoriteRepo.Create(ctx, favorite)
}

func (s *favoriteService) RemoveFromFavorites(ctx context.Context, userID, listingID int64) error {
	return s.favoriteRepo.Delete(ctx, userID, listingID)
}

func (s *favoriteService) GetFavoriteListings(ctx context.Context, userID int64) (models.ListingSlice, error) {
	return s.favoriteRepo.GetFavoriteListingsByUserID(ctx, userID)
}
