package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"homemie/db/models"
	"homemie/db/models/dto"
	"homemie/pkg/logger"
	"math"
	"time"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/ericlagergren/decimal"
	"github.com/lib/pq"
)

type IListingRepository interface {
	Create(ctx context.Context, listing *models.Listing) error
	FindAll(ctx context.Context) (models.ListingSlice, error)
	FindByID(ctx context.Context, id int64) (*models.Listing, error)
	Update(ctx context.Context, listing *models.Listing) error
	Delete(ctx context.Context, listing *models.Listing) error
	SearchAndFilter(ctx context.Context, filter *dto.SearchFilterListing) (models.ListingSlice, *dto.Pagination, error)
}

type listingRepo struct {
	db *sql.DB
}

func NewListingRepo(db *sql.DB) IListingRepository {
	return &listingRepo{db}
}

func (r *listingRepo) Create(ctx context.Context, listing *models.Listing) (err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Create listing",
			"function", "Create",
			"params", listing,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	err = listing.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return fmt.Errorf("failed to create listing: %w", err)
	}
	return nil
}

func (r *listingRepo) FindAll(ctx context.Context) (listings models.ListingSlice, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Find all listings",
			"function", "FindAll",
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	listings, err = models.Listings().All(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to find all listings: %w", err)
	}
	return
}

func (r *listingRepo) FindByID(ctx context.Context, id int64) (listing *models.Listing, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Find listing by ID",
			"function", "FindByID",
			"params", id,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	listing, err = models.Listings(models.ListingWhere.ID.EQ(id)).One(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to find listing by id: %w", err)
	}
	return
}

func (r *listingRepo) Update(ctx context.Context, listing *models.Listing) (err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Update listing",
			"function", "Update",
			"params", listing,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	_, err = listing.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return fmt.Errorf("failed to update listing: %w", err)
	}
	return nil
}

func (r *listingRepo) Delete(ctx context.Context, listing *models.Listing) (err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Delete listing",
			"function", "Delete",
			"params", listing,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	_, err = listing.Delete(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to delete listing: %w", err)
	}
	return nil
}

func (r *listingRepo) SearchAndFilter(ctx context.Context, filter *dto.SearchFilterListing) (listings models.ListingSlice, pagination *dto.Pagination, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Search and filter listings",
			"function", "SearchAndFilter",
			"params", filter,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())

	queryMods := []qm.QueryMod{
		models.ListingWhere.Status.EQ(null.StringFrom("active")),
	}

	if filter.Keyword != "" {
		queryMods = append(queryMods, qm.Where("search_vector @@ plainto_tsquery('simple', ?)", filter.Keyword))
	}

	queryMods = append(queryMods, qm.InnerJoin("addresses a ON listings.address_id = a.id"))
	if filter.CityID != nil {
		queryMods = append(queryMods, qm.Where("a.city_id = ?", *filter.CityID))
	}
	if len(filter.WardIDs) > 0 {
		wardIDs := make([]interface{}, len(filter.WardIDs))
		for i, v := range filter.WardIDs {
			wardIDs[i] = v
		}
		queryMods = append(queryMods, qm.WhereIn("a.ward_id IN ?", wardIDs...))
	}
	if len(filter.AreaIDs) > 0 {
		areaIDs := make([]interface{}, len(filter.AreaIDs))
		for i, v := range filter.AreaIDs {
			areaIDs[i] = v
		}
		queryMods = append(queryMods, qm.WhereIn("a.area_id IN ?", areaIDs...))
	}

	if filter.RadiusKM != nil && filter.Lat != nil && filter.Lon != nil {
		distanceQuery := `6371 * acos(cos(radians(?)) * cos(radians(a.latitude)) * cos(radians(a.longitude) - radians(?)) + sin(radians(?)) * sin(radians(a.latitude)))`
		queryMods = append(queryMods, qm.Where(distanceQuery+" <= ?", *filter.Lat, *filter.Lon, *filter.Lat, *filter.RadiusKM))

		distanceExpression := fmt.Sprintf("6371 * acos(cos(radians(%f)) * cos(radians(a.latitude)) * cos(radians(a.longitude) - radians(%f)) + sin(radians(%f)) * sin(radians(a.latitude)))", *filter.Lat, *filter.Lon, *filter.Lat)
		queryMods = append(queryMods, qm.Select(models.ListingColumns.ID, models.ListingColumns.OwnerID, models.ListingColumns.Title, models.ListingColumns.Description, models.ListingColumns.PropertyType, models.ListingColumns.IsShared, models.ListingColumns.Price, models.ListingColumns.AreaM2, models.ListingColumns.AddressID, models.ListingColumns.ContactPhone, models.ListingColumns.ContactEmail, models.ListingColumns.ContactName, models.ListingColumns.NumBedrooms, models.ListingColumns.NumBathrooms, models.ListingColumns.NumFloors, models.ListingColumns.HasBalcony, models.ListingColumns.HasParking, models.ListingColumns.Amenities, models.ListingColumns.PetAllowed, models.ListingColumns.AllowedPetTypes, models.ListingColumns.ListingType, models.ListingColumns.DepositAmount, models.ListingColumns.Status, models.ListingColumns.IsFeatured, models.ListingColumns.ViewCount, models.ListingColumns.CreatedAt, models.ListingColumns.UpdatedAt, models.ListingColumns.PublishedAt, models.ListingColumns.ExpiresAt, models.ListingColumns.SearchVector, distanceExpression+" as distance_km"))
		queryMods = append(queryMods, qm.OrderBy("distance_km ASC"))
	} else {
		queryMods = append(queryMods, qm.OrderBy("created_at DESC"))
	}

	if len(filter.PropertyType) > 0 {
		queryMods = append(queryMods, qm.Where("property_type = ANY(?)", pq.Array(filter.PropertyType)))
	}
	if filter.IsShared != nil {
		queryMods = append(queryMods, models.ListingWhere.IsShared.EQ(null.BoolFromPtr(filter.IsShared)))
	}
	if filter.MinPrice != nil {
		minPriceDecimal := new(decimal.Big)
		minPriceDecimal.SetFloat64(*filter.MinPrice)
		queryMods = append(queryMods, models.ListingWhere.Price.GTE(types.NewDecimal(minPriceDecimal)))
	}
	if filter.MaxPrice != nil {
		maxPriceDecimal := new(decimal.Big)
		maxPriceDecimal.SetFloat64(*filter.MaxPrice)
		queryMods = append(queryMods, models.ListingWhere.Price.LTE(types.NewDecimal(maxPriceDecimal)))
	}
	if filter.MinArea != nil {
		minAreaDecimal := new(decimal.Big)
		minAreaDecimal.SetFloat64(*filter.MinArea)
		queryMods = append(queryMods, models.ListingWhere.AreaM2.GTE(types.NewNullDecimal(minAreaDecimal)))
	}
	if filter.MaxArea != nil {
		maxAreaDecimal := new(decimal.Big)
		maxAreaDecimal.SetFloat64(*filter.MaxArea)
		queryMods = append(queryMods, models.ListingWhere.AreaM2.LTE(types.NewNullDecimal(maxAreaDecimal)))
	}
	if filter.NumBedrooms != nil {
		queryMods = append(queryMods, models.ListingWhere.NumBedrooms.EQ(null.IntFromPtr(filter.NumBedrooms)))
	}
	if filter.NumBathrooms != nil {
		queryMods = append(queryMods, models.ListingWhere.NumBathrooms.EQ(null.IntFromPtr(filter.NumBathrooms)))
	}
	if filter.NumFloors != nil {
		queryMods = append(queryMods, models.ListingWhere.NumFloors.EQ(null.IntFromPtr(filter.NumFloors)))
	}
	if filter.HasBalcony != nil {
		queryMods = append(queryMods, models.ListingWhere.HasBalcony.EQ(null.BoolFromPtr(filter.HasBalcony)))
	}
	if filter.HasParking != nil {
		queryMods = append(queryMods, models.ListingWhere.HasParking.EQ(null.BoolFromPtr(filter.HasParking)))
	}
	if len(filter.Amenities) > 0 {
		amenitiesJSON, _ := json.Marshal(filter.Amenities)
		queryMods = append(queryMods, qm.Where("amenities @> ?", string(amenitiesJSON)))
	}
	if len(filter.AllowedPetTypes) > 0 {
		petTypesJSON, _ := json.Marshal(filter.AllowedPetTypes)
		queryMods = append(queryMods, qm.Where("allowed_pet_types @> ?", string(petTypesJSON)))
	}
	if filter.PetAllowed != nil {
		queryMods = append(queryMods, models.ListingWhere.PetAllowed.EQ(null.BoolFromPtr(filter.PetAllowed)))
	}
	if filter.ListingType != "" {
		queryMods = append(queryMods, models.ListingWhere.ListingType.EQ(null.StringFrom(string(filter.ListingType))))
	}
	total, err := models.Listings(queryMods...).Count(ctx, r.db)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count listings: %w", err)
	}

	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.Limit == 0 {
		filter.Limit = 20
	}
	offset := (filter.Page - 1) * filter.Limit
	queryMods = append(queryMods, qm.Limit(filter.Limit), qm.Offset(offset))

	listings, err = models.Listings(queryMods...).All(ctx, r.db)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to search listings: %w", err)
	}

	pagination = &dto.Pagination{
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(filter.Limit))),
	}

	return
}