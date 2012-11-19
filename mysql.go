package main

import (
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	"time"
)

//Every GPS location for a trip
type Location struct {
	Longitude float64
	Latitude  float64
	Title     string
	Details   string
	Color     string
}

//A single trip for a single map
type Trip struct {
	TripName    string
	Zoom        int
	CenterLong  float64
	CenterLat   float64
	Coordinates []Location
}

func Connect() mysql.Conn {
	//Set up database connection
	db := mysql.New("tcp", "", "127.0.0.1:3306", "root", "rootroot", "gps")
	err := db.Connect()
	if err != nil {
		fmt.Println("ERROR CONNECTING:", err)
		panic(err)
	}

	return db
}

//TODO: "latestTweet" is an outdated name. No longer using Twitter
func GetLatestId() int {
	db := Connect()
	defer db.Close()

	rows, _, err := db.Query("SELECT id FROM latestTweet")
	if err != nil {
		panic(err)
	}

	if len(rows) < 1 {
		return -1
	} else if len(rows) > 1 {
		//delete all rows, the table is messed up
		stmt, err := db.Prepare("DELETE FROM latestTweet")
		_, err = stmt.Run()
		if err != nil {
			panic(err)
		}
	}

	return rows[0].Int(0)
}

//TODO: "latestTweet" is an outdated name. No longer using Twitter
func SaveLatestId(id int) {
	db := Connect()
	defer db.Close()

	//Delete all rows
	stmt, err := db.Prepare("DELETE FROM latestTweet")
	_, err = stmt.Run()
	if err != nil {
		panic(err)
	}

	//Insert new row
	stmt, err = db.Prepare("INSERT INTO latestTweet (id) VALUES (?)")
	_, err = stmt.Run(id)
	if err != nil {
		panic(err)
	}
}

func AddGPS(longitude, latitude float64, message, msgType string, time int64) {
	db := Connect()
	defer db.Close()

	//Get the current trip, if it exists
	rows, _, err := db.Query("select id from trips where current = 1")
	if err != nil {
		panic(err)
	}

	//Get the foreign key to the current trip
	tripKey := -1
	if len(rows) > 1 {
		fmt.Println("More than one row!!")
	} else if len(rows) == 0 {
		fmt.Println("0 rows!")
	} else {
		tripKey = (rows[0]).Int(0)
	}

	//Add the GPS row with data
	if tripKey == -1 {
		stmt, err := db.Prepare("INSERT INTO gps (longitude, latitude, details, timestamp, type) VALUES (?, ?, ?, ?, ?)")
		_, err = stmt.Run(longitude, latitude, message, time, msgType)
		if err != nil {
			panic(err)
		}
	} else {
		stmt, err := db.Prepare("INSERT INTO gps (longitude, latitude, details, trip, timestamp, type) VALUES (?, ?, ?, ?, ?, ?)")
		_, err = stmt.Run(longitude, latitude, message, tripKey, time, msgType)
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
	fmt.Println("Starting trip", name)

	//Insert new trip
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
		tripId := row.Str(0)
		stmt, err := db.Prepare("UPDATE trips SET current = 0 WHERE id=" + tripId)
		_, err = stmt.Run()
		if err != nil {
			panic(err)
		}
	}
}

//TODO: this can be split up and organized better
//Also, could default to the current trip, or allow specific trips to be returned, that would allow multiple maps on the site
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

		checkinType := row.Str(5)

		//Formatting the infowindow bubble.
		timestamp := fmt.Sprintf("%s, %s %d, %d at %d:%d:%d", mytime.Weekday().String(), month, day, year, hour, min, sec)
		details := "<p><b>" + checkinType + "</b> <br />" + timestamp + "<br />" + row.Str(6) + "</ p>"

		//Customizing colors in Go. Could do this in javascript, but I don't like javascript at all
		color := "red"
		if checkinType == "OK" {
			color = "blue"
		}

		//Add new GPS location
		myTrip.Coordinates = append(myTrip.Coordinates, Location{row.Float(2), row.Float(3), row.Str(5), details, color})

		//Info for centering and scaling the map
		//Again, I could do this in javascript, but I really don't like javascript
		//longitude
		if longLow > row.Float(2) {
			longLow = row.Float(2)
		}
		if longHigh < row.Float(2) {
			longHigh = row.Float(2)
		}

		//latitude
		if latLow > row.Float(3) {
			latLow = row.Float(3)
		}
		if latHigh < row.Float(3) {
			latHigh = row.Float(3)
		}
	}

	//Zoom the map based on the longest total distance. It errs on the side of showing more map
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

	//Center the map 
	averageLong := (longLow + longHigh) / 2
	averageLat := (latLow + latHigh) / 2

	myTrip.CenterLat = averageLat
	myTrip.CenterLong = averageLong

	return myTrip
}
