package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

//In the JSON, Message is an array of object(s). 
//In XML, Messages is an array of object(s) Message.
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
	Messages      Messages `json:"messages"`
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
	//This is generic to allow for both a list of 
	//messages and a single message. The SPOT API's
	//naming convention for JSON is unusual and 
	//doesn't make sense. The XML makes more sense.
	Message json.RawMessage
}

type Message struct {
	AtClientUnixTime string  `json:"@clientUnixTime"`
	Id               int     `json:"id"`
	MessengerId      string  `json:"messengerId"`
	MessengerName    string  `json:"messengerName"`
	UnixTime         int64   `json:"unixTime"`
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

func getMessages(feedId string) ([]Message, error) {
	url := "https://api.findmespot.com/spot-main-web/consumer/rest-api/2.0/public/feed/" + feedId + "/message.json"

	//https - skip the verification for ease of use
	//shouldn't really be any reason this needs to be secure
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

	jsonBlob := spotResp.Resp.FeedMsgResp.Messages.Message

	fmt.Println("count from json", spotResp.Resp.FeedMsgResp.Count)

	if spotResp.Resp.FeedMsgResp.Count == 1 {
		//Single message in JSON
		var msg Message
		err := json.Unmarshal(jsonBlob, &msg)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		list := make([]Message, 0)
		list = append(list, msg)
		fmt.Println("single message returned!", len(list))
		return list, nil
	} else if spotResp.Resp.FeedMsgResp.Count > 1 {
		var list []Message
		//Multiple messages in JSON
		err := json.Unmarshal(jsonBlob, &list)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return list, nil
	}

	return nil, nil
}

func GetNewLocations(feedId string, id int) ([]Message, error) {
	allMsgs, err := getMessages(feedId)
	fmt.Println("allMsgs count", len(allMsgs))
	if err != nil {
		return nil, err
	}

	list := make([]Message, 0)

	//Only return messages that are "new"
	//i.e. messages with Id greater than the previous "new" message
	for _, mes := range allMsgs {
		if mes.Id > id {
			//New message!
			list = append(list, mes)
		}
	}

	fmt.Println(len(list))

	return list, nil
}
