package main

import (
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
	// _ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine
)

//Need to add more parameters to insert into DB
func AddGPS(longitude float32, latitude float32) {
	fmt.Println("NEW LOCATION:", longitude, latitude)

	db := mysql.New("tcp", "", "127.0.0.1:3306", "root", "rootroot", "gps")
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	//Need to add other column info
	stmt, err := db.Prepare("INSERT INTO gps (longitude, latitude) VALUES (?, ?)")

	stmt.Run(longitude, latitude)

}
