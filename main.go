package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var client *http.Client
var param string
var latestId string

func main() {
	latestId = ""
	for {

		//Get all new tweets
		list, err := GetNewTweets(latestId)

		if list != nil && len(*list) > 0 {
			latestId = (*list)[0].Id_str
		}
		fmt.Println("Latest Tweet:", latestId)
		//TODO: make sure the latestId is correct
		/////////////////////////////////////////////////////
		/////////////////////////////////////////////////////
		/////////////////////////////////////////////////////

		if err != nil {
			continue
		}

		//TODO: sort the tweets
		/////////////////////////////////////////////////////
		/////////////////////////////////////////////////////
		/////////////////////////////////////////////////////

		//put the new tweets in the DB
		for _, tweet := range *list {
			if strings.HasPrefix(tweet.Text, "Start") {
				//Create a new trip

			} else if strings.HasPrefix(tweet.Text, "End") {
				//End the current trip

			} else if strings.HasPrefix((tweet.Text), "http://t.co/") {
				//Add location to current trip
				//If no current trip, GPS data is stored with no pointer to a trip

				s := (tweet.Text)[0:20]
				firstRedirect := true
				client = &http.Client{
					CheckRedirect: func(req *http.Request, via []*http.Request) error {
						if firstRedirect {
							param = req.URL.Path
							firstRedirect = false
						}
						//fmt.Println("REDIRECTED! param:", param)
						//return errors.New("don't follow")
						return nil
					},
				}
				_, err := client.Get(s)

				if err != nil {
					continue
				}

				long, lat, err := GetGPSLocationFromId(param)

				if err != nil {
					continue
				}

				//NEED TO IMPLEMENT THIS
				AddGPS(long, lat)
			}
		}

		//Wait 10 seconds
		//TODO: increase this time
		time.Sleep(10000 * time.Millisecond)
	}
}
