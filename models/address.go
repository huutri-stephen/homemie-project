package models

import "time"

type Address struct {
	ID           uint   `json:"id" gorm:"primary_key"`
	CityID       uint   `json:"city_id"`
	WardID       uint   `json:"ward_id"`
	AreaID       *uint  `json:"area_id"`
	Street       string `json:"street"`
	HouseNumber  string `json:"house_number"`
	BuildingName string `json:"building_name"`
	FloorNumber  int    `json:"floor_number"`
	RoomNumber   string `json:"room_number"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
