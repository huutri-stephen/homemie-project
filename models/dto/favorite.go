package dto

import "time"

type Favorite struct {
	UserID    int64     `json:"user_id"`
	ListingID int64     `json:"listing_id"`
	CreatedAt time.Time `json:"created_at"`
}