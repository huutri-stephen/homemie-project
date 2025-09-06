package request

type CreateBookingRequest struct {
	ListingID      int64  `json:"listing_id" binding:"required"`
	ScheduledTime  string `json:"scheduled_time" binding:"required"`
	MessageFromRenter string `json:"message_from_renter"`
}

type RespondBookingRequest struct {
	Status          string `json:"status" binding:"required,oneof=accepted rejected"`
	ResponseMessage string `json:"response_message_from_owner"`
}
