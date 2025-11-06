package domain

import (
	"time"

	"github.com/google/uuid"
)

type Favorite struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ListingID uuid.UUID `json:"listing_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
