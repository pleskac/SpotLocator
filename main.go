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
	go endpoint()

	//update latestId
	latestId = GetLatestId()

	for {
		newLocations, err := GetNewLocations("0oCHzmaKo1zRkSHQglD2qqXkT2yJPvzpK", latestId)

		if err != nil {
			fmt.Println("Error getting new locations:", err)
		}

		//save the latest tweet

		for _, location := range newLocations {
			//TODO: uncomment
			fmt.Println("New location!")
			AddGPS(location.Longitude, location.Latitude, location.MessageContent, location.UnixTime)

			if location.Id > latestId {
				latestId = location.Id
				SaveLatestId(latestId)
			}
		}

		//Wait 100 seconds or so
		time.Sleep(1000000 * time.Millisecond)
	}
}
