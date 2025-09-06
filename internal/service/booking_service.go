package service

import (
	"errors"
	"homemie/internal/domain"
	"homemie/models/dto"
	"homemie/models/request"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BookingService struct {
	bookingRepo domain.BookingRepository
	listingRepo domain.ListingRepository
	logger      *zap.Logger
}

func NewBookingService(bookingRepo domain.BookingRepository, listingRepo domain.ListingRepository, logger *zap.Logger) *BookingService {
	return &BookingService{bookingRepo, listingRepo, logger}
}

func (s *BookingService) CreateBooking(renterID int64, input request.CreateBookingRequest) (booking *dto.Booking, err error) {
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

	booking = &dto.Booking{
		ListingID:         input.ListingID,
		RenterID:          renterID,
		ScheduledTime:     scheduledTime,
		MessageFromRenter: input.MessageFromRenter,
		Status:            dto.BookingStatusPending,
	}

	err = s.bookingRepo.Create(booking)
	return
}

func (s *BookingService) GetMyBookings(userID int64) (bookings []dto.Booking, err error) {
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

func (s *BookingService) GetOwnerBookings(ownerID int64) (bookings []dto.Booking, err error) {
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

func (s *BookingService) RespondToBooking(bookingID int64, ownerID int64, req request.RespondBookingRequest) (booking *dto.Booking, err error) {
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

	if booking.Status != dto.BookingStatusPending {
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

func (s *BookingService) CancelBooking(bookingID int64, userID int64, userRole string) (booking *dto.Booking, err error) {
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

	if booking.Status != dto.BookingStatusPending && booking.Status != dto.BookingStatusAccepted {
		s.logger.Warn("Booking cannot be cancelled", zap.String("status", booking.Status))
		return nil, errors.New("booking cannot be cancelled")
	}

	now := time.Now()
	booking.Status = dto.BookingStatusCancelled
	booking.RespondedAt = &now
	booking.RespondedBy = &userID

	err = s.bookingRepo.Update(booking)
	return
}

func (s *BookingService) AutoCompleteBookings() {
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
		booking.Status = dto.BookingStatusCompleted
		if err := s.bookingRepo.Update(&booking); err != nil {
			s.logger.Error("Failed to auto-complete booking", zap.Int64("booking_id", booking.ID), zap.Error(err))
		}
	}
}