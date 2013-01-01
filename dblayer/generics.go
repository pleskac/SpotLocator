package dblayer

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	z_mysql "github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"net/http"
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

type TimeZoneResponse struct {
	DstOffset    float64 `json:"dstOffset"`
	RawOffset    float64 `json:"rawOffset"`
	Status       string  `json:"status"`
	TimeZoneId   string  `json:"timeZoneId"`
	TimeZoneName string  `json:"timeZoneName"`
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

func getTimeZoneTime(long, lat float64, utcTime int64) (int64, string) {
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/timezone/json?location=%f,%f&timestamp=%d&sensor=false", lat, long, utcTime)

	//https - skip the verification for ease of use
	//shouldn't really be any reason this needs to be secure
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return utcTime, ""
	}

	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	tzResp := &TimeZoneResponse{}
	if err := dec.Decode(tzResp); err != nil {
		fmt.Println(err)
		return utcTime, ""
	}

	timeZone := tzResp.TimeZoneId + " (" + tzResp.TimeZoneName + ")"

	return utcTime + int64(tzResp.RawOffset+tzResp.DstOffset), timeZone
}
