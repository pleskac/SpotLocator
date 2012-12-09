package mysql

import (
	"fmt"
	z_mysql "github.com/ziutek/mymysql/mysql"
	//This reference is necessary, otherwise a panic will occur upon calling mysql.New()
	_ "github.com/ziutek/mymysql/native"
	"strconv"
	"time"
)

//internal constants
const (
	lATEST_SPOT = "SPOT"
	pASSWORD    = "PASSWORD"
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
	TripId      int
	TripName    string
	IsCurrent   int
	Coordinates []Location
}

func Connect() z_mysql.Conn {
	//Set up database connection
	db := z_mysql.New("tcp", "", "127.0.0.1:3306", "root", "rootroot", "gps")
	err := db.Connect()
	if err != nil {
		fmt.Println("ERROR CONNECTING:", err)
		panic(err)
	}

	return db
}

func GetLatestSpotId() int {
	spotStr := getValue(lATEST_SPOT)
	spotId, err := strconv.Atoi(spotStr)

	if err != nil {
		fmt.Println(err)
		return -1
	}

	return spotId
}

func getValue(idType string) string {
	db := Connect()
	defer db.Close()

	rows, _, err := db.Query("SELECT v FROM kvp WHERE k = " + idType)
	if err != nil {
		panic(err)
	}

	if len(rows) < 1 {
		return "-1"
	}

	return rows[0].Str(0)
}

func SaveLatestSpotId(id int) {
	idStr := fmt.Sprintf("%d", id)
	saveValue(lATEST_SPOT, idStr)
}

func saveValue(key string, newValue string) {
	db := Connect()
	defer db.Close()

	//Delete that row, if it exists
	stmt, err := db.Prepare("DELETE FROM kvp WHERE k = " + key)
	_, err = stmt.Run()
	if err != nil {
		//row did not exist, we'll just add it later
	}

	//Insert new row
	stmt, err = db.Prepare("INSERT INTO kvp (k, v) VALUES (?, ?)")
	_, err = stmt.Run(key, newValue)
	if err != nil {
		panic(err)
	}
}

func AddGPS(longitude, latitude float64, message, msgType string, time int64) {
	db := Connect()
	defer db.Close()

	//Get the current trip, if it exists
	rows, _, err := db.Query("SELECT id FROM trips WHERE current = 1")
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
//Only allow this to access MySQL, move formatting data somewhere else
//Also, could default to the current trip, or allow specific trips to be returned, that would allow multiple maps on the site
func GetCurrentTripId() int {
	db := Connect()
	defer db.Close()

	//Get the current trip, if it exists
	rows, _, err := db.Query("SELECT * FROM trips WHERE current = 1")
	if err != nil {
		panic(err)
	}

	if len(rows) > 1 {
		fmt.Println("More than one row!! WRONG!")
		return -1
	} else if len(rows) == 0 {
		fmt.Println("0 rows! No current trip to return")
		return -1
	}
	id := rows[0].Int(0) //the first(only) row. the first element is the id.

	return id
}

func FindTrip(name string) int {
	db := Connect()
	defer db.Close()

	query := "SELECT * FROM trips WHERE name LIKE '%" + name + "%'"

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

func GetTrip(id int) Trip {
	if id < 0 {
		//ids cannot be negative
		return Trip{}
	}

	db := Connect()
	defer db.Close()

	tripQuery := fmt.Sprintf("SELECT * FROM trips WHERE id = %d", id)
	gpsQuery := fmt.Sprintf("SELECT * FROM gps WHERE trip = %d", id)

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
		timestamp := fmt.Sprintf("%s, %s %d, %d at %d:%d:%d", mytime.Weekday().String(), month, day, year, hour, min, sec)
		details := "<p><b>" + checkinType + "</b> <br />" + timestamp + "<br />" + row.Str(6) + "</ p>"

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

func GetTripList() []Trip {
	var list []Trip
	query := "SELECT * FROM trips ORDER BY id DESC"

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
