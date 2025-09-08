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
	SearchAndFilter(filter *dto.SearchFilterListing) ([]dto.Listing, *dto.Pagination, error)
	DB() *gorm.DB
}
