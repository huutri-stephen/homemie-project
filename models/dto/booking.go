package dto

import (
	"time"
)

const (
	BookingStatusPending   = "pending"
	BookingStatusAccepted  = "accepted"
	BookingStatusRejected  = "rejected"
	BookingStatusCancelled = "cancelled"
	BookingStatusCompleted = "completed"
)

type Booking struct {
	ID                       int64     `gorm:"primaryKey"`
	ListingID                int64     `gorm:"not null"`
	RenterID                 int64     `gorm:"not null"`
	ScheduledTime            time.Time `gorm:"not null"`
	Status                   string    `gorm:"type:booking_status_enum;default:'pending'"`
	MessageFromRenter        string    `gorm:"type:text"`
	ResponseMessageFromOwner string    `gorm:"type:text"`
	RespondedAt              *time.Time
	RespondedBy              *int64
	CreatedAt                time.Time
	UpdatedAt                time.Time
}
