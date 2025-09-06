package repository

import (
	"homemie/internal/domain"
	"homemie/models/dto"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type bookingRepo struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewBookingRepo(db *gorm.DB, logger *zap.Logger) domain.BookingRepository {
	return &bookingRepo{db, logger}
}

func (r *bookingRepo) Create(booking *dto.Booking) (err error) {
	defer func(start time.Time) {
		r.logger.Info("Create booking",
			zap.String("function", "Create"),
			zap.Any("params", booking),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	return r.db.Create(booking).Error
}

func (r *bookingRepo) FindByUserID(userID int64) (bookings []dto.Booking, err error) {
	defer func(start time.Time) {
		r.logger.Info("Find bookings by user ID",
			zap.String("function", "FindByUserID"),
			zap.Int64("params", userID),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	err = r.db.Preload("Listing").Where("renter_id = ?", userID).Find(&bookings).Error
	return
}

func (r *bookingRepo) FindByOwnerID(ownerID int64) (bookings []dto.Booking, err error) {
	defer func(start time.Time) {
		r.logger.Info("Find bookings by owner ID",
			zap.String("function", "FindByOwnerID"),
			zap.Int64("params", ownerID),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	err = r.db.Joins("JOIN listings ON listings.id = bookings.listing_id").
		Where("listings.owner_id = ?", ownerID).
		Preload("User").
		Preload("Listing").
		Find(&bookings).Error
	return
}

func (r *bookingRepo) FindByID(id int64) (booking *dto.Booking, err error) {
	defer func(start time.Time) {
		r.logger.Info("Find booking by ID",
			zap.String("function", "FindByID"),
			zap.Int64("params", id),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	booking = &dto.Booking{}
	err = r.db.Preload("Listing").First(booking, id).Error
	return
}

func (r *bookingRepo) Update(booking *dto.Booking) (err error) {
	defer func(start time.Time) {
		r.logger.Info("Update booking",
			zap.String("function", "Update"),
			zap.Any("params", booking),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	return r.db.Save(booking).Error
}

func (r *bookingRepo) FindCompletableBookings() (bookings []dto.Booking, err error) {
	defer func(start time.Time) {
		r.logger.Info("Find completable bookings",
			zap.String("function", "FindCompletableBookings"),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	err = r.db.Where("status = ? AND scheduled_time < NOW() - INTERVAL '1 day'", dto.BookingStatusAccepted).Find(&bookings).Error
	return
}
