package api

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Safayet-Shawn/Parking-Lot/db"
	"github.com/Safayet-Shawn/Parking-Lot/helper/api_client"
	"github.com/Safayet-Shawn/Parking-Lot/helper/redis"
	"github.com/go-chi/chi/v5"
)

func GetSlot(w http.ResponseWriter, r *http.Request) {
	UserLocation := chi.URLParam(r, "latlon")
	acamp, err := GetAvailableParkingSlots()
	if err != nil {
		log.Printf("Error on Getting Parking Slots  where and reason:%v", err)
		w.WriteHeader(http.StatusBadRequest)
		resp := api_client.Message(false, err.Error())
		api_client.Respond(w, resp)
		return
	}
	resp := api_client.Message(true, "get all available Parking Slot")
	LatLon := strings.Split(UserLocation, ":")
	userLat := LatLon[0]
	userLon := LatLon[1]

	min := 0.00
	var dis float64
	for i, v := range acamp {
		dis = Distance(userLat, userLon, v)
		if dis < min {
			dis = min

		}
		acamp[i].Distance = dis
	}
	var meg string
	for _, vl := range acamp {
		if dis == vl.Distance {
			meg = fmt.Sprintf("Nearest Slot Name:%v , Slot Row %v and Slot Column: %v", vl.Name, vl.Row, vl.Colum)
		}
	}
	resp["Instruction"] = meg
	resp["data"] = acamp
	api_client.Respond(w, resp)
}

func GetAvailableParkingSlots() ([]ParkingSlot, error) {
	var ret []ParkingSlot
	err := db.GetDB().Table("parking_slots").Where("status = ? ", 1).Find(&ret).Error
	if err != nil {
		return ret, err
	}
	// applog.Infof(r.Context(), "len of campaigns: %v", len(ret))
	return ret, nil
}
func Distance(UserLat, UserLon string, slot ParkingSlot) float64 {
	stringFloats := []string{UserLat, UserLon, slot.Lat, slot.Lon}

	// Convert strings to floats
	floats, err := stringsToFloats(stringFloats)
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}
	distance := Haversine(floats[0], floats[1], floats[2], floats[3])
	// Print the converted floats
	// fmt.Sprintf("Slot Name:%v , Slot Row %v and Slot Column: %v",slot.Name)
	slot.Distance = distance
	fmt.Println("============inside compare==========", slot.Distance)
	return distance

}
func Haversine(Lat1, Lon1, Lat2, Lon2 float64) float64 {
	const earthRadius = 6371 // in kilometers

	// Convert latitude and longitude from degrees to radians
	lat1 := degreesToRadians(Lat1)
	lon1 := degreesToRadians(Lon1)
	lat2 := degreesToRadians(Lat2)
	lon2 := degreesToRadians(Lon2)

	// Calculate differences
	dLat := lat2 - lat1
	dLon := lon2 - lon1

	// Calculate distance using Haversine formula
	a := math.Pow(math.Sin(dLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c

	return distance
}
func degreesToRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}
func stringsToFloats(strings []string) ([]float64, error) {
	floats := make([]float64, len(strings))
	for i, str := range strings {
		floatValue, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, err
		}
		floats[i] = floatValue
	}
	return floats, nil
}
func Park(w http.ResponseWriter, r *http.Request) {
	uid, _ := api_client.GetToken(r)

	c := bookSlot{}
	p := db.ParkingSlot{}
	//here err
	err := api_client.NewApiClient(r).JsonBind(&c)
	if err != nil {
		log.Printf("failed to bind json, Reason: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to unmarshal"))
		return
	}
	pid := chi.URLParam(r, "pid")
	sid := chi.URLParam(r, "sid")
	err = UpdateParkSlot(sid, p)
	if err != nil {
		log.Printf("failed to insert into db, Reason: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("unable to insert into db %v", err)))
		return
	}
	err = insertIntoBookSlot(c, pid, sid, uid)
	if err != nil {
		log.Printf("failed to insert into db, Reason: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("unable to insert into db %v", err)))
		return
	}
	resp := api_client.Message(true, fmt.Sprintf("succesfully Parkrd vehicle in Parking Slot %v", sid))
	api_client.Respond(w, resp)

}
func insertIntoBookSlot(c bookSlot, pid, sid string, uid int) error {
	// strconv.Atoi(str)
	p, _ := strconv.Atoi(pid)
	s, _ := strconv.Atoi(sid)
	// id, _ := strconv.Atoi(uid)
	currentTime := time.Now()
	currentDate := time.Now().Format("2006-01-02")
	slot := db.BookSlot{
		ParkingLotId:  p,
		ParkingSlotId: s,
		UserId:        uid,
		Vehicle:       c.Vehicle,
		Date:          currentDate,
		Status:        c.Status,
		CreatedAt:     currentTime,
		UpdatedAt:     currentTime,
	}
	epochTime := currentTime.Unix()
	ctx := context.Background()
	err := redis.GetRedis().Set(ctx, sid, epochTime, 0)
	if err != nil {
		log.Printf("Failed set Parking start time in redis where time was %v and err:%v", currentTime, err)
	}
	d := db.GetDB()
	return d.Create(&slot).Error
}
func Unpark(w http.ResponseWriter, r *http.Request) {
	sid := chi.URLParam(r, "sid")
	usr := db.BookSlot{}
	err := api_client.NewApiClient(r).JsonBind(&usr)
	if err != nil {
		log.Printf("failed to bind json, Reason: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		msg := api_client.Message(false, err.Error())
		api_client.Respond(w, msg)
		return
	}
	id, err := strconv.Atoi(sid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := api_client.Message(false, err.Error())
		api_client.Respond(w, msg)
		return
	}
	ctx := context.Background()
	val, err := redis.GetRedis().Get(ctx, sid).Result()
	if err != nil {
		log.Printf("failed to value from redis ", err)
		msg := api_client.Message(false, err.Error())
		api_client.Respond(w, msg)
		return
	}
	v, _ := strconv.ParseInt(val, 10, 64)
	StartTime := time.Unix(v, 0)
	parkedTime := ParkedTime(StartTime)
	cost := parkedTime * 10

	err = UnparkSlot(id, usr, parkedTime, cost)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		api_client.Respond(w, api_client.Message(false, err.Error()))
		return
	}
	err = redis.GetRedis().Del(context.Background(), sid).Err()
	if err != nil {
		log.Printf("failed to remove key from redis")
		msg := api_client.Message(false, err.Error())
		api_client.Respond(w, msg)
		return
	}
	msg := fmt.Sprintf("Sucessfully Unparker From Parking Lot\n You have been charged for %v hours and cost is : %v RS", parkedTime, cost)
	resp := api_client.Message(true, msg)

	api_client.Respond(w, resp)
}

func UnparkSlot(id int, usr db.BookSlot, durration, cost int) error {
	usr.Cost = cost
	usr.ParkingDuration = durration
	usr.Status = "1"
	err := db.GetDB().Table("book_slots").Select("id = ?", id).Where("id = ? ", id).Update(usr).Error
	if err != nil {
		return err
	}
	return nil
}
func UpdateParkSlot(id string, usr db.ParkingSlot) error {
	usr.Status = 2
	err := db.GetDB().Table("parking_slots").Select("id = ?", id).Where("id = ? ", id).Update(usr).Error
	if err != nil {
		return err
	}
	return nil
}
func ParkedTime(start time.Time) int {
	// Parse provided time string into a time.Time object
	providedTime, err := time.Parse("2006-01-02 15:04:05 -0700 -07", start.String())
	if err != nil {
		log.Println("Error in Time Parse:", err)
	}

	currentTime := time.Now()
	var d int
	// Calculate the time difference
	timeDifference := currentTime.Sub(providedTime)
	minutesDifference := float64(timeDifference.Minutes())
	duration := (minutesDifference / 60)
	FirstDigitAfterDecimal := firstDigitAfterDecimal(duration)
	if FirstDigitAfterDecimal > 0 {
		d = int(duration + 1)
	} else {
		d = int(duration)

	}

	return d
}
func firstDigitAfterDecimal(f float64) int {
	// Get the fractional part of the float number
	fractionalPart := math.Mod(math.Abs(f), 1)

	// Multiply fractional part by 10 to get the first digit after the decimal point
	firstDigit := int(math.Floor(fractionalPart * 10))

	return firstDigit
}

func UpdateSlot(w http.ResponseWriter, r *http.Request) {
	_, UserType := api_client.GetToken(r)
	if UserType != 1 {
		w.WriteHeader(http.StatusUnauthorized)
		resp := api_client.Message(false, "Sorry, You are Unauthorized, only Manager can Update Parking Slot status which is maintenance/working")
		api_client.Respond(w, resp)
	} else {
		sid := chi.URLParam(r, "sid")
		usr := db.ParkingSlot{}
		queryParams := r.URL.Query()
		status := queryParams.Get("status")
		err := api_client.NewApiClient(r).JsonBind(&usr)
		if err != nil {
			log.Printf("failed to bind json, Reason: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			msg := api_client.Message(false, err.Error())
			api_client.Respond(w, msg)
			return
		}
		id, err := strconv.Atoi(sid)
		if err != nil {
			log.Printf("failed to convert into int from string,where user Id: %v, Reason: %v", sid, err)
			w.WriteHeader(http.StatusBadRequest)
			msg := api_client.Message(false, err.Error())
			api_client.Respond(w, msg)
			return
		}
		var s int
		if status == "maintenance" {
			s = 3
		} else if status == "working" {
			s = 1
		} else {
			w.WriteHeader(http.StatusBadRequest)
			api_client.Respond(w, api_client.Message(false, fmt.Sprintf("You must choose between maintenance or working")))
			return
		}

		err = updateSlotStatus(id, usr, s)
		if err != nil {
			log.Printf("failed to update status,where Id: %v, Reason: %v", id, err)
			w.WriteHeader(http.StatusBadRequest)
			api_client.Respond(w, api_client.Message(false, err.Error()))
			return
		}

		resp := api_client.Message(true, fmt.Sprintf("sucessfully Updated Parking Slot: %v ", sid))
		api_client.Respond(w, resp)
	}
}

func updateSlotStatus(id int, usr db.ParkingSlot, status int) error {
	usr.Status = status
	err := db.GetDB().Table("parking_slots").Select("id = ?", id).Where("id = ? ", id).Update(usr).Error
	if err != nil {
		return err
	}
	return nil
}
func GetSlotInfo(w http.ResponseWriter, r *http.Request) {
	_, UserType := api_client.GetToken(r)
	if UserType != 1 {
		w.WriteHeader(http.StatusUnauthorized)
		resp := api_client.Message(false, "Sorry, You are Unauthorized, only Manager can Update Parking Slot status which is maintenance/working")
		api_client.Respond(w, resp)
	} else {
		queryParams := r.URL.Query()

		date := queryParams.Get("date")

		acamp, err := GetParkingSlots(date)
		if err != nil {
			log.Printf("Error on Getting All Parking Solt list   reason:%v", err)
			w.WriteHeader(http.StatusBadRequest)
			resp := api_client.Message(false, err.Error())
			api_client.Respond(w, resp)
			return
		}
		resp := api_client.Message(true, "Get All Parking Solt")
		var totalHour int
		var totalAmount int
		var totalVehicleParked = 0
		for _, v := range acamp {
			totalHour = totalHour + v.ParkingDuration
			totalAmount = totalAmount + v.Cost
			totalVehicleParked++
		}
		if date == "" {
			resp["data"] = acamp
			api_client.Respond(w, resp)
		} else {
			msg := fmt.Sprintf("total number of vehicles parked: %v, total parking time: %v and the total fee collected: %v  on  day: %v", totalVehicleParked, totalHour, totalAmount, date)
			resp["Info"] = msg
			resp["data"] = acamp
			api_client.Respond(w, resp)
		}
	}
}

func GetParkingSlots(date string) ([]bookSlot, error) {
	var ret []bookSlot
	if date == "" {
		err := db.GetDB().Table("book_slots").Find(&ret).Error
		if err != nil {
			return ret, err
		}
		return ret, nil

	} else {
		err := db.GetDB().Table("book_slots").Where("date = ? ", date).Find(&ret).Error
		if err != nil {
			return ret, err
		}
		return ret, nil
	}

}
