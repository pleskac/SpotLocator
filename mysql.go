package main

import (
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
	// _ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine
)

func AddGPS(longitude float32, latitude float32) {
	fmt.Println("NEW LOCATION:", longitude, latitude)

	db := mysql.New("tcp", "", "127.0.0.1:3306", "root", "rootroot", "gps")
	err := db.Connect()
	if err != nil {
		panic(err)
	}
	/*
		stmt, err := db.Prepare("INSERT INTO gps VALUES (?, ?)")

		stmt.Run("longitude", longitude)
	*/
	//NEED TO IMPLEMENT THIS!!

}
