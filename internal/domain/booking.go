package domain

import (
    "homemie/models/dto"
)

type BookingRepository interface {
	Create(*dto.Booking) error
	FindByUserID(userID int64) ([]dto.Booking, error)
	FindByOwnerID(ownerID int64) ([]dto.Booking, error)
}
