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

	//LOL this is so bad
	go newMain()

	//OLD CODE TO KEEP THIS BEAST RUNNING
	//update latestSpotId
	latestSpotId = dblayer.GetLatestSpotId()

	for {
		//DO ALL SPOT UPDATES
		newLocations, err := GetNewLocations("0oCHzmaKo1zRkSHQglD2qqXkT2yJPvzpK", latestSpotId)

		if err != nil {
			fmt.Println("Error getting new locations:", err)
		}

		for _, location := range newLocations {
			fmt.Println("Adding new GPS location", location.MessageType)

			//SPOT returns time in UTC. This will correct the time to the localized time.
			dblayer.AddGPS_UTC(location.Longitude, location.Latitude, location.MessageContent, location.MessageType, location.UnixTime)

			if location.Id > latestSpotId {
				latestSpotId = location.Id
				dblayer.SaveLatestSpotId(latestSpotId)
			}
		}

		//DO OTHER UPDATES FROM OTHER DEVICES 

		//Wait 30 seconds
		time.Sleep(30000 * time.Millisecond)
	}
}

func newMain() {
	//I don't have testing, which I could fix.
	//I'm also testing in production. 
	//No-nos that I'm too lazy to change right now.

	//get users
	users := dblayer.GetAllUsers()

	//for every user, update every device
	for _, user := range users {
		fmt.Println(user)
	}
}
