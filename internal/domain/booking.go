package domain

import (
    "homemie/models"
)

type BookingRepository interface {
	Create(*models.Booking) error
	FindByUserID(userID uint) ([]models.Booking, error)
	FindByOwnerID(ownerID uint) ([]models.Booking, error)
}
