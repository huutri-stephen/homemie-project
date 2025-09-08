package dto

import "time"

type ListingImage struct {
	ID        int64  `gorm:"primaryKey"`
	ListingID int64  `gorm:"not null"`
	ImageURL  string `gorm:"type:varchar(500);not null"` // Update to 500 characters
	IsMain    bool   `gorm:"default:false"`
	SortOrder int32  `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
