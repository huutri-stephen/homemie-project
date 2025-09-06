package service

import (
	"errors"
	"homemie/internal/domain"
	"homemie/models/dto"
	"homemie/models/request"
	"time"
)

type BookingService struct {
	bookingRepo domain.BookingRepository
	listingRepo domain.ListingRepository
}

func NewBookingService(bookingRepo domain.BookingRepository, listingRepo domain.ListingRepository) *BookingService {
	return &BookingService{bookingRepo, listingRepo}
}

func (s *BookingService) CreateBooking(renterID int64, input request.CreateBookingRequest) (*dto.Booking, error) {
	scheduledTime, err := time.Parse(time.RFC3339, input.ScheduledTime)
	if err != nil {
		return nil, errors.New("invalid scheduled_time format")
	}

	booking := &dto.Booking{
		ListingID:         input.ListingID,
		RenterID:          renterID,
		ScheduledTime:     scheduledTime,
		MessageFromRenter: input.MessageFromRenter,
		Status:            dto.BookingStatusPending,
	}

	err = s.bookingRepo.Create(booking)
	return booking, err
}

func (s *BookingService) GetMyBookings(userID int64) ([]dto.Booking, error) {
	return s.bookingRepo.FindByUserID(userID)
}

func (s *BookingService) GetOwnerBookings(ownerID int64) ([]dto.Booking, error) {
	return s.bookingRepo.FindByOwnerID(ownerID)
}

func (s *BookingService) RespondToBooking(bookingID int64, ownerID int64, req request.RespondBookingRequest) (*dto.Booking, error) {
	booking, err := s.bookingRepo.FindByID(bookingID)
	if err != nil {
		return nil, errors.New("booking not found")
	}

	listing, err := s.listingRepo.FindByID(booking.ListingID)
	if err != nil {
		return nil, errors.New("listing not found")
	}

	if listing.OwnerID != ownerID {
		return nil, errors.New("unauthorized")
	}

	if booking.Status != dto.BookingStatusPending {
		return nil, errors.New("booking cannot be responded to")
	}

	now := time.Now()
	booking.Status = req.Status
	booking.ResponseMessageFromOwner = req.ResponseMessage
	booking.RespondedAt = &now
	booking.RespondedBy = &ownerID

	if err := s.bookingRepo.Update(booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *BookingService) CancelBooking(bookingID int64, userID int64, userRole string) (*dto.Booking, error) {
	booking, err := s.bookingRepo.FindByID(bookingID)
	if err != nil {
		return nil, errors.New("booking not found")
	}

	listing, err := s.listingRepo.FindByID(booking.ListingID)
	if err != nil {
		return nil, errors.New("listing not found")
	}

	isOwner := listing.OwnerID == userID
	isRenter := booking.RenterID == userID

	if !isOwner && !isRenter {
		return nil, errors.New("unauthorized")
	}

	if booking.Status != dto.BookingStatusPending && booking.Status != dto.BookingStatusAccepted {
		return nil, errors.New("booking cannot be cancelled")
	}

	now := time.Now()
	booking.Status = dto.BookingStatusCancelled
	booking.RespondedAt = &now
	booking.RespondedBy = &userID

	if err := s.bookingRepo.Update(booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *BookingService) AutoCompleteBookings() {
	bookings, err := s.bookingRepo.FindCompletableBookings()
	if err != nil {
		// Log the error
		return
	}

	for _, booking := range bookings {
		booking.Status = dto.BookingStatusCompleted
		s.bookingRepo.Update(&booking)
	}
}
