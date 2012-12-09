package dblayer

import (
	"fmt"
	z_mysql "github.com/ziutek/mymysql/mysql"
	//This reference is necessary, otherwise a panic will occur upon calling mysql.New()
	_ "github.com/ziutek/mymysql/native"
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
