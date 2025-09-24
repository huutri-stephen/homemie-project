package service

import (
	"errors"
	"homemie/db/models"
	"homemie/db/models/dto"
	"homemie/internal/repo"

	"homemie/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IListingService interface {
	Create(input dto.CreateListingRequest) (*models.Listing, error)
	GetAll() ([]models.Listing, error)
	SearchAndFilter(filter *dto.SearchFilterListing) ([]models.Listing, *dto.Pagination, error)
	GetByID(id int64) (*models.Listing, error)
	Update(id int64, userID int64, input dto.CreateListingRequest) (*models.Listing, error)
	Delete(id int64, userID int64) error
}

type listingService struct {
	listingRepo      repo.IListingRepository
	addressRepo      repo.IAddressRepository
	listingImageRepo repo.IListingImageRepository
	logger           *zap.Logger
}

func NewListingService(listingRepo repo.IListingRepository, addressRepo repo.IAddressRepository, listingImageRepo repo.IListingImageRepository, logger *zap.Logger) IListingService {
	return &listingService{listingRepo: listingRepo, addressRepo: addressRepo, listingImageRepo: listingImageRepo, logger: logger}
}

func (s *listingService) Create(input dto.CreateListingRequest) (listing *models.Listing, err error) {
	defer func(start time.Time) {
		s.logger.Info("Create listing",
			zap.String("function", "Create"),
			zap.Any("params", input),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	// todo: find current addrest if address exist
	addr := models.Address{
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

	listing = &models.Listing{
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

	err = s.listingRepo.Create(listing)
	if err != nil {
		s.logger.Error("Failed to create listing", zap.Error(err))
		return nil, err
	}

	var listingImages []models.ListingImage
	for _, image := range input.Images {
		listingImages = append(listingImages, models.ListingImage{
			ListingID: listing.ID,
			ImageURL:  image.ImageURL,
			IsMain:    image.IsMain,
			SortOrder: image.SortOrder,
		})
	}

	if len(listingImages) > 0 {
		_, err = s.listingImageRepo.AddListingImages(listingImages)
		if err != nil {
			s.logger.Error("Failed to create listing images", zap.Error(err))
			// Here we might want to decide if we should delete the created listing
			return nil, err
		}
	}

	return listing, err
}

func (s *listingService) GetAll() (listings []models.Listing, err error) {
	defer func(start time.Time) {
		s.logger.Info("Get all listings",
			zap.String("function", "GetAll"),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	return s.listingRepo.FindAll()
}

func (s *listingService) SearchAndFilter(filter *dto.SearchFilterListing) (listings []models.Listing, pagination *dto.Pagination, err error) {
	defer func(start time.Time) {
		s.logger.Info("Search and filter listings",
			zap.String("function", "SearchAndFilter"),
			zap.Any("params", filter),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	return s.listingRepo.SearchAndFilter(filter)
}

func (s *listingService) GetByID(id int64) (listing *models.Listing, err error) {
	defer func(start time.Time) {
		s.logger.Info("Get listing by ID",
			zap.String("function", "GetByID"),
			zap.Int64("params", id),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	return s.listingRepo.FindByID(id)
}

func (s *listingService) Update(id int64, userID int64, input dto.CreateListingRequest) (listing *models.Listing, err error) {
	defer func(start time.Time) {
		s.logger.Info("Update listing",
			zap.String("function", "Update"),
			zap.Any("params", gin.H{"id": id, "user_id": userID, "input": input}),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	listing, err = s.listingRepo.FindByID(id)
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

	err = s.listingRepo.Update(listing)
	return listing, err
}

func (s *listingService) Delete(id int64, userID int64) (err error) {
	defer func(start time.Time) {
		s.logger.Info("Delete listing",
			zap.String("function", "Delete"),
			zap.Any("params", gin.H{"id": id, "user_id": userID}),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	listing, err := s.listingRepo.FindByID(id)
	if err != nil {
		return err
	}
	if listing.OwnerID != userID {
		return errors.New("unauthorized")
	}
	return s.listingRepo.Delete(listing)
}
