package repo

import (
	"homemie/db/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IListingImageRepository interface {
	AddListingImages(listingImages []models.ListingImage) ([]models.ListingImage, error)
}

type listingImageRepo struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewListingImageRepository(db *gorm.DB, logger *zap.Logger) IListingImageRepository {
	return &listingImageRepo{db, logger}
}

func (r *listingImageRepo) AddListingImages(listingImages []models.ListingImage) ([]models.ListingImage, error) {
	if err := r.db.Create(&listingImages).Error; err != nil {
		return nil, err
	}
	return listingImages, nil
}
