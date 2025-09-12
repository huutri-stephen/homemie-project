package request

type AddListingImagesRequest struct {
	ListingID int64   `json:"listing_id" binding:"required"`
	Images []Image `json:"images" binding:"required"`
}

type ListingImageRequest struct {
	ImageURL  string `json:"image_url" binding:"required"`
	IsMain    bool   `json:"is_main"`
	SortOrder int32    `json:"sort_order"`
}

type Image struct {
	ImageURL string 	`json:"image_url" binding:"required"`
	IsMain    bool   	`json:"is_main"`
	SortOrder int32    	`json:"sort_order"`
}
