package repository

import (
	"homemie/internal/domain"
	"homemie/models/dto"
	"homemie/pkg/utils"
	"math"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type listingRepo struct {
	db *gorm.DB
}

func NewListingRepo(db *gorm.DB) domain.ListingRepository {
	return &listingRepo{db}
}

func (r *listingRepo) Create(listing *dto.Listing) error {
	return r.db.Create(listing).Error
}

func (r *listingRepo) FindAll() ([]dto.Listing, error) {
	var listings []dto.Listing
	err := r.db.Find(&listings).Error
	return listings, err
}

func (r *listingRepo) FindByID(id int64) (*dto.Listing, error) {
	var listing dto.Listing
	err := r.db.First(&listing, id).Error
	if err != nil {
		return nil, err
	}
	return &listing, nil
}

func (r *listingRepo) Update(listing *dto.Listing) error {
	return r.db.Save(listing).Error
}

func (r *listingRepo) Delete(listing *dto.Listing) error {
	return r.db.Delete(listing).Error
}

func (r *listingRepo) SearchAndFilter(filter *dto.SearchFilterListing) ([]dto.Listing, *dto.Pagination, error) {
	var listings []dto.Listing
	var total int64

	// base query
	db := r.db.Model(&dto.Listing{}).Where("listings.status = ?", "active")

	// keyword full-text search
	if filter.Keyword != "" {
		db = db.Where("search_vector @@ plainto_tsquery('simple', ?)", filter.Keyword)
	}

	// join with addresses
	db = db.Joins("JOIN addresses a ON listings.address_id = a.id")
	if filter.CityID != nil {
		db = db.Where("a.city_id = ?", *filter.CityID)
	}
	if len(filter.WardIDs) > 0 {
		db = db.Where("a.ward_id IN (?)", filter.WardIDs)
	}
	if len(filter.AreaIDs) > 0 {
		db = db.Where("a.area_id IN (?)", filter.AreaIDs)
	}

	// radius filter
	if filter.RadiusKM != nil && filter.Lat != nil && filter.Lon != nil {
		distanceQuery := `
			6371 * acos(
				cos(radians(?)) * cos(radians(a.latitude)) * cos(radians(a.longitude) - radians(?)) 
				+ sin(radians(?)) * sin(radians(a.latitude))
			) <= ?
		`
		db = db.Where(distanceQuery, *filter.Lat, *filter.Lon, *filter.Lat, *filter.RadiusKM)
		db = db.Select(`listings.*, 
			6371 * acos(
				cos(radians(?)) * cos(radians(a.latitude)) * cos(radians(a.longitude) - radians(?)) 
				+ sin(radians(?)) * sin(radians(a.latitude))
			) as distance_km
		`, *filter.Lat, *filter.Lon, *filter.Lat).
			Order("distance_km ASC")
	} else {
		db = db.Select("listings.*").Order("created_at DESC")
	}

	// property filters
	if len(filter.PropertyType) > 0 {
		db = db.Where("property_type = ANY(?)", pq.Array(filter.PropertyType))
	}
	if filter.IsShared != nil {
		db = db.Where("is_shared = ?", *filter.IsShared)
	}
	if filter.MinPrice != nil {
		db = db.Where("price >= ?", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		db = db.Where("price <= ?", *filter.MaxPrice)
	}
	if filter.MinArea != nil {
		db = db.Where("area_m2 >= ?", *filter.MinArea)
	}
	if filter.MaxArea != nil {
		db = db.Where("area_m2 <= ?", *filter.MaxArea)
	}
	if filter.NumBedrooms != nil {
		db = db.Where("num_bedrooms = ?", *filter.NumBedrooms)
	}
	if filter.NumBathrooms != nil {
		db = db.Where("num_bathrooms = ?", *filter.NumBathrooms)
	}
	if filter.NumFloors != nil {
		db = db.Where("num_floors = ?", *filter.NumFloors)
	}
	if filter.HasBalcony != nil {
		db = db.Where("has_balcony = ?", *filter.HasBalcony)
	}
	if filter.HasParking != nil {
		db = db.Where("has_parking = ?", *filter.HasParking)
	}
	if len(filter.Amenities) > 0 {
    db = db.Where("amenities @> ?", utils.ConvertStringArrayToJSON(filter.Amenities))
	}
	if len(filter.AllowedPetTypes) > 0 {
		db = db.Where("allowed_pet_types @> ?", utils.ConvertStringArrayToJSON(filter.AllowedPetTypes))
	}
	if filter.PetAllowed != nil {
		db = db.Where("pet_allowed = ?", *filter.PetAllowed)
	}
	if filter.ListingType != "" {
		db = db.Where("listing_type = ?", filter.ListingType)
	}

	countDB := db.Session(&gorm.Session{})
	countDB.Select("*").Count(&total)

	// pagination
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.Limit == 0 {
		filter.Limit = 20
	}
	offset := (filter.Page - 1) * filter.Limit

	// main query
	if err := db.Offset(offset).Limit(filter.Limit).Find(&listings).Error; err != nil {
		return nil, nil, err
	}

	// pagination response
	pagination := &dto.Pagination{
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(filter.Limit))),
	}

	return listings, pagination, nil
}

func (r *listingRepo) DB() *gorm.DB {
	return r.db
}