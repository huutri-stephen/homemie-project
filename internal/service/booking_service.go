package service

import (
	"errors"
	"homemie/db/models"
	"homemie/db/models/dto"
	"homemie/internal/repo"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IBookingService interface {
	CreateBooking(renterID int64, input dto.CreateBookingRequest) (*models.Booking, error)
	GetMyBookings(userID int64) ([]models.Booking, error)
	GetOwnerBookings(ownerID int64) ([]models.Booking, error)
	RespondToBooking(bookingID int64, ownerID int64, req dto.RespondBookingRequest) (*models.Booking, error)
	CancelBooking(bookingID int64, userID int64, userRole string) (*models.Booking, error)
	AutoCompleteBookings()
}

type bookingService struct {
	bookingRepo repo.IBookingRepository
	listingRepo repo.IListingRepository
	logger      *zap.Logger
}

func NewBookingService(bookingRepo repo.IBookingRepository, listingRepo repo.IListingRepository, logger *zap.Logger) IBookingService {
	return &bookingService{bookingRepo, listingRepo, logger}
}

func (s *bookingService) CreateBooking(renterID int64, input dto.CreateBookingRequest) (booking *models.Booking, err error) {
	defer func(start time.Time) {
		s.logger.Info("Create booking",
			zap.String("function", "CreateBooking"),
			zap.Any("params", gin.H{"renter_id": renterID, "input": input}),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	scheduledTime, err := time.Parse(time.RFC3339, input.ScheduledTime)
	if err != nil {
		s.logger.Error("Invalid scheduled_time format", zap.Error(err))
		return nil, errors.New("invalid scheduled_time format")
	}

	booking = &models.Booking{
		ListingID:         input.ListingID,
		RenterID:          renterID,
		ScheduledTime:     scheduledTime,
		MessageFromRenter: input.MessageFromRenter,
		Status:            models.BookingStatusEnumPending,
	}

	err = s.bookingRepo.Create(booking)
	return
}

func (s *bookingService) GetMyBookings(userID int64) (bookings []models.Booking, err error) {
	defer func(start time.Time) {
		s.logger.Info("Get my bookings",
			zap.String("function", "GetMyBookings"),
			zap.Int64("params", userID),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	return s.bookingRepo.FindByUserID(userID)
}

func (s *bookingService) GetOwnerBookings(ownerID int64) (bookings []models.Booking, err error) {
	defer func(start time.Time) {
		s.logger.Info("Get owner bookings",
			zap.String("function", "GetOwnerBookings"),
			zap.Int64("params", ownerID),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())
	return s.bookingRepo.FindByOwnerID(ownerID)
}

func (s *bookingService) RespondToBooking(bookingID int64, ownerID int64, req dto.RespondBookingRequest) (booking *models.Booking, err error) {
	defer func(start time.Time) {
		s.logger.Info("Respond to booking",
			zap.String("function", "RespondToBooking"),
			zap.Any("params", gin.H{"booking_id": bookingID, "owner_id": ownerID, "req": req}),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	booking, err = s.bookingRepo.FindByID(bookingID)
	if err != nil {
		s.logger.Error("Booking not found", zap.Error(err))
		return nil, errors.New("booking not found")
	}

	listing, err := s.listingRepo.FindByID(booking.ListingID)
	if err != nil {
		s.logger.Error("Listing not found", zap.Error(err))
		return nil, errors.New("listing not found")
	}

	if listing.OwnerID != ownerID {
		s.logger.Warn("Unauthorized booking response attempt")
		return nil, errors.New("unauthorized")
	}

	if booking.Status != models.BookingStatusEnumPending {
		s.logger.Warn("Booking cannot be responded to", zap.String("status", booking.Status))
		return nil, errors.New("booking cannot be responded to")
	}

	now := time.Now()
	booking.Status = req.Status
	booking.ResponseMessageFromOwner = req.ResponseMessage
	booking.RespondedAt = &now
	booking.RespondedBy = &ownerID

	err = s.bookingRepo.Update(booking)
	return
}

func (s *bookingService) CancelBooking(bookingID int64, userID int64, userRole string) (booking *models.Booking, err error) {
	defer func(start time.Time) {
		s.logger.Info("Cancel booking",
			zap.String("function", "CancelBooking"),
			zap.Any("params", gin.H{"booking_id": bookingID, "user_id": userID, "user_role": userRole}),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
	}(time.Now())

	booking, err = s.bookingRepo.FindByID(bookingID)
	if err != nil {
		s.logger.Error("Booking not found", zap.Error(err))
		return nil, errors.New("booking not found")
	}

	listing, err := s.listingRepo.FindByID(booking.ListingID)
	if err != nil {
		s.logger.Error("Listing not found", zap.Error(err))
		return nil, errors.New("listing not found")
	}

	isOwner := listing.OwnerID == userID
	isRenter := booking.RenterID == userID

	if !isOwner && !isRenter {
		s.logger.Warn("Unauthorized booking cancellation attempt")
		return nil, errors.New("unauthorized")
	}

	if booking.Status != models.BookingStatusEnumPending && booking.Status != models.BookingStatusEnumAccepted {
		s.logger.Warn("Booking cannot be cancelled", zap.String("status", booking.Status))
		return nil, errors.New("booking cannot be cancelled")
	}

	now := time.Now()
	booking.Status = models.BookingStatusEnumCancelled
	booking.RespondedAt = &now
	booking.RespondedBy = &userID

	err = s.bookingRepo.Update(booking)
	return
}

func (s *bookingService) AutoCompleteBookings() {
	defer func(start time.Time) {
		s.logger.Info("Auto-complete bookings cron job",
			zap.String("function", "AutoCompleteBookings"),
			zap.Duration("duration", time.Since(start)),
		)
	}(time.Now())

	bookings, err := s.bookingRepo.FindCompletableBookings()
	if err != nil {
		s.logger.Error("Failed to find completable bookings", zap.Error(err))
		return
	}

	s.logger.Info("Found bookings to auto-complete", zap.Int("count", len(bookings)))

	for _, booking := range bookings {
		booking.Status = models.BookingStatusEnumCompleted
		if err := s.bookingRepo.Update(&booking); err != nil {
			s.logger.Error("Failed to auto-complete booking", zap.Int64("booking_id", booking.ID), zap.Error(err))
		}
	}
}
