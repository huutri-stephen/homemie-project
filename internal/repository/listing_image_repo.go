package repository

import (
	"homemie/models/dto"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ListingImageRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewListingImageRepository(db *gorm.DB, logger *zap.Logger) *ListingImageRepository {
	return &ListingImageRepository{db, logger}
}

func (r *ListingImageRepository) Create(listingImage *dto.ListingImage) (createdImage *dto.ListingImage, err error) {
	defer func(start time.Time) {
		r.logger.Info("Create listing image",
			zap.String("function", "Create"),
			zap.Any("params", listingImage),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	if err = r.db.Create(listingImage).Error; err != nil {
		return nil, err
	}
	return listingImage, nil
}
