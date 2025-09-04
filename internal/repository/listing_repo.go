package repository

import (
	"gorm.io/gorm"
	"homemie/internal/domain"
	"homemie/models/dto"
)

type listingRepo struct {
	db *gorm.DB
}

func NewListingRepo(db *gorm.DB) domain.ListingRepository {
	return &listingRepo{db}
}

func (r *listingRepo) Create(listing *dto.Listing) error {
	return r.db.Create(listing).Error
}

func (r *listingRepo) FindAll() ([]dto.Listing, error) {
	var listings []dto.Listing
	err := r.db.Preload("Owner").Find(&listings).Error
	return listings, err
}

func (r *listingRepo) FindByID(id int64) (*dto.Listing, error) {
	var listing dto.Listing
	err := r.db.Preload("Owner").First(&listing, id).Error
	if err != nil {
		return nil, err
	}
	return &listing, nil
}

func (r *listingRepo) Update(listing *dto.Listing) error {
	return r.db.Save(listing).Error
}

func (r *listingRepo) Delete(listing *dto.Listing) error {
	return r.db.Delete(listing).Error
}

func (r *listingRepo) DB() *gorm.DB {
	return r.db
}
