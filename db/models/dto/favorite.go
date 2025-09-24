package dto

type AddFavoriteRequest struct {
	ListingID int64 `json:"listing_id" binding:"required"`
}