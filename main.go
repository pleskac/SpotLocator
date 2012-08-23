package main

import (
	"fmt"
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
				s := (tweet.Text)[11:16]
				fmt.Println(s)
				GetGPSLocationFromId(s)
				fmt.Println("NEW LOCATION:", tweet)
			}
		}

		//Wait 10 seconds
		//TODO: increase this time
		time.Sleep(10000 * time.Millisecond)
	}
}
