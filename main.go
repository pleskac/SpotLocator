package main

import (
	"fmt"
	"time"
)

func main() {
	for{
		fmt.Println("Pulling tweets... NOT!")
		GetNewTweets() //save as something
		//SaveTweets(newTweets)
		time.Sleep(10000  * time.Millisecond)
	}
}
