package models

import "time"

type Booking struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"not null"` // Người thuê
    ListingID uint      `gorm:"not null"`
    Schedule  time.Time `gorm:"not null"` // Lịch hẹn xem
    Status    string    `gorm:"type:varchar(20);default:'pending'"` // pending, accepted, rejected

    User    User    `gorm:"foreignKey:UserID"`
    Listing Listing `gorm:"foreignKey:ListingID"`
}
