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
	latestId = GetLatestTweet()

	for {

		//Get all new tweets
		list, err := GetNewTweets(latestId)

		if list != nil && len(*list) > 0 {
			latestId = (*list)[0].Id_str
			SaveLatestTweet(latestId) //This could be in a new thread... only done for restarts

		}
		fmt.Println("Latest Tweet:", latestId)

		if err != nil {
			continue
		}

		//put the new tweets in the DB
		//iterate through the tweets backwards so they are oldest to newest
		for n := len(*list) - 1; n >= 0; n-- {
			tweet := (*list)[n]
			fmt.Println("Looking at tweet", tweet.Id_str, tweet.Text)
			//|| strings.HasPrefix(tweet.Text, "Create") || strings.HasPrefix(tweet.Text, "create")
			if strings.HasPrefix(tweet.Text, "Start") || strings.HasPrefix(tweet.Text, "start") {
				//Create a new trip
				CreateTrip(tweet.Text[6:len(tweet.Text)])

			} else if strings.HasPrefix(tweet.Text, "End") || strings.HasPrefix(tweet.Text, "end") {
				//End the current trip
				EndTrips()

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

				long, lat, msg, err := GetGPSLocationFromId(param)

				if err != nil {
					continue
				}

				//NEED TO IMPLEMENT THIS
				AddGPS(long, lat, msg)
			}
		}

		//Wait 10 seconds
		//TODO: increase this time
		time.Sleep(10000 * time.Millisecond)
	}
}
