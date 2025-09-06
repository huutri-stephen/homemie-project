package repository

import (
	"homemie/internal/domain"
	"homemie/models/dto"

	"gorm.io/gorm"
)

type bookingRepo struct {
	db *gorm.DB
}

func NewBookingRepo(db *gorm.DB) domain.BookingRepository {
	return &bookingRepo{db}
}

func (r *bookingRepo) Create(booking *dto.Booking) error {
	return r.db.Create(booking).Error
}

func (r *bookingRepo) FindByUserID(userID int64) ([]dto.Booking, error) {
	var bookings []dto.Booking
	err := r.db.Preload("Listing").Where("renter_id = ?", userID).Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepo) FindByOwnerID(ownerID int64) ([]dto.Booking, error) {
	var bookings []dto.Booking
	err := r.db.Joins("JOIN listings ON listings.id = bookings.listing_id").
		Where("listings.owner_id = ?", ownerID).
		Preload("User").
		Preload("Listing").
		Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepo) FindByID(id int64) (*dto.Booking, error) {
	var booking dto.Booking
	err := r.db.First(&booking, id).Error
	return &booking, err
}

func (r *bookingRepo) Update(booking *dto.Booking) error {
	return r.db.Save(booking).Error
}