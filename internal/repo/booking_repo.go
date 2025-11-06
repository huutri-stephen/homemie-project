package repo

import (
	"context"
	"database/sql"
	"fmt"
	"homemie/db/models"
	"homemie/pkg/logger"
	"time"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

type IBookingRepository interface {
	Create(ctx context.Context, booking *models.Booking) error
	FindByUserID(ctx context.Context, userID int64) (models.BookingSlice, error)
	FindByOwnerID(ctx context.Context, ownerID int64) (models.BookingSlice, error)
	FindByID(ctx context.Context, id int64) (*models.Booking, error)
	Update(ctx context.Context, booking *models.Booking) error
	FindCompletableBookings(ctx context.Context) (models.BookingSlice, error)
}

type bookingRepo struct {
	db *sql.DB
}

func NewBookingRepo(db *sql.DB) IBookingRepository {
	return &bookingRepo{db}
}

func (r *bookingRepo) Create(ctx context.Context, booking *models.Booking) (err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Create booking",
			"function", "Create",
			"params", booking,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	err = booking.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}
	return nil
}

func (r *bookingRepo) FindByUserID(ctx context.Context, userID int64) (bookings models.BookingSlice, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Find bookings by user ID",
			"function", "FindByUserID",
			"params", userID,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	bookings, err = models.Bookings(
		models.BookingWhere.RenterID.EQ(userID),
		qm.Load(models.BookingRels.Listing),
	).All(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to find bookings by user id: %w", err)
	}
	return
}

func (r *bookingRepo) FindByOwnerID(ctx context.Context, ownerID int64) (bookings models.BookingSlice, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Find bookings by owner ID",
			"function", "FindByOwnerID",
			"params", ownerID,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	bookings, err = models.Bookings(
		qm.InnerJoin("listings l on l.id = bookings.listing_id"),
		qm.Where("l.owner_id = ?", ownerID),
		qm.Load(models.BookingRels.Renter),
		qm.Load(models.BookingRels.Listing),
	).All(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to find bookings by owner id: %w", err)
	}
	return
}

func (r *bookingRepo) FindByID(ctx context.Context, id int64) (booking *models.Booking, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Find booking by ID",
			"function", "FindByID",
			"params", id,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	booking, err = models.Bookings(
		models.BookingWhere.ID.EQ(id),
		qm.Load(models.BookingRels.Listing),
	).One(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to find booking by id: %w", err)
	}
	return
}

func (r *bookingRepo) Update(ctx context.Context, booking *models.Booking) (err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Update booking",
			"function", "Update",
			"params", booking,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	_, err = booking.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}
	return nil
}

func (r *bookingRepo) FindCompletableBookings(ctx context.Context) (bookings models.BookingSlice, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Find completable bookings",
			"function", "FindCompletableBookings",
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	bookings, err = models.Bookings(
		models.BookingWhere.Status.EQ(null.StringFrom(models.BookingStatusEnumAccepted)),
		qm.Where("scheduled_time < NOW() - INTERVAL '1 day'"),
	).All(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to find completable bookings: %w", err)
	}
	return
}