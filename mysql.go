package main

import (
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
	// _ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine
)

//Need to add more parameters to insert into DB
func AddGPS(longitude float32, latitude float32, message string) {
	fmt.Println("NEW LOCATION:", longitude, latitude, message)

	db := mysql.New("tcp", "", "127.0.0.1:3306", "root", "rootroot", "gps")
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	//Get the current trip
	rows, _, err := db.Query("select id from trips where current = 1")
	if err != nil {
		panic(err)
	}

	var tripKey int
	if len(rows) > 1 {
		fmt.Println("More than one row!!")
	} else if len(rows) == 0 {
		fmt.Println("0 rows!")
	} else {
		//get dat foreign key to dat trip]

		tripKey = (rows[0]).Int(0)
	}
	fmt.Println("Trip Key:", tripKey)
	//Add the GPS row
	stmt, err := db.Prepare("INSERT INTO gps (longitude, latitude, details) VALUES (?, ?, ?)")
	_, err = stmt.Run(longitude, latitude, message)
	if err != nil {
		panic(err)
	}

}
