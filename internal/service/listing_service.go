package service

import (
	"context"
	"encoding/json"
	"errors"
	"homemie/db/models"
	"homemie/db/models/dto"
	"homemie/internal/repo"
	"homemie/pkg/logger"
	"homemie/pkg/utils"
	"time"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/gin-gonic/gin"
)

type IListingService interface {
	Create(ctx context.Context, input dto.CreateListingRequest) (*models.Listing, error)
	GetAll(ctx context.Context) (models.ListingSlice, error)
	SearchAndFilter(ctx context.Context, filter *dto.SearchFilterListing) (models.ListingSlice, *dto.Pagination, error)
	GetByID(ctx context.Context, id int64) (*models.Listing, error)
	Update(ctx context.Context, id int64, userID int64, input dto.CreateListingRequest) (*models.Listing, error)
	Delete(ctx context.Context, id int64, userID int64) error
}

type listingService struct {
	listingRepo      repo.IListingRepository
	addressRepo      repo.IAddressRepository
	listingImageRepo repo.IListingImageRepository
}

func NewListingService(listingRepo repo.IListingRepository, addressRepo repo.IAddressRepository, listingImageRepo repo.IListingImageRepository) IListingService {
	return &listingService{listingRepo: listingRepo, addressRepo: addressRepo, listingImageRepo: listingImageRepo}
}

func (s *listingService) Create(ctx context.Context, input dto.CreateListingRequest) (listing *models.Listing, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Create listing",
			"function", "Create",
			"params", input,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())

	// todo: find current addrest if address exist

	addrReq := input.Address

	addr := models.Address{
		CityID:       addrReq.CityID,
		WardID:       addrReq.WardID,
		AreaID:       null.Int64FromPtr(addrReq.AreaID),
		Street:       null.StringFrom(addrReq.Street),
		HouseNumber:  null.StringFrom(addrReq.HouseNumber),
		BuildingName: null.StringFrom(addrReq.BuildingName),
		FloorNumber:  null.IntFrom(int(addrReq.FloorNumber)),
		RoomNumber:   null.StringFrom(addrReq.RoomNumber),
		Latitude:     types.NewNullDecimal(utils.FloatToDecimal(addrReq.Latitude)),
		Longitude:    types.NewNullDecimal(utils.FloatToDecimal(addrReq.Longitude)),
	}
	address, err := s.addressRepo.Create(ctx, &addr)
	if err != nil {
		logger.Errorw(ctx, "Failed to create address", "error", err)
		return nil, err
	}

	amenitiesJSON, _ := json.Marshal(input.Amenities)
	petTypesJSON, _ := json.Marshal(input.AllowedPetTypes)

	listing = &models.Listing{
		OwnerID:         null.Int64From(input.OwnerID),
		Title:           input.Title,
		Description:     null.StringFrom(input.Description),
		PropertyType:    null.StringFrom(input.PropertyType),
		IsShared:        null.BoolFrom(input.IsShared),
		Price:           types.NewDecimal(utils.FloatToDecimal(input.Price)),
		AreaM2:          types.NewNullDecimal(utils.FloatToDecimal(input.AreaM2)),
		AddressID:       address.ID,
		ContactPhone:    null.StringFrom(input.ContactPhone),
		ContactEmail:    null.StringFrom(input.ContactEmail),
		ContactName:     null.StringFrom(input.ContactName),
		NumBedrooms:     null.IntFrom(int(input.NumBedrooms)),
		NumBathrooms:    null.IntFrom(int(input.NumBathrooms)),
		NumFloors:       null.IntFrom(int(input.NumFloors)),
		HasBalcony:      null.BoolFrom(input.HasBalcony),
		HasParking:      null.BoolFrom(input.HasParking),
		Amenities:       null.JSONFrom(amenitiesJSON),
		PetAllowed:      null.BoolFrom(input.PetAllowed),
		AllowedPetTypes: null.JSONFrom(petTypesJSON),
		ListingType:     null.StringFrom(input.ListingType),
		DepositAmount:   types.NewNullDecimal(utils.FloatToDecimal(input.DepositAmount)),
	}

	err = s.listingRepo.Create(ctx, listing)
	if err != nil {
		logger.Errorw(ctx, "Failed to create listing", "error", err)
		return nil, err
	}

	var listingImages []models.ListingImage
	for _, image := range input.Images {
		listingImages = append(listingImages, models.ListingImage{
			ListingID: listing.ID,
			ImageURL:  image.ImageURL,
			IsMain:    null.BoolFrom(image.IsMain),
			SortOrder: null.IntFrom(int(image.SortOrder)),
		})
	}

	if len(listingImages) > 0 {
		_, err = s.listingImageRepo.AddListingImages(ctx, listingImages)
		if err != nil {
			logger.Errorw(ctx, "Failed to create listing images", "error", err)
			// Here we might want to decide if we should delete the created listing
			return nil, err
		}
	}

	return listing, err
}

func (s *listingService) GetAll(ctx context.Context) (listings models.ListingSlice, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Get all listings",
			"function", "GetAll",
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	return s.listingRepo.FindAll(ctx)
}

func (s *listingService) SearchAndFilter(ctx context.Context, filter *dto.SearchFilterListing) (listings models.ListingSlice, pagination *dto.Pagination, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Search and filter listings",
			"function", "SearchAndFilter",
			"params", filter,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	return s.listingRepo.SearchAndFilter(ctx, filter)
}

func (s *listingService) GetByID(ctx context.Context, id int64) (listing *models.Listing, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Get listing by ID",
			"function", "GetByID",
			"params", id,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	return s.listingRepo.FindByID(ctx, id)
}

func (s *listingService) Update(ctx context.Context, id int64, userID int64, input dto.CreateListingRequest) (listing *models.Listing, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Update listing",
			"function", "Update",
			"params", gin.H{"id": id, "user_id": userID, "input": input},
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())

	listing, err = s.listingRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if listing.OwnerID.Int64 != userID {
		return nil, errors.New("unauthorized")
	}

	listing.Title = input.Title
	listing.Description = null.StringFrom(input.Description)
	listing.Price = types.NewDecimal(utils.FloatToDecimal(input.Price))
	// listing.Address = input.Address
	// listing.City = input.City

	err = s.listingRepo.Update(ctx, listing)
	return listing, err
}

func (s *listingService) Delete(ctx context.Context, id int64, userID int64) (err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Delete listing",
			"function", "Delete",
			"params", gin.H{"id": id, "user_id": userID},
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())

	listing, err := s.listingRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if listing.OwnerID.Int64 != userID {
		return errors.New("unauthorized")
	}
	return s.listingRepo.Delete(ctx, listing)
}
