package dto

import (
	"time"
)

type Address struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	CityID       int64     `gorm:"not null;index"`
	WardID       int64     `gorm:"not null;index"`
	AreaID       *int64    `gorm:"index"`
	Street       string    `gorm:"type:varchar(100)"`
	HouseNumber  string    `gorm:"type:varchar(50)"`
	BuildingName string    `gorm:"type:varchar(100)"`
	FloorNumber  int32     `gorm:""`
	RoomNumber   string    `gorm:"type:varchar(20)"`
	Latitude     float64   `gorm:"type:decimal(10,8)"`
	Longitude    float64   `gorm:"type:decimal(10,8)"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}
