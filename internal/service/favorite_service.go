package service

import (
	"github.com/google/uuid"

	"homemie/internal/domain"
	"homemie/internal/repository"
	"homemie/models/request"
)

type FavoriteService interface {
	Create(userID uuid.UUID, req request.Favorite) error
	GetByUserID(userID uuid.UUID, page int, pageSize int) ([]domain.Favorite, error)
}

type favoriteService struct {
	favoriteRepo *repository.favoriteRepository
	listingRepo  *repository.listingRepository
}

func NewFavoriteService(favoriteRepo *repository.favoriteRepository, listingRepo *repository.listingRepository) FavoriteService {
	return &favoriteService{favoriteRepo, listingRepo}
}

func (s *favoriteService) Create(userID uuid.UUID, req request.Favorite) error {
	// Check if listing exists
	_, err := s.listingRepo.GetByID(req.ListingID)
	if err != nil {
		return err
	}

	// Check if favorite already exists
	_, err = s.favoriteRepo.FindByUserIDAndListingID(userID, req.ListingID)
	if err == nil {
		return nil
	}

	favorite := &domain.Favorite{
		ID:        uuid.New(),
		UserID:    userID,
		ListingID: req.ListingID,
	}
	return s.favoriteRepo.Create(favorite)
}

func (s *favoriteService) GetByUserID(userID uuid.UUID, page int, pageSize int) ([]domain.Favorite, error) {
	return s.favoriteRepo.GetByUserID(userID, page, pageSize)
}
