package main

import (
	"./dblayer"
	"fmt"
	"net/http"
	"strconv"
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
			dblayer.AddGPS_UTC(location.Longitude, location.Latitude, location.MessageContent, location.MessageType, "markpleskac@gmail.com", location.UnixTime)

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

		devices := dblayer.GetDevices(user)
		for _, device := range devices {
			//this is just one SPOT or .... something else! 
			//make it generic!
			//get a DEVICE struct
			if device.Type == "SPOT" {
				fmt.Println(device)
				myint, conversionErr := strconv.ParseInt(device.MostRecent, 10, 64)
				if conversionErr != nil {
					newLocations, err := GetNewLocations(device.Key, int(myint))
					if err != nil {
						fmt.Println("Error getting new locations:", err)
					}

					for _, location := range newLocations {
						fmt.Println("Adding new GPS location", location.MessageType)

						//SPOT returns time in UTC. This will correct the time to the localized time.
						//TODO: THIS NEEDS TO INCORPORATE A USER
						dblayer.AddGPS_UTC(location.Longitude, location.Latitude, location.MessageContent, location.MessageType, user, location.UnixTime)

						if location.Id > latestSpotId {
							latestSpotId = location.Id
							dblayer.SaveLatestSpotId(latestSpotId)
						}
					}

				}

			}

		}

	}
}
