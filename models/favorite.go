package models

import "time"

type Favorite struct {
	UserID    uint `gorm:"primaryKey"`
	ListingID uint `gorm:"primaryKey"`
	CreatedAt time.Time
}