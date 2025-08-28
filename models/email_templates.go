package models

import (
    "time"
)

type EmailTemplate struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Name      string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
    Subject   string    `gorm:"type:varchar(255);not null" json:"subject"`
    Body      string    `gorm:"type:text;not null" json:"body"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
