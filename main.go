package main

import (
	"fmt"
	"net/http"
	"time"
)

var client *http.Client
var param string
var latestId int

func main() {
	go endpoint()

	//update latestId

	for {
		_, err := GetNewLocations("0oCHzmaKo1zRkSHQglD2qqXkT2yJPvzpK", latestId)

		if err != nil {
			fmt.Println(err)
		}

		//Wait 100 seconds or so
		time.Sleep(1000000 * time.Millisecond)
	}
}
