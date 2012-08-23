package main

import (
	"fmt"
	_ "io/ioutil"
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
							param = req.URL.Path[1:len(param)]
							firstRedirect = false
						}
						fmt.Println("REDIRECTED! param:", param)
						return nil
					},
				}
				fmt.Println("Twitter URL:", s)
				_, err := client.Get(s)

				//body, err := ioutil.ReadAll(resp.Body)
				//fmt.Println(string(body))
				if err != nil {
					continue
				}

				fmt.Println("Spot URL:", param)

				GetGPSLocationFromId(param)
				fmt.Println("NEW LOCATION:", tweet)
			}
		}

		//Wait 10 seconds
		//TODO: increase this time
		time.Sleep(10000 * time.Millisecond)
	}
}
