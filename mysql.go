package main

import (
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
	// _ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine
)

type Location struct {
	//other info
	Longitude string
	Latitude  string
}

type Trip struct {
	TripName    string
	Coordinates []Location
}

func GetLatestTweet() string {
	//set up database connection
	db := mysql.New("tcp", "", "127.0.0.1:3306", "root", "rootroot", "gps")
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	//query
	rows, _, err := db.Query("select id from latestTweet")
	if err != nil {
		panic(err)
	}

	if len(rows) > 1 {
		fmt.Println("More than one row!! WRONG!")

	} else if len(rows) == 0 {
		fmt.Println("0 rows!")
		return ""
	} else {
		//get dat foreign key to dat trip
		fmt.Println("Latest Tweet:", (rows[0]).Str(0))
		return (rows[0]).Str(0)
	}

	return ""
}

func SaveLatestTweet(tweetId string) {
	db := mysql.New("tcp", "", "127.0.0.1:3306", "root", "rootroot", "gps")
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	//delete all rows
	stmt, err := db.Prepare("DELETE FROM latestTweet")
	_, err = stmt.Run()
	if err != nil {
		panic(err)
	}

	//insert new row
	stmt, err = db.Prepare("INSERT INTO latestTweet (id) VALUES (?)")
	_, err = stmt.Run(latestId)
	if err != nil {
		panic(err)
	}
}

//Need to add more parameters to insert into DB
func AddGPS(longitude float32, latitude float32, message string) {
	fmt.Println("NEW LOCATION:", longitude, latitude, message)

	db := mysql.New("tcp", "", "127.0.0.1:3306", "root", "rootroot", "gps")
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	//Get the current trip, if it exists
	rows, _, err := db.Query("select id from trips where current = 1")
	if err != nil {
		panic(err)
	}
	tripKey := -1
	if len(rows) > 1 {
		fmt.Println("More than one row!!")
	} else if len(rows) == 0 {
		fmt.Println("0 rows!")
	} else {
		//get dat foreign key to dat trip
		tripKey = (rows[0]).Int(0)
	}

	fmt.Println("Trip Key:", tripKey)

	//Add the GPS row
	if tripKey == -1 {
		stmt, err := db.Prepare("INSERT INTO gps (longitude, latitude, details) VALUES (?, ?, ?)")
		_, err = stmt.Run(longitude, latitude, message)
		if err != nil {
			panic(err)
		}
	} else {
		stmt, err := db.Prepare("INSERT INTO gps (longitude, latitude, details, trip) VALUES (?, ?, ?, ?)")
		_, err = stmt.Run(longitude, latitude, message, tripKey)
		if err != nil {
			panic(err)
		}
	}

}

func CreateTrip(name string) {
	//End all trips
	EndTrips()

	//Create new trip, set it as current trip
	fmt.Println("starting trip!")

	//set up database connection
	db := mysql.New("tcp", "", "127.0.0.1:3306", "root", "rootroot", "gps")
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	//insert
	fmt.Println("Name:", name)
	stmt, err := db.Prepare("INSERT INTO trips (name, details, current) VALUES (?, ?, ?)")
	_, err = stmt.Run(name, "", 1)
	if err != nil {
		panic(err)
	}
}

func EndTrips() {
	fmt.Println("Ending all trips")
	//end all trips
	db := mysql.New("tcp", "", "127.0.0.1:3306", "root", "rootroot", "gps")
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	rows, _, err := db.Query("SELECT id FROM trips WHERE current = 1")
	if err != nil {
		panic(err)
	}
	for _, row := range rows {
		fmt.Println("REMOVING A CURRENT TRIP!")
		tripId := row.Str(0)
		stmt, err := db.Prepare("UPDATE trips SET current = 0 WHERE id=" + tripId)
		_, err = stmt.Run()
		if err != nil {
			panic(err)
		}
	}
}

func GetCurrentTrip() Trip {
	db := mysql.New("tcp", "", "127.0.0.1:3306", "root", "rootroot", "gps")
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	//Get the current trip, if it exists
	rows, _, err := db.Query("select * from trips where current = 1")
	if err != nil {
		panic(err)
	}

	if len(rows) > 1 {
		fmt.Println("More than one row!! WRONG!")
		return Trip{}
	} else if len(rows) == 0 {
		fmt.Println("0 rows! No current trip to return")
		return Trip{}
	}

	name := (rows[0]).Str(1)
	myTrip := Trip{name, nil}

	//Get the GPS coordinates of that trip
	id := (rows[0]).Int(0)
	rows, _, err = db.Query("select * from gps where trip = " + id)
	if err != nil {
		panic(err)
	}

	for _, row := range rows {
		myTrip.Coordinates = append(myTrip.Coordinates, Location{row.Str(2), row.Str(3)})
	}

	return myTrip
}
