package domain

import "mihome/models"

type ListingRepository interface {
	Create(listing *models.Listing) error
	FindAll() ([]models.Listing, error)
	FindByID(id uint) (*models.Listing, error)
	Update(listing *models.Listing) error
	Delete(listing *models.Listing) error
}
