package main

import (
	"net/http"
	"strings"
	"time"
)

var client *http.Client
var param string

func main() {
	for {

		//Get all new tweets
		list, err := GetNewTweets()

		if err != nil {
			continue
		}

		//sort the tweets

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

				AddGPS(long, lat)
			}
		}

		//Wait 10 seconds
		//TODO: increase this time
		time.Sleep(10000 * time.Millisecond)
	}
}
