package models

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
	ID                       uint `gorm:"primaryKey"`
	ListingID                uint `gorm:"not null"`
	RenterID                 uint `gorm:"not null"`
	ScheduledTime            time.Time `gorm:"not null"`
	Status                   string `gorm:"type:booking_status_enum;default:'pending'"`
	MessageFromRenter        string `gorm:"type:text"`
	ResponseMessageFromOwner string `gorm:"type:text"`
	RespondedAt              *time.Time
	RespondedBy              *uint
	CreatedAt                time.Time
	UpdatedAt                time.Time
}