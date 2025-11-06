package request

import "github.com/google/uuid"

type Favorite struct {
	ListingID uuid.UUID `json:"listing_id" binding:"required"`
}
