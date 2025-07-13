package service

import (
	"errors"
	"mihome/internal/domain"
	"mihome/models"
)

type ListingService struct {
	repo domain.ListingRepository
}

func NewListingService(repo domain.ListingRepository) *ListingService {
	return &ListingService{repo}
}

type CreateListingInput struct {
	Title       string
	Description string
	Price       float64
	Address     string
	City        string
	OwnerID     uint
}

func (s *ListingService) Create(input CreateListingInput) (*models.Listing, error) {
	listing := &models.Listing{
		Title:       input.Title,
		Description: input.Description,
		Price:       input.Price,
		Address:     input.Address,
		City:        input.City,
		OwnerID:     input.OwnerID,
	}
	err := s.repo.Create(listing)
	return listing, err
}

func (s *ListingService) GetAll() ([]models.Listing, error) {
	return s.repo.FindAll()
}

func (s *ListingService) GetByID(id uint) (*models.Listing, error) {
	return s.repo.FindByID(id)
}

func (s *ListingService) Update(id uint, userID uint, input CreateListingInput) (*models.Listing, error) {
	listing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if listing.OwnerID != userID {
		return nil, errors.New("unauthorized")
	}

	listing.Title = input.Title
	listing.Description = input.Description
	listing.Price = input.Price
	listing.Address = input.Address
	listing.City = input.City

	err = s.repo.Update(listing)
	return listing, err
}

func (s *ListingService) Delete(id uint, userID uint) error {
	listing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if listing.OwnerID != userID {
		return errors.New("unauthorized")
	}
	return s.repo.Delete(listing)
}
