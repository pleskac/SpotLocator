package main

import (
	"fmt"
	"time"
)

func main() {
	for{
		fmt.Println("Pulling tweets... ")
		list, err := GetNewTweets()

		if err != nil{
			continue
		}
		
		for _, tweet := range *list{
			fmt.Println(tweet)
		}
		

		//SaveTweets(newTweets)
		time.Sleep(10000  * time.Millisecond)
	}
}
