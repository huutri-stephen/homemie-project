package request

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
	Latitude        float64               `json:"latitude"`
	Longitude       float64               `json:"longitude"`
	ListingType     string                `json:"listing_type"`
	DepositAmount   float64               `json:"deposit_amount"`
	Address         AddressRequest        `json:"address"`
	Images          []ListingImageRequest `json:"images"`
}
