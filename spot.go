package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	MessagesResponse MessagesResponse
}

type MessagesResponse struct {
	Messages Messages
}

type Messages struct {
	latitude  float32
	longitude float32
}

func GetGPSLocationFromId(id string) (float32, float32, error) {
	url := "http://share.findmespot.com/spot-adventures/rest-api/1.0/public/location/" + id
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error pulling GPS location:", err)
		return 0.0, 0.0, err
	}

	//Closes the http response at the end of the function
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	res := new(Response)

	if err = dec.Decode(res); err != nil {
		fmt.Println("Error decoding:", err)
		return 0.0, 0.0, err
	}
	fmt.Println("RESPOSE: ", res)
	fmt.Println("Latitude:", res.MessagesResponse.Messages.latitude)

	return 0.0, 0.0, nil
}
