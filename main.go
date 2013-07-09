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
	//I don't have testing, which I could fix.
	//I'm also testing in production. 
	//No-nos that I'm too lazy to change right now.

	//Get users
	users := dblayer.GetAllUsers()

	//For every user, update every device
	for _, user := range users {
		fmt.Println(user)

		devices := dblayer.GetDevices(user)
		//Iterate through all of the users' devices
		for _, device := range devices {
			//Device is a SPOT GPS. First supported device.
			if device.Type == "SPOT" {
				//Get the latest spot id
				lastSpotId, conversionErr := strconv.ParseInt(device.MostRecent, 10, 64)
				if conversionErr != nil {
					newLocations, err := GetNewLocations(device.Key, int(lastSpotId))
					if err != nil {
						fmt.Println("Error getting new locations:", err)
					}

					for _, location := range newLocations {
						//fmt.Println("Adding new GPS location", location.MessageType)

						//SPOT returns time in UTC. This will correct the time to the localized time.
						dblayer.AddGPS_UTC(location.Longitude, location.Latitude, location.MessageContent, location.MessageType, user, location.UnixTime)

						if location.Id > latestSpotId {
							latestSpotId = location.Id
							dblayer.SaveLatestSpotId(latestSpotId)
						}
					}
				}
			} //ELSE IF THIS IS ANOTHER TYPE OF DEVICE, ADD SUPPORT HERE
		}

		time.Sleep(30000 * time.Millisecond)

	} //End of the infinite 'for', wait between users. Bad because the more users, the less refresh rate.

	/*
		//OLD CODE TO KEEP THIS BEAST RUNNING, KEPT FOR SAFE KEEPING
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
		}*/
}
