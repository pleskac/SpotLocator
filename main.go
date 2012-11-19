package main

import (
	"fmt"
	"net/http"
	"time"
)

var client *http.Client
var param string
var latestId int

func main() {
	//API endpoint
	go endpoint()

	//update latestId
	latestId = GetLatestId()

	for {
		newLocations, err := GetNewLocations("0oCHzmaKo1zRkSHQglD2qqXkT2yJPvzpK", latestId)

		if err != nil {
			fmt.Println("Error getting new locations:", err)
		}

		for _, location := range newLocations {
			AddGPS(location.Longitude, location.Latitude, location.MessageContent, location.MessageType, location.UnixTime)

			if location.Id > latestId {
				latestId = location.Id
				SaveLatestId(latestId)
			}
		}

		//Wait 30 seconds
		time.Sleep(30000 * time.Millisecond)
	}
}
