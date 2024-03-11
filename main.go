package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Safayet-Shawn/Parking-Lot/db"
	"github.com/Safayet-Shawn/Parking-Lot/route"
)

func main() {
	r := route.NewRouter()
	err := dbTableCreateIfNotExist()
	if err != nil {
		log.Fatalln("failed to create db tables")
	}
	fmt.Println("Starting Praking Lot App Server at Port: 8080 ...")

	http.ListenAndServe(":8080", r)
}
func dbTableCreateIfNotExist() error {
	db.InitDB()
	db.AutoMigrateDatabase()

	return nil
}
