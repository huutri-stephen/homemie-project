package models

import (
	"time"
)

type Listing struct {
	ID              uint `gorm:"primaryKey"`
	OwnerID         uint `gorm:"not null"`
	Title           string `gorm:"type:varchar(255);not null"`
	Description     string `gorm:"type:text"`
	PropertyType    string `gorm:"type:property_type_enum;default:'rented_room'"`
	IsShared        bool   `gorm:"default:false"`
	Price           float64 `gorm:"type:decimal(15,2);not null"`
	AreaM2          float64 `gorm:"type:decimal(10,2)"`
	AddressID       uint `gorm:"not null"`
	ContactPhone    string `gorm:"type:varchar(20)"`
	ContactEmail    string `gorm:"type:varchar(255)"`
	ContactName     string `gorm:"type:varchar(100)"`
	NumBedrooms     int
	NumBathrooms    int
	NumFloors       int
	HasBalcony      bool `gorm:"default:false"`
	HasParking      bool `gorm:"default:false"`
	Amenities       string `gorm:"type:json"`
	PetAllowed      bool   `gorm:"default:false"`
	AllowedPetTypes string `gorm:"type:json"`
	Latitude        float64 `gorm:"type:decimal(10,8)"`
	Longitude       float64 `gorm:"type:decimal(10,8)"`
	ListingType     string `gorm:"type:listing_type_enum;default:'for_rent'"`
	DepositAmount   float64 `gorm:"type:decimal(15,2)"`
	Status          string `gorm:"type:listing_status_enum;default:'pending'"`
	IsFeatured      bool `gorm:"default:false"`
	ViewCount       int `gorm:"default:0"`
	PublishedAt     *time.Time
	ExpiresAt       *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}