package service

import (
	"homemie/internal/domain"
	"homemie/models/dto"
	"homemie/models/request"

	"go.uber.org/zap"
)

type ListingImageService interface {
	AddListingImages(req request.AddListingImagesRequest) ([]dto.ListingImage, error)
}

type listingImageService struct {
	repo   domain.ListingImageRepository
	logger *zap.Logger
}

func NewListingImageService(repo domain.ListingImageRepository, logger *zap.Logger) ListingImageService {
	return &listingImageService{repo, logger}
}

func (s *listingImageService) AddListingImages(req request.AddListingImagesRequest) ([]dto.ListingImage, error) {
	var listingImages []dto.ListingImage
	for _, image := range req.Images {
		listingImages = append(listingImages, dto.ListingImage{
			ListingID: req.ListingID,
			ImageURL:  image.ImageURL,
			IsMain: image.IsMain,
			SortOrder: image.SortOrder,
		})
	}
	return s.repo.AddListingImages(listingImages)
}
