package models

import (
    "time"

    "gorm.io/gorm"
)

type User struct {
    ID        uint           `gorm:"primaryKey"`
    Name      string         `gorm:"type:varchar(100);not null"`
    Email     string         `gorm:"uniqueIndex;type:varchar(100);not null"`
    Password  string         `gorm:"type:varchar(255);not null"`
    Phone     string         `gorm:"type:varchar(20)"`
    Role      string         `gorm:"type:varchar(20);default:'renter'"` // renter | owner | admin
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
