package models

import (
    "time"

    "gorm.io/gorm"
)

type Booking struct {
    ID        uint           `gorm:"primaryKey"`
    UserID    uint           `gorm:"not null"`
    User      User           `gorm:"foreignKey:UserID"`
    ListingID uint           `gorm:"not null"`
    Listing   Listing        `gorm:"foreignKey:ListingID"`
    StartDate time.Time      `gorm:"not null"`
    EndDate   time.Time      `gorm:"not null"`
    Status    string         `gorm:"type:varchar(20);default:'pending'"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
