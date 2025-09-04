package repository

import (
	"homemie/models/dto"

	"gorm.io/gorm"
)

type AddressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) *AddressRepository {
	return &AddressRepository{db}
}

func (r *AddressRepository) Create(address *dto.Address) (*dto.Address, error) {
	if err := r.db.Create(address).Error; err != nil {
		return nil, err
	}
	return address, nil
}
