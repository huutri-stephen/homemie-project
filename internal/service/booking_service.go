package service

import (
	"context"
	"errors"
	"homemie/db/models"
	"homemie/db/models/dto"
	"homemie/internal/repo"
	"homemie/pkg/logger"
	"time"

	"github.com/aarondl/null/v8"
	"github.com/gin-gonic/gin"
)

type IBookingService interface {
	CreateBooking(ctx context.Context, renterID int64, input dto.CreateBookingRequest) (*models.Booking, error)
	GetMyBookings(ctx context.Context, userID int64) (models.BookingSlice, error)
	GetOwnerBookings(ctx context.Context, ownerID int64) (models.BookingSlice, error)
	RespondToBooking(ctx context.Context, bookingID int64, ownerID int64, req dto.RespondBookingRequest) (*models.Booking, error)
	CancelBooking(ctx context.Context, bookingID int64, userID int64, userRole string) (*models.Booking, error)
	AutoCompleteBookings(ctx context.Context)
}

type bookingService struct {
	bookingRepo repo.IBookingRepository
	listingRepo repo.IListingRepository
}

func NewBookingService(bookingRepo repo.IBookingRepository, listingRepo repo.IListingRepository) IBookingService {
	return &bookingService{bookingRepo, listingRepo}
}

func (s *bookingService) CreateBooking(ctx context.Context, renterID int64, input dto.CreateBookingRequest) (booking *models.Booking, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Create booking",
			"function", "CreateBooking",
			"params", gin.H{"renter_id": renterID, "input": input},
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())

	scheduledTime, err := time.Parse(time.RFC3339, input.ScheduledTime)
	if err != nil {
		logger.Errorw(ctx, "Invalid scheduled_time format", "error", err)
		return nil, errors.New("invalid scheduled_time format")
	}

	booking = &models.Booking{
		ListingID:         input.ListingID,
		RenterID:          renterID,
		ScheduledTime:     scheduledTime,
		MessageFromRenter: null.StringFrom(input.MessageFromRenter),
		Status:            null.StringFrom(models.BookingStatusEnumPending),
	}

	err = s.bookingRepo.Create(ctx, booking)
	return
}

func (s *bookingService) GetMyBookings(ctx context.Context, userID int64) (bookings models.BookingSlice, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Get my bookings",
			"function", "GetMyBookings",
			"params", userID,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	return s.bookingRepo.FindByUserID(ctx, userID)
}

func (s *bookingService) GetOwnerBookings(ctx context.Context, ownerID int64) (bookings models.BookingSlice, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Get owner bookings",
			"function", "GetOwnerBookings",
			"params", ownerID,
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())
	return s.bookingRepo.FindByOwnerID(ctx, ownerID)
}

func (s *bookingService) RespondToBooking(ctx context.Context, bookingID int64, ownerID int64, req dto.RespondBookingRequest) (booking *models.Booking, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Respond to booking",
			"function", "RespondToBooking",
			"params", gin.H{"booking_id": bookingID, "owner_id": ownerID, "req": req},
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())

	booking, err = s.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		logger.Errorw(ctx, "Booking not found", "error", err)
		return nil, errors.New("booking not found")
	}

	listing, err := s.listingRepo.FindByID(ctx, booking.ListingID)
	if err != nil {
		logger.Errorw(ctx, "Listing not found", "error", err)
		return nil, errors.New("listing not found")
	}

	if listing.OwnerID.Int64 != ownerID {
		logger.Warnw(ctx, "Unauthorized booking response attempt")
		return nil, errors.New("unauthorized")
	}

	if booking.Status.String != models.BookingStatusEnumPending {
		logger.Warnw(ctx, "Booking cannot be responded to", "status", booking.Status.String)
		return nil, errors.New("booking cannot be responded to")
	}

	now := time.Now()
	booking.Status = null.StringFrom(req.Status)
	booking.ResponseMessageFromOwner = null.StringFrom(req.ResponseMessage)
	booking.RespondedAt = null.TimeFrom(now)
	booking.RespondedBy = null.Int64From(ownerID)

	err = s.bookingRepo.Update(ctx, booking)
	return
}

func (s *bookingService) CancelBooking(ctx context.Context, bookingID int64, userID int64, userRole string) (booking *models.Booking, err error) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Cancel booking",
			"function", "CancelBooking",
			"params", gin.H{"booking_id": bookingID, "user_id": userID, "user_role": userRole},
			"duration", time.Since(start),
			"error", err,
		)
	}(time.Now())

	booking, err = s.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		logger.Errorw(ctx, "Booking not found", "error", err)
		return nil, errors.New("booking not found")
	}

	listing, err := s.listingRepo.FindByID(ctx, booking.ListingID)
	if err != nil {
		logger.Errorw(ctx, "Listing not found", "error", err)
		return nil, errors.New("listing not found")
	}

	isOwner := listing.OwnerID.Int64 == userID
	isRenter := booking.RenterID == userID

	if !isOwner && !isRenter {
		logger.Warnw(ctx, "Unauthorized booking cancellation attempt")
		return nil, errors.New("unauthorized")
	}

	if booking.Status.String != models.BookingStatusEnumPending && booking.Status.String != models.BookingStatusEnumAccepted {
		logger.Warnw(ctx, "Booking cannot be cancelled", "status", booking.Status.String)
		return nil, errors.New("booking cannot be cancelled")
	}

	now := time.Now()
	booking.Status = null.StringFrom(models.BookingStatusEnumCancelled)
	booking.RespondedAt = null.TimeFrom(now)
	booking.RespondedBy = null.Int64From(userID)

	err = s.bookingRepo.Update(ctx, booking)
	return
}

func (s *bookingService) AutoCompleteBookings(ctx context.Context) {
	defer func(start time.Time) {
		logger.FromContext(ctx).Infow("Auto-complete bookings cron job",
			"function", "AutoCompleteBookings",
			"duration", time.Since(start),
		)
	}(time.Now())

	bookings, err := s.bookingRepo.FindCompletableBookings(ctx)
	if err != nil {
		logger.Errorw(ctx, "Failed to find completable bookings", "error", err)
		return
	}

	logger.Infow(ctx, "Found bookings to auto-complete", "count", len(bookings))

	for _, booking := range bookings {
		booking.Status = null.StringFrom(models.BookingStatusEnumCompleted)
		if err := s.bookingRepo.Update(ctx, booking); err != nil {
			logger.Errorw(ctx, "Failed to auto-complete booking", "booking_id", booking.ID, "error", err)
		}
	}
}