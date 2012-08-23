package main

import (
	"encoding/json"
	"fmt"
	_ "io/ioutil"
	"net/http"
)

type TweetList []Tweet

type Tweet struct {
	Id_str     string
	Text       string
	Created_at string
}

func GetNewTweets(lastTweet string) (*TweetList, error) {
	//read last tweet id
	var resp *http.Response
	var err error
	if lastTweet == "" {
		resp, err = http.Get("http://twitter.com/statuses/user_timeline/markpleskac.json?include_entities=true&include_rts=true&trim_user=true")
	} else {
		resp, err = http.Get("http://twitter.com/statuses/user_timeline/markpleskac.json?include_entities=true&include_rts=true&trim_user=true&since_id=" + lastTweet)
	}

	if err != nil {
		//fmt.Println("Error pulling new tweets:", err)
		return nil, err
	}

	//Closes the http response at the end of the function
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	list := new(TweetList)

	//Decodes the JSON, if there is an error return it
	if err = dec.Decode(list); err != nil {
		fmt.Println("Error decoding:", err)
		return nil, err
	}
	return list, nil
}
