package request

import "time"

type ListingImageRequest struct {
	ID        int64  `gorm:"primaryKey"`
	ListingID int64  `gorm:"not null"`
	ImageURL  string `gorm:"type:varchar(500);not null"`
	IsMain    bool   `gorm:"default:false"`
	SortOrder int32  `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
