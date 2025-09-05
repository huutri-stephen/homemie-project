package dto

import (
	"time"
)

type Listing struct {
	ID              int64 `gorm:"primaryKey"`
	OwnerID         int64 `gorm:"not null"`
	Title           string `gorm:"type:varchar(255);not null"`
	Description     string `gorm:"type:text"`
	PropertyType    string `gorm:"type:property_type_enum;default:'rented_room'"`
	IsShared        bool   `gorm:"default:false"`
	Price           float64 `gorm:"type:decimal(15,2);not null"`
	AreaM2          float64 `gorm:"type:decimal(10,2)"`
	AddressID       int64 `gorm:"not null"`
	ContactPhone    string `gorm:"type:varchar(20)"`
	ContactEmail    string `gorm:"type:varchar(255)"`
	ContactName     string `gorm:"type:varchar(100)"`
	NumBedrooms     int32
	NumBathrooms    int32
	NumFloors       int32
	HasBalcony      bool `gorm:"default:false"`
	HasParking      bool `gorm:"default:false"`
	Amenities       string `gorm:"type:json"`
	PetAllowed      bool   `gorm:"default:false"`
	AllowedPetTypes string `gorm:"type:json"`
	Latitude        float64 `gorm:"type:decimal(10,8)"`
	Longitude       float64 `gorm:"type:decimal(10,8)"`
	ListingType     string `gorm:"type:listing_type_enum;default:'for_rent'"`
	DepositAmount   float64 `gorm:"type:decimal(15,2)"`
	Status          string `gorm:"type:listing_status_enum;default:'pending'"`
	IsFeatured      bool `gorm:"default:false"`
	ViewCount       int32 `gorm:"default:0"`
	PublishedAt     *time.Time
	ExpiresAt       *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DistanceKM      *float64 `gorm:"-" json:"distance_km,omitempty"`
}

type SearchFilterListing struct {
	Keyword         string   `form:"keyword"`
	PropertyType    []string `form:"property_type"`
	IsShared        *bool    `form:"is_shared"`
	MinPrice        *float64 `form:"min_price"`
	MaxPrice        *float64 `form:"max_price"`
	MinArea         *float64 `form:"min_area"`
	MaxArea         *float64 `form:"max_area"`
	NumBedrooms     *int     `form:"num_bedrooms"`
	NumBathrooms    *int     `form:"num_bathrooms"`
	NumFloors       *int     `form:"num_floors"`
	HasBalcony      *bool    `form:"has_balcony"`
	HasParking      *bool    `form:"has_parking"`
	Amenities       []string `form:"amenities"`
	PetAllowed      *bool    `form:"pet_allowed"`
	AllowedPetTypes []string `form:"allowed_pet_types"`
	ListingType     string   `form:"listing_type"`
	Page            int      `form:"page"`
	Limit           int      `form:"limit"`
	CityID          *int     `form:"city_id"`
	WardIDs         []int    `form:"ward_ids"`
	AreaIDs         []int    `form:"area_ids"`
	RadiusKM        *float64 `form:"radius_km"`
	Lat             *float64 `form:"lat"`
	Lon             *float64 `form:"lon"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}