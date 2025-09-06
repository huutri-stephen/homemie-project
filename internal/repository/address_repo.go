package repository

import (
	"homemie/models/dto"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AddressRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewAddressRepository(db *gorm.DB, logger *zap.Logger) *AddressRepository {
	return &AddressRepository{db, logger}
}

func (r *AddressRepository) Create(address *dto.Address) (createdAddress *dto.Address, err error) {
	defer func(start time.Time) {
		r.logger.Info("Create address",
			zap.String("function", "Create"),
			zap.Any("params", address),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	if err = r.db.Create(address).Error; err != nil {
		return nil, err
	}
	return address, nil
}