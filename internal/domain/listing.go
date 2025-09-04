package domain

import (
	"homemie/models/dto"

	"gorm.io/gorm"
)

type ListingRepository interface {
	Create(listing *dto.Listing) error
	FindAll() ([]dto.Listing, error)
	FindByID(id int64) (*dto.Listing, error)
	Update(listing *dto.Listing) error
	Delete(listing *dto.Listing) error
	DB() *gorm.DB
}
