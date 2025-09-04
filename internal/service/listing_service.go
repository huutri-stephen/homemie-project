package service

import (
	"errors"
	"homemie/internal/domain"
	"homemie/internal/repository"
	"homemie/models/request"
	"homemie/models/dto"
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

func (s *ListingService) Create(input request.CreateListingRequest) (*dto.Listing, error) {
	
	// todo: find current addrest if address exist
	addr := dto.Address{
		CityID:       input.Address.CityID,
		WardID:       input.Address.WardID,
		AreaID:       input.Address.AreaID,
		Street:       input.Address.Street,
		HouseNumber:  input.Address.HouseNumber,
		BuildingName: input.Address.BuildingName,
		FloorNumber:  input.Address.FloorNumber,
		RoomNumber:   input.Address.RoomNumber,
		Latitude:     input.Address.Latitude,
		Longitude:    input.Address.Longitude,
	}
	address, err := s.addressRepo.Create(&addr)
	if err != nil {
		return nil, err
	}

	listing := &dto.Listing{
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

	for _, image := range input.Images {
		img := dto.ListingImage{
			ListingID: listing.ID,
			ImageURL:  image.ImageURL,
			IsMain:    image.IsMain,
			SortOrder: image.SortOrder,
		}
		_, err := s.listingImageRepo.Create(&img)
		if err != nil {
			return nil, err
		}
	}

	return listing, err
}

func (s *ListingService) GetAll() ([]dto.Listing, error) {
	return s.repo.FindAll()
}

func (s *ListingService) GetByID(id int64) (*dto.Listing, error) {
	return s.repo.FindByID(id)
}

func (s *ListingService) Update(id int64, userID int64, input request.CreateListingRequest) (*dto.Listing, error) {
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

func (s *ListingService) Delete(id int64, userID int64) error {
	listing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if listing.OwnerID != userID {
		return errors.New("unauthorized")
	}
	return s.repo.Delete(listing)
}
