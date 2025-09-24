package service

import (
	"homemie/db/models"
	"homemie/db/models/dto"
	"homemie/internal/repo"

	"github.com/aarondl/null/v8"
	"go.uber.org/zap"
)

type IListingImageService interface {
	AddListingImages(req dto.AddListingImagesRequest) ([]models.ListingImage, error)
}

type listingImageService struct {
	repo   repo.IListingImageRepository
	logger *zap.Logger
}

func NewListingImageService(repo repo.IListingImageRepository, logger *zap.Logger) IListingImageService {
	return &listingImageService{repo, logger}
}

func (s *listingImageService) AddListingImages(req dto.AddListingImagesRequest) ([]models.ListingImage, error) {
	var listingImages []models.ListingImage
	for _, image := range req.Images {
		listingImages = append(listingImages, models.ListingImage{
			ListingID: req.ListingID,
			ImageURL:  image.ImageURL,
			IsMain:    null.BoolFrom(image.IsMain),
			SortOrder: null.IntFrom(int(image.SortOrder)),
		})
	}
	return s.repo.AddListingImages(listingImages)
}
