package main

import (
	"./mysql"
	"fmt"
	"net/http"
	"time"
)

var client *http.Client
var param string
var latestSpotId int

func main() {
	//API endpoint
	go endpoint()

	//update latestSpotId
	latestSpotId = mysql.GetLatestSpotId()

	for {
		newLocations, err := GetNewLocations("0oCHzmaKo1zRkSHQglD2qqXkT2yJPvzpK", latestSpotId)

		if err != nil {
			fmt.Println("Error getting new locations:", err)
		}

		for _, location := range newLocations {
			mysql.AddGPS(location.Longitude, location.Latitude, location.MessageContent, location.MessageType, location.UnixTime)

			if location.Id > latestSpotId {
				latestSpotId = location.Id
				mysql.SaveLatestSpotId(latestSpotId)
			}
		}

		//Wait 30 seconds
		time.Sleep(30000 * time.Millisecond)
	}
}
