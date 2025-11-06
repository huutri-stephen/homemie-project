package dto

// Request DTOs
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

type CreateListingRequest struct {
	OwnerID         int64
	Title           string                `json:"title" binding:"required"`
	Description     string                `json:"description"`
	PropertyType    string                `json:"property_type"`
	IsShared        bool                  `json:"is_shared"`
	Price           float64               `json:"price" binding:"required"`
	AreaM2          float64               `json:"area_m2"`
	ContactPhone    string                `json:"contact_phone"`
	ContactEmail    string                `json:"contact_email"`
	ContactName     string                `json:"contact_name"`
	NumBedrooms     int32                 `json:"num_bedrooms"`
	NumBathrooms    int32                 `json:"num_bathrooms"`
	NumFloors       int32                 `json:"num_floors"`
	HasBalcony      bool                  `json:"has_balcony"`
	HasParking      bool                  `json:"has_parking"`
	Amenities       []string              `json:"amenities"`
	PetAllowed      bool                  `json:"pet_allowed"`
	AllowedPetTypes []string              `json:"allowed_pet_types"`
	// Latitude        float64               `json:"latitude"`
	// Longitude       float64               `json:"longitude"`
	ListingType     string                `json:"listing_type"`
	DepositAmount   float64               `json:"deposit_amount"`
	Address         AddressRequest        `json:"address"`
	Images          []ListingImageRequest `json:"images"`
}


// Response DTOs
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}