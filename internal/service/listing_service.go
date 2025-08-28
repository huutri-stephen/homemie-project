package service

import (
	"errors"
	"homemie/internal/domain"
	"homemie/internal/repository"
	"homemie/models"
	"homemie/pkg/utils"
)

type ListingService struct {
	repo             domain.ListingRepository
	addressRepo      *repository.AddressRepository
	listingImageRepo *repository.ListingImageRepository
}

func NewListingService(repo domain.ListingRepository, addressRepo *repository.AddressRepository, listingImageRepo *repository.ListingImageRepository) *ListingService {
	return &ListingService{repo, addressRepo, listingImageRepo}
}

type CreateListingInput struct {
	OwnerID         uint
	Title           string
	Description     string
	PropertyType    string
	IsShared        bool
	Price           float64
	AreaM2          float64
	ContactPhone    string
	ContactEmail    string
	ContactName     string
	NumBedrooms     int
	NumBathrooms    int
	NumFloors       int
	HasBalcony      bool
	HasParking      bool
	Amenities       []string
	PetAllowed      bool
	AllowedPetTypes []string
	Latitude        float64
	Longitude       float64
	ListingType     string
	DepositAmount   float64
	Address         models.Address // todo: change to 
	Images          []models.ListingImage // todo: change to 
}


func (s *ListingService) Create(input CreateListingInput) (*models.Listing, error) {
	// Create Address
	address, err := s.addressRepo.Create(&input.Address)
	if err != nil {
		return nil, err
	}

	listing := &models.Listing{
		OwnerID:         input.OwnerID,
		Title:           input.Title,
		Description:     input.Description,
		PropertyType:    input.PropertyType,
		IsShared:        input.IsShared,
		Price:           input.Price,
		AreaM2:          input.AreaM2,
		AddressID:       address.ID,
		ContactPhone:    input.ContactPhone,
		ContactEmail:    input.ContactEmail,
		ContactName:     input.ContactName,
		NumBedrooms:     input.NumBedrooms,
		NumBathrooms:    input.NumBathrooms,
		NumFloors:       input.NumFloors,
		HasBalcony:      input.HasBalcony,
		HasParking:      input.HasParking,
		Amenities:       utils.ConvertStringArrayToJSON(input.Amenities),
		PetAllowed:      input.PetAllowed,
		AllowedPetTypes: utils.ConvertStringArrayToJSON(input.AllowedPetTypes),
		// Latitude:        input.Latitude,
		// Longitude:       input.Longitude,
		ListingType:     input.ListingType,
		DepositAmount:   input.DepositAmount,
	}

	err = s.repo.Create(listing)
	if err != nil {
		return nil, err
	}

	// Create Listing Images
	for _, image := range input.Images {
		image.ListingID = listing.ID
		_, err := s.listingImageRepo.Create(&image)
		if err != nil {
			return nil, err
		}
	}

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
	// listing.Address = input.Address
	// listing.City = input.City

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
