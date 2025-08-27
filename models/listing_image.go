package models

import "time"

type ListingImage struct {
	ID        uint `gorm:"primaryKey"`
	ListingID uint `gorm:"not null"`
	ImageURL  string `gorm:"type:varchar(255);not null"`
	IsMain    bool   `gorm:"default:false"`
	SortOrder int    `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}