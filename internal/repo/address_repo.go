package repo

import (
	"homemie/db/models"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IAddressRepository interface {
	Create(address *models.Address) (*models.Address, error)
}

type addressRepo struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewAddressRepository(db *gorm.DB, logger *zap.Logger) IAddressRepository {
	return &addressRepo{db, logger}
}

func (r *addressRepo) Create(address *models.Address) (createdAddress *models.Address, err error) {
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
