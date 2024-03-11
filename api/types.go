package api

import "time"

type User struct {
	UserType int    `json:"user_type,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Address  string `json:"address,omitempty"`
	// Country        string `json:"country,omitempty"`
	Password       string `json:"password,omitempty"`
	RepeatPassword string `json:"repeat_password,omitempty"`
}
type Login struct {
	UserType int    `json:"user_type,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Identity string `json:"identity,omitempty"`
	Password string `json:"password,omitempty"`
}
type ParkingLot struct {
	Name string `json:"name,omitempty"`
	Lat  string `json:"lat,omitempty"`
	Lon  string `json:"lon,omitempty"`

	Status int `json:"status,omitempty"`
}
type ParkingSlot struct {
	ParkingLotId int     `json:"parkinglot_id,omitempty"`
	Name         string  `json:"name,omitempty"`
	Lat          string  `json:"lat,omitempty"`
	Lon          string  `json:"lon,omitempty"`
	Row          int     `json:"row,omitempty"`
	Colum        int     `json:"colum,omitempty"`
	Status       int     `json:"status,omitempty"`
	Distance     float64 `json:"distance,omitempty"`
}
type bookSlot struct {
	ParkingLotId    int    `json:"parkinglot_id,omitempty"`
	ParkingSlotId   int    `json:"parkingslot_id,omitempty"`
	UserId          string `json:"user_id,omitempty"` // have to get from auth
	Date            string `json:"date,not null"`
	Vehicle         string `json:"vehicle",omitempty`
	ParkingDuration int    `json:"parking_duration",omitempty`

	Cost      int    `json:"cost",omitempty`
	Status    string `json:"status",omitempty` //its for parked/unparked
	IsDeleted bool   `json:"payment,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
