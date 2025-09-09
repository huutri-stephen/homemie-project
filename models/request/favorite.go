package request

// AddFavoriteRequest defines the request body for adding a favorite.
type AddFavoriteRequest struct {
	ListingID int64 `json:"listing_id" binding:"required"`
}
