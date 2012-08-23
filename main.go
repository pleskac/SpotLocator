package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	for {
		fmt.Println("Pulling tweets... ")
		list, err := GetNewTweets()

		if err != nil {
			continue
		}

		for _, tweet := range *list {
			if strings.HasPrefix(tweet.Text, "Start") {

			} else if strings.HasPrefix(tweet.Text, "End") {

			} else if strings.HasPrefix((tweet.Text), "http://t.co/") {
				fmt.Println("NEW LOCATION:", tweet)
			}
		}

		//SaveTweets(newTweets)
		time.Sleep(10000 * time.Millisecond)
	}
}
