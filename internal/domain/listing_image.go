package domain

import (
	"homemie/models/dto"
)

type ListingImageRepository interface {
	AddListingImages(listingImages []dto.ListingImage) ([]dto.ListingImage, error)
}
