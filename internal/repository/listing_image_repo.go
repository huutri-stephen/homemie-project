package repository

import (
	"homemie/internal/domain"
	"homemie/models/dto"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type listingImageRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewListingImageRepository(db *gorm.DB, logger *zap.Logger) domain.ListingImageRepository {
	return &listingImageRepository{db, logger}
}

func (r *listingImageRepository) AddListingImages(listingImages []dto.ListingImage) ([]dto.ListingImage, error) {
	if err := r.db.Create(&listingImages).Error; err != nil {
		return nil, err
	}
	return listingImages, nil
}
