package models

import (
    "time"

    "gorm.io/gorm"
)

type Listing struct {
    ID          uint           `gorm:"primaryKey"`
    Title       string         `gorm:"type:varchar(255);not null"`
    Description string         `gorm:"type:text"`
    Price       float64        `gorm:"not null"`
    Address     string         `gorm:"type:varchar(255)"`
    City        string         `gorm:"type:varchar(100)"`
    OwnerID     uint           `gorm:"not null"` // FK đến User
    Owner       User           `gorm:"foreignKey:OwnerID"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`
}
