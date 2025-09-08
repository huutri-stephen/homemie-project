package request

import "time"

type AddressRequest struct {
	ID           int64     `json:"id" gorm:"primary_key"`
	CityID       int64     `json:"city_id"`
	WardID       int64     `json:"ward_id"`
	AreaID       *int64    `json:"area_id"`
	Street       string    `json:"street"`
	HouseNumber  string    `json:"house_number"`
	BuildingName string    `json:"building_name"`
	FloorNumber  int32     `json:"floor_number"`
	RoomNumber   string    `json:"room_number"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
