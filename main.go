package main

import (
	"./dblayer"
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
	latestSpotId = dblayer.GetLatestSpotId()

	for {
		newLocations, err := GetNewLocations("0oCHzmaKo1zRkSHQglD2qqXkT2yJPvzpK", latestSpotId)

		if err != nil {
			fmt.Println("Error getting new locations:", err)
		}

		for _, location := range newLocations {
			fmt.Println("Adding new GPS location", location.MessageType)

			dblayer.AddGPS(location.Longitude, location.Latitude, location.MessageContent, location.MessageType, location.UnixTime)

			if location.Id > latestSpotId {
				latestSpotId = location.Id
				dblayer.SaveLatestSpotId(latestSpotId)
			}
		}

		//Wait 30 seconds
		time.Sleep(30000 * time.Millisecond)
	}
}
