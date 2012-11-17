package main

import (
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"time"
)

type Location struct {
	//other info
	Longitude float64
	Latitude  float64
	Title     string
	Details   string
}

type Trip struct {
	TripName    string
	Zoom        int
	CenterLong  float64
	CenterLat   float64
	Coordinates []Location
}

func Connect() mysql.Conn {
	//set up database connection
	db := mysql.New("tcp", "", "127.0.0.1:3306", "root", "rootroot", "gps")
	err := db.Connect()
	if err != nil {
		fmt.Println("ERROR CONNECTING:", err)
		panic(err)
	}

	return db
}

func GetLatestId() int {
	db := Connect()
	defer db.Close()

	rows, _, err := db.Query("SELECT id FROM latestTweet")
	if err != nil {
		panic(err)
	}

	if len(rows) != 1 {
		return -1
	}

	return rows[0].Int(0)
}

func SaveLatestId(id int) {
	db := Connect()
	defer db.Close()

	//delete all rows
	stmt, err := db.Prepare("DELETE FROM latestTweet")
	_, err = stmt.Run()
	if err != nil {
		panic(err)
	}

	//insert new row
	stmt, err = db.Prepare("INSERT INTO latestTweet (id) VALUES (?)")
	_, err = stmt.Run(id)
	if err != nil {
		panic(err)
	}
}

func AddGPS(longitude float64, latitude float64, message string, time int64) {
	db := Connect()
	defer db.Close()

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
		stmt, err := db.Prepare("INSERT INTO gps (longitude, latitude, details, timestamp) VALUES (?, ?, ?, ?)")
		_, err = stmt.Run(longitude, latitude, message, time)
		if err != nil {
			panic(err)
		}
	} else {
		stmt, err := db.Prepare("INSERT INTO gps (longitude, latitude, details, trip, timestamp) VALUES (?, ?, ?, ?, ?)")
		_, err = stmt.Run(longitude, latitude, message, tripKey, time)
		if err != nil {
			panic(err)
		}
	}

}

func CreateTrip(name string) {
	db := Connect()
	defer db.Close()

	//End all trips
	EndTrips()

	//Create new trip, set it as current trip
	fmt.Println("starting trip!")

	//insert
	fmt.Println("Name:", name)
	stmt, err := db.Prepare("INSERT INTO trips (name, details, current) VALUES (?, ?, ?)")
	_, err = stmt.Run(name, "", 1)
	if err != nil {
		panic(err)
	}
}

func EndTrips() {
	db := Connect()
	defer db.Close()

	fmt.Println("Ending all trips")

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
	db := Connect()
	defer db.Close()

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
	myTrip := Trip{name, 10, -96.7, 40.8, nil}

	//Get the GPS coordinates of that trip
	id := (rows[0]).Str(0)
	fmt.Println("select * from gps where trip =", id)
	rows, _, err = db.Query("select * from gps where trip = " + id)
	if err != nil {
		panic(err)
	}

	latLow := 90.0     //the MAX lat value
	latHigh := -90.0   //the MIN lat value
	longLow := 180.0   //the MAX long value
	longHigh := -180.0 //the MIN long value

	//Add every GPS location
	for _, row := range rows {
		//Timestamp: Details
		mytime := time.Unix(row.Int64(4), 0)
		year, month, day := mytime.Date()
		hour, min, sec := mytime.Clock()

		//fmt.Println(mytime.Weekday(), month, day, year, hour, min, sec
		timestamp := fmt.Sprintf("%s, %s %d, %d at %d:%d:%d", mytime.Weekday().String(), month, day, year, hour, min, sec)
		details := timestamp //+ " – " + row.Str(6)

		myTrip.Coordinates = append(myTrip.Coordinates, Location{row.Float(2), row.Float(3), row.Str(5), details})

		//CENTER AND SCALE THE MAP
		//long
		if longLow > row.Float(2) {
			longLow = row.Float(2)
		}
		if longHigh < row.Float(2) {
			longHigh = row.Float(2)
		}

		//lat
		if latLow > row.Float(3) {
			latLow = row.Float(3)
		}
		if latHigh < row.Float(3) {
			latHigh = row.Float(3)
		}
	}

	totalDistance := longHigh - longLow

	if (latHigh - latLow) > (longHigh - longLow) {
		totalDistance = latHigh - latLow
	}
	if totalDistance < 0.05 {
		myTrip.Zoom = 15
	} else if totalDistance < 0.125 {
		myTrip.Zoom = 12
	} else if totalDistance < 0.25 {
		myTrip.Zoom = 11
	} else if totalDistance < 0.5 {
		myTrip.Zoom = 10
	} else if totalDistance < 1.1 {
		myTrip.Zoom = 9
	} else if totalDistance < 2.2 {
		myTrip.Zoom = 8
	} else if totalDistance < 4.5 {
		myTrip.Zoom = 7
	} else if totalDistance < 9 {
		myTrip.Zoom = 6
	} else if totalDistance < 17 {
		myTrip.Zoom = 5
	} else if totalDistance < 34 {
		myTrip.Zoom = 4
	} else {
		myTrip.Zoom = 3
	}

	averageLong := (longLow + longHigh) / 2
	averageLat := (latLow + latHigh) / 2

	myTrip.CenterLat = averageLat
	myTrip.CenterLong = averageLong
	//myTrip.CenterLong = fmt.Sprintf("%f", averageLong)
	//myTrip.CenterLat = fmt.Sprintf("%f", averageLat)

	return myTrip
}
