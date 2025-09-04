package repository

import (
	"homemie/models/dto"

	"gorm.io/gorm"
)

type ListingImageRepository struct {
	db *gorm.DB
}

func NewListingImageRepository(db *gorm.DB) *ListingImageRepository {
	return &ListingImageRepository{db}
}

func (r *ListingImageRepository) Create(listingImage *dto.ListingImage) (*dto.ListingImage, error) {
	if err := r.db.Create(listingImage).Error; err != nil {
		return nil, err
	}
	return listingImage, nil
}
