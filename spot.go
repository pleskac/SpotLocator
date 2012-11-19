package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

//In the JSON, Message is an array of objects. 
//In XML, Messages is an array of objects Message.
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
	Message json.RawMessage
	//Message []Message `json:"message"`
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

	jsonBlob := spotResp.Resp.FeedMsgResp.Messages.Message

	//list := make([]Message, 0)
	var list []Message

	if spotResp.Resp.FeedMsgResp.Count == 1 {
		var msg Message
		err := json.Unmarshal(jsonBlob, &msg)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		list := make([]Message, 1)
		list[0] = msg
		//add singleMessage to msg
	} else if spotResp.Resp.FeedMsgResp.Count > 1 {
		err := json.Unmarshal(jsonBlob, &list)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

	}

	return list, nil
}

func GetNewLocations(feedId string, id int) ([]Message, error) {
	allMsgs, err := getMessages(feedId)

	if err != nil {
		return nil, err
	}

	list := make([]Message, 0)

	//FILTER OUT ALREADY FOUND ONESa
	for _, mes := range list {
		fmt.Println(mes.Id)

		if mes.Id > id {
			//add it
			list = append(list, mes)
			fmt.Println("LIST:", list)
		}
	}

	return list, nil
}
