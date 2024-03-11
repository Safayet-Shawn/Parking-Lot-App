package db

import "time"

type ParkingLot struct {
	Id        int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name      string `gorm:"size:100" json:"name,omitempty"`
	Lat       string `gorm:"type:varchar(100)" json:"lat",omitempty`
	Lon       string `gorm:"type:varchar(100)" json:"lon,omitempty"`
	Status    int    `gorm:"type:tinyint" json:"status,omitempty"`
	IsDeleted bool   `gorm:"type:bool" json:"payment,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
type ParkingSlot struct {
	Id           int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ParkingLotId int    `gorm:"type:int(100)" json:"parkinglot_id,omitempty"`
	Name         string `gorm:"type:varchar(100)" json:"name",omitempty`
	// UserId         string `gorm:"type:int(100)" json:"user_id,omitempty"`
	// Car         string `gorm:"type:varchar(100)" json:"name",omitempty`
	Lat    string `gorm:"type:varchar(100)" json:"lat",omitempty`
	Lon    string `gorm:"type:varchar(100)" json:"lon",omitempty`
	Row    int    `gorm:"type:int(100)" json:"row",omitempty`
	Column int    `gorm:"type:int(100)"json:"column",omitempty`
	Status int    `gorm:"type:int(100)"json:"status",omitempty`
}

// 3 status type for Parking slot
// 1=>active
// 2=>engaged
// 3=>underMaitainence

type BookSlot struct {
	Id              int       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ParkingLotId    int       `gorm:"type:int(100)" json:"parkinglot_id,omitempty"`
	ParkingSlotId   int       `gorm:"type:int(100)" json:"parkingslot_id,omitempty"`
	Date            string    `gorm:"type:varchar(30)" json:"date,not null"`
	UserId          int       `gorm:"type:int(100)" json:"user_id,omitempty"`
	Vehicle         string    `gorm:"type:varchar(100)" json:"vehicle,omitempty"`
	ParkingDuration int       `gorm:"type:int(100)" json:"parkinh_duration,omitempty"`
	Cost            int       `gorm:"type:int(100)" json:"cost,omitempty"`
	Status          string    `gorm:"type:varchar(100)" json:"status,omitempty"`
	IsDeleted       bool      `gorm:"type:bool" json:"is_deleted,omitempty"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
}
