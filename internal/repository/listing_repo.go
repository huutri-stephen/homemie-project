package repository

import (
	"gorm.io/gorm"
	"homemie/internal/domain"
	"homemie/models"
)

type listingRepo struct {
	db *gorm.DB
}

func NewListingRepo(db *gorm.DB) domain.ListingRepository {
	return &listingRepo{db}
}

func (r *listingRepo) Create(listing *models.Listing) error {
	return r.db.Create(listing).Error
}

func (r *listingRepo) FindAll() ([]models.Listing, error) {
	var listings []models.Listing
	err := r.db.Preload("Owner").Find(&listings).Error
	return listings, err
}

func (r *listingRepo) FindByID(id uint) (*models.Listing, error) {
	var listing models.Listing
	err := r.db.Preload("Owner").First(&listing, id).Error
	if err != nil {
		return nil, err
	}
	return &listing, nil
}

func (r *listingRepo) Update(listing *models.Listing) error {
	return r.db.Save(listing).Error
}

func (r *listingRepo) Delete(listing *models.Listing) error {
	return r.db.Delete(listing).Error
}
