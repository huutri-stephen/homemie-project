package service

import (
	"context"
	"homemie/db/models"
	"homemie/db/models/dto"
	"homemie/internal/repo"

	"github.com/aarondl/null/v8"
)

type IListingImageService interface {
	AddListingImages(ctx context.Context, req dto.AddListingImagesRequest) ([]models.ListingImage, error)
}

type listingImageService struct {
	repo repo.IListingImageRepository
}

func NewListingImageService(repo repo.IListingImageRepository) IListingImageService {
	return &listingImageService{repo}
}

func (s *listingImageService) AddListingImages(ctx context.Context, req dto.AddListingImagesRequest) ([]models.ListingImage, error) {
	var listingImages []models.ListingImage
	for _, image := range req.Images {
		listingImages = append(listingImages, models.ListingImage{
			ListingID: req.ListingID,
			ImageURL:  image.ImageURL,
			IsMain:    null.BoolFrom(image.IsMain),
			SortOrder: null.IntFrom(int(image.SortOrder)),
		})
	}
	return s.repo.AddListingImages(ctx, listingImages)
}