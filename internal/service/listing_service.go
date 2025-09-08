package service

import (
	"errors"
	"homemie/internal/domain"
	"homemie/internal/repository"
	"homemie/models/dto"
	"homemie/models/request"
	"homemie/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ListingService struct {
	repo             domain.ListingRepository
	addressRepo      *repository.AddressRepository
	listingImageRepo *repository.ListingImageRepository
	logger           *zap.Logger
}

func NewListingService(repo domain.ListingRepository, addressRepo *repository.AddressRepository, listingImageRepo *repository.ListingImageRepository, logger *zap.Logger) *ListingService {
	return &ListingService{repo, addressRepo, listingImageRepo, logger}
}

func (s *ListingService) Create(input request.CreateListingRequest) (listing *dto.Listing, err error) {
	defer func(start time.Time) {
		s.logger.Info("Create listing",
			zap.String("function", "Create"),
			zap.Any("params", input),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

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
		s.logger.Error("Failed to create address", zap.Error(err))
		return nil, err
	}

	listing = &dto.Listing{
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
		ListingType:   input.ListingType,
		DepositAmount: input.DepositAmount,
	}

	err = s.repo.Create(listing)
	if err != nil {
		s.logger.Error("Failed to create listing", zap.Error(err))
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
			s.logger.Error("Failed to create listing image", zap.Error(err))
			return nil, err
		}
	}

	return listing, err
}

func (s *ListingService) GetAll() (listings []dto.Listing, err error) {
	defer func(start time.Time) {
		s.logger.Info("Get all listings",
			zap.String("function", "GetAll"),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	return s.repo.FindAll()
}

func (s *ListingService) SearchAndFilter(filter *dto.SearchFilterListing) (listings []dto.Listing, pagination *dto.Pagination, err error) {
	defer func(start time.Time) {
		s.logger.Info("Search and filter listings",
			zap.String("function", "SearchAndFilter"),
			zap.Any("params", filter),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	return s.repo.SearchAndFilter(filter)
}

func (s *ListingService) GetByID(id int64) (listing *dto.Listing, err error) {
	defer func(start time.Time) {
		s.logger.Info("Get listing by ID",
			zap.String("function", "GetByID"),
			zap.Int64("params", id),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	return s.repo.FindByID(id)
}

func (s *ListingService) Update(id int64, userID int64, input request.CreateListingRequest) (listing *dto.Listing, err error) {
	defer func(start time.Time) {
		s.logger.Info("Update listing",
			zap.String("function", "Update"),
			zap.Any("params", gin.H{"id": id, "user_id": userID, "input": input}),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	listing, err = s.repo.FindByID(id)
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

func (s *ListingService) Delete(id int64, userID int64) (err error) {
	defer func(start time.Time) {
		s.logger.Info("Delete listing",
			zap.String("function", "Delete"),
			zap.Any("params", gin.H{"id": id, "user_id": userID}),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	listing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if listing.OwnerID != userID {
		return errors.New("unauthorized")
	}
	return s.repo.Delete(listing)
}
