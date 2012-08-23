package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

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
				
				param string
				client = &http.Client{
					CheckRedirect: func(req *http.Request, via []*http.Request) errot {
						param = req.URL
					},
				}

				fmt.Println("URL:", param)

				resp, err := http.Get(s)
				
				GetGPSLocationFromId(param)
				fmt.Println("NEW LOCATION:", tweet)
			}
		}

		//Wait 10 seconds
		//TODO: increase this time
		time.Sleep(10000 * time.Millisecond)
	}
}
