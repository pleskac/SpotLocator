package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

type SpotApiResponse struct {
	Resp Response `json:"response"`
}

type Response struct {
	FeedMsgResp FeedMessageResponse `json:"feedMessageResponse"`
}

type FeedMessageResponse struct {
	Count         int      `json:"count"`
	Feed          Feed     `json:"feed"`
	TotalCount    int      `json:"totalCount"`
	ActivityCount int      `json:"activityCount"`
	Messages      Messages `json:"messages"` // this ought to be an array but it's not from Spot? this might break...
}

type Feed struct {
	Id                   string `json:"id"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	Status               string `json:"status"`
	Usage                int    `json:"usage"`
	DaysRange            int    `json:"daysRange"`
	DetailedMessageShown bool   `json:"detailedMessageShown"`
}

type Messages struct {
	Message []Message `json:"message"`
}

type Message struct {
	AtClientUnixTime string  `json:"@clientUnixTime"`
	Id               int     `json:"id"`
	MessengerId      string  `json:"messengerId"`
	MessengerName    string  `json:"messengerName"`
	UnixTime         int     `json:"unixTime"`
	MessageType      string  `json:"messageType"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	ShowCustomMsg    string  `json:"showCustomMsg"`
	DateTime         string  `json:"dateTime"`
	MessageDetail    string  `json:"messageDetail"`
	Selected         bool    `json:"selected"`
	Altitude         int     `json:"altitude"`
	Hidden           int     `json:"hidden"`
	MessageContent   string  `json:"messageContent"`
}

type NewLocations []Message

func getJsonObject(feedId string) (*SpotApiResponse, error) {
	url := "https://api.findmespot.com/spot-main-web/consumer/rest-api/2.0/public/feed/" + feedId + "/message.json"

	fmt.Println("URL:", url)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	spotResp := &SpotApiResponse{}
	if err := dec.Decode(spotResp); err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(spotResp)

	return spotResp, nil
}

func GetNewLocations(feedId string, id int) (*NewLocations, error) {
	json, err := getJsonObject(feedId)

	if err != nil {
		return nil, err
	}
	//FILTER

	for _, mes := range json.Resp.FeedMsgResp.Messages.Message {
		fmt.Println(mes.Id)
	}

	/*for mes := range jsonStructure.Response.FeedMessageResponse.Messages.Message {
		fmt.Println(mes)
	}
	*/
	//jsonStructure.Response.MessagesResponse.Messages
	return nil, nil
}

/*
func GetGPSLocationFromId(id string) (float32, float32, string, int64, error) {
	url := "http://share.findmespot.com/spot-adventures/rest-api/1.0/public/location" + id
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error pulling GPS location:", err)
		return 0.0, 0.0, "", 0, err
	}

	//Closes the http response at the end of the function
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	res := new(TopStructure)

	if err = dec.Decode(res); err != nil {
		fmt.Println("Error decoding:", err)
		return 0.0, 0.0, "", 0, err
	}
	//fmt.Println("RESPOSE: ", res)
	//fmt.Println("Latitude:", res.Response.MessagesResponse.Messages.Message.Latitude)
	long := res.Response.MessagesResponse.Messages.Message.Longitude
	lat := res.Response.MessagesResponse.Messages.Message.Latitude
	msg := res.Response.MessagesResponse.Messages.Message.MessageContent
	tm := res.Response.MessagesResponse.Messages.Message.Time

	fmt.Println("Timestamp:", tm)
	return long, lat, msg, tm, nil
}
*/
