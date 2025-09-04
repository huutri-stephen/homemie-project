package request

import (
	"time"
)

type CreateBookingRequest struct {
	ListingID 			int64   	`json:"listing_id" binding:"required"`
	MessageFromRenter 	string  	`json:"message_from_renter"`
	ScheduledTime       time.Time 	`json:"start_date" binding:"required"`
}