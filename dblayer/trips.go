package dblayer

import (
	"fmt"
	_ "github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"time"
)

func AddGPSNow(longitude, latitude float64, message, msgType, user string) {
	//NOW is in UTC
	now := time.Now().Unix()
	AddGPS_UTC(longitude, latitude, message, msgType, user, now)
}

func AddGPS_UTC(longitude, latitude float64, message, msgType, user string, utcTime int64) {
	fixedTime, timeZone := getTimeZoneTime(longitude, latitude, utcTime)
	addGPS(longitude, latitude, message, msgType, timeZone, user, fixedTime)
}

func addGPS(longitude, latitude float64, message, msgType, timeZone, user string, time int64) {
	db := Connect()
	defer db.Close()

	fmt.Println(user)

	//Get the current trip, if it exists
	query := "SELECT id FROM trips WHERE user = '" + user + "' AND current = 1"
	rows, _, err := db.Query(query)
	if err != nil {
		panic(err)
	}

	//Get the foreign key to the current trip
	tripKey := -1
	if len(rows) != 1 {
		tripKey = -1
	} else {
		tripKey = (rows[0]).Int(0)
	}

	//Add the GPS row with data
	if tripKey == -1 {
		stmt, err := db.Prepare("INSERT INTO gps (longitude, latitude, details, timestamp, type, timezone) VALUES (?, ?, ?, ?, ?, ?)")
		_, err = stmt.Run(longitude, latitude, message, time, msgType, timeZone)
		if err != nil {
			panic(err)
		}
	} else {
		stmt, err := db.Prepare("INSERT INTO gps (longitude, latitude, details, trip, timestamp, type, timezone) VALUES (?, ?, ?, ?, ?, ?, ?)")
		_, err = stmt.Run(longitude, latitude, message, tripKey, time, msgType, timeZone)
		if err != nil {
			panic(err)
		}
	}

}

func CreateTrip(name, details, user string) {
	db := Connect()
	defer db.Close()

	//End all trips
	EndTrips(user)

	//Create new trip, set it as current trip
	fmt.Println("Starting trip", name)

	//Insert new trip
	stmt, err := db.Prepare("INSERT INTO trips (name, details, user, current) VALUES (?, ?, ?, ?)")
	_, err = stmt.Run(name, details, user, 1)
	if err != nil {
		panic(err)
	}
}

func EndTrips(user string) {
	db := Connect()
	defer db.Close()

	fmt.Println("Ending all trips")

	query := "SELECT id FROM trips WHERE current = 1 AND user = '" + user + "'"
	rows, _, err := db.Query(query)
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

//TODO: user...
//TODO: this can be split up and organized better
//Only allow this to access MySQL, move formatting data somewhere else
//Also, could default to the current trip, or allow specific trips to be returned, that would allow multiple maps on the site
func GetCurrentTripId(user string) int {
	db := Connect()
	defer db.Close()

	//Get the current trip, if it exists
	currentTripQuery := "SELECT * FROM trips WHERE current = 1 AND user = '" + user + "'"
	rows, _, err := db.Query(currentTripQuery)
	if err != nil {
		panic(err)
	}

	if len(rows) != 1 {
		return -1
	}

	id := rows[0].Int(0) //the first(only) row. the first element is the id.

	return id
}

func FindTrip(name, user string) int {
	db := Connect()
	defer db.Close()

	query := "SELECT * FROM trips WHERE user = '" + user + "' AND name LIKE '%" + name + "%'"

	rows, _, err := db.Query(query)
	if err != nil {
		panic(err)
	}

	if len(rows) == 0 {
		return -1
	}

	//return first match
	return rows[0].Int(0)
}

//Do I need the user??... only for authentication...
func GetTrip(id int) Trip {
	if id < 0 {
		//ids cannot be negative
		return Trip{}
	}

	db := Connect()
	defer db.Close()

	tripQuery := fmt.Sprintf("SELECT * FROM trips WHERE id = %d", id)
	gpsQuery := fmt.Sprintf("SELECT * FROM gps WHERE trip = %d ORDER BY timestamp", id)

	rows, _, err := db.Query(tripQuery)
	if err != nil {
		panic(err)
	}

	if len(rows) < 1 {
		//Trip with that id does not exist
		return Trip{}
	}

	name := (rows[0]).Str(1)
	isCurrent := (rows[0]).Int(3)
	myTrip := Trip{id, name, isCurrent, nil}

	rows, _, err = db.Query(gpsQuery)
	if err != nil {
		panic(err)
	}

	//Add every GPS location
	for _, row := range rows {
		//Timestamp: Details
		mytime := time.Unix(row.Int64(4), 0)
		year, month, day := mytime.Date()
		hour, min, sec := mytime.Clock()

		checkinType := row.Str(5)

		//Formatting the infowindow bubble.
		timestamp := fmt.Sprintf("%s, %s %d, %d at %d:%02d:%02d<br />%s", mytime.Weekday().String(), month, day, year, hour, min, sec, row.Str(7))

		msgDetails := row.Str(6)
		if msgDetails == "" {
			msgDetails = "No details available"
		}

		details := "<p><strong>" + checkinType + "</strong>  (" + row.Str(0) + ") <br />" +
			timestamp + "<br /><br />" +
			"Latitude: " + row.Str(2) + "<br />" +
			"Longitude: " + row.Str(3) + "<br /><br />" +
			msgDetails + "</ p>"
		//TODO: add related pictures

		//Customizing colors in Go. Could do this in javascript, but I don't like javascript at all
		color := "Red"
		if checkinType == "OK" {
			color = "RoyalBlue"
		} else if checkinType == "TRACK" {
			color = "DarkOliveGreen"
		} else if checkinType == "CUSTOM" {
			color = "Orange"
		}

		//Add new GPS location
		myTrip.Coordinates = append(myTrip.Coordinates, Location{row.Float(2), row.Float(3), row.Str(5), details, color})
	}

	return myTrip
}

func GetTripList(user string) []Trip {
	var list []Trip
	query := "SELECT * FROM trips WHERE user = '" + user + "' ORDER BY id DESC"

	db := Connect()
	defer db.Close()

	rows, _, err := db.Query(query)
	if err != nil {
		panic(err)
	}

	for _, row := range rows {
		//Add every id of every trip to the list
		list = append(list, Trip{row.Int(0), row.Str(1), row.Int(3), nil})
	}

	return list
}
