package service

import (
	"homemie/db/models"
	"homemie/db/models/dto"
	"homemie/internal/repo"

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
			IsMain:    image.IsMain,
			SortOrder: image.SortOrder,
		})
	}
	return s.repo.AddListingImages(listingImages)
}
