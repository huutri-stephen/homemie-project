package service

import (
	"homemie/internal/domain"
	"homemie/models"
	"time"
)

type BookingService struct {
	repo domain.BookingRepository
}

func NewBookingService(repo domain.BookingRepository) *BookingService {
	return &BookingService{repo}
}

type CreateBookingInput struct {
	UserID    uint
	ListingID uint
	StartDate time.Time
	EndDate   time.Time
}

func (s *BookingService) CreateBooking(input CreateBookingInput) (*models.Booking, error) {
	booking := &models.Booking{
		// UserID:    input.UserID,
		// ListingID: input.ListingID,
		// StartDate: input.StartDate,
		// EndDate:   input.EndDate,
		// Status:    models.BookingStatusPending,
	}
	err := s.repo.Create(booking)
	return booking, err
}

func (s *BookingService) GetMyBookings(userID uint) ([]models.Booking, error) {
	return s.repo.FindByUserID(userID)
}

func (s *BookingService) GetOwnerBookings(ownerID uint) ([]models.Booking, error) {
	return s.repo.FindByOwnerID(ownerID)
}
