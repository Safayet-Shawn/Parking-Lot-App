package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Safayet-Shawn/Parking-Lot/db"
	"github.com/Safayet-Shawn/Parking-Lot/helper/api_client"
	"github.com/go-chi/chi/v5"
)

func CreateParking(w http.ResponseWriter, r *http.Request) {
	_, UserType := api_client.GetToken(r)
	if UserType != 1 {
		w.WriteHeader(http.StatusUnauthorized)
		resp := api_client.Message(false, "Sorry, You are Unauthorized, only Manager create parking")
		api_client.Respond(w, resp)
	} else {
		c := ParkingLot{}
		err := api_client.NewApiClient(r).JsonBind(&c)
		if err != nil {
			log.Printf("failed to bind json, Reason: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("unable to unmarshal"))
			return
		}
		err = insertIntoDb(c)
		if err != nil {
			log.Printf("failed to insert into db, Reason: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("unable to insert into db %v", err)))
			return
		}
		resp := api_client.Message(true, "succesfully created Parking Lot")
		api_client.Respond(w, resp)
	}

}

func insertIntoDb(c ParkingLot) error {

	lot := db.ParkingLot{
		Name:   c.Name,
		Lat:    c.Lat,
		Lon:    c.Lon,
		Status: c.Status,
	}
	d := db.GetDB()
	return d.Create(&lot).Error
}
func CreateParkingSlot(w http.ResponseWriter, r *http.Request) {
	_, UserType := api_client.GetToken(r)
	if UserType != 1 {
		w.WriteHeader(http.StatusUnauthorized)
		resp := api_client.Message(false, "Sorry, You are Unauthorized, only Manager create Parking Slot")
		api_client.Respond(w, resp)
	} else {
		c := ParkingSlot{}
		parkingLotID := chi.URLParam(r, "id")
		id, _ := strconv.Atoi(parkingLotID)
		err := api_client.NewApiClient(r).JsonBind(&c)
		if err != nil {
			log.Printf("failed to bind json, Reason: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("unable to unmarshal"))
			return
		}
		err = insertIntoSlot(c, id)
		if err != nil {
			log.Printf("failed to insert into db, Reason: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("unable to insert into db %v", err)))
			return
		}
		resp := api_client.Message(true, "succesfully created Parking Slot")
		api_client.Respond(w, resp)
	}
}
func insertIntoSlot(c ParkingSlot, lotId int) error {
	lot := db.ParkingSlot{
		ParkingLotId: lotId,
		Name:         c.Name,
		Lat:          c.Lat,
		Lon:          c.Lon,
		Row:          c.Row,
		Column:       c.Colum,
		Status:       c.Status,
	}
	d := db.GetDB()
	return d.Create(&lot).Error
}
