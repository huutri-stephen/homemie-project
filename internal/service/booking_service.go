package service

import (
	"homemie/internal/domain"
	"homemie/models/dto"
)

type BookingService struct {
	repo domain.BookingRepository
}

func NewBookingService(repo domain.BookingRepository) *BookingService {
	return &BookingService{repo}
}

func (s *BookingService) CreateBooking(input dto.Booking) (*dto.Booking, error) {
	err := s.repo.Create(&input)
	return &input, err
}

func (s *BookingService) GetMyBookings(userID int64) ([]dto.Booking, error) {
	return s.repo.FindByUserID(userID)
}

func (s *BookingService) GetOwnerBookings(ownerID int64) ([]dto.Booking, error) {
	return s.repo.FindByOwnerID(ownerID)
}
