package main

import (
	"encoding/json"
	"fmt"
	_ "io/ioutil"
	"net/http"
)

type TweetList []Tweet

type Tweet struct {
	Id string `json:"id_str"`
}

func GetNewTweets() (*TweetList, error) {
	//read last tweet id
	resp, err := http.Get("http://twitter.com/statuses/user_timeline/markpleskac.json?include_entities=true&include_rts=true&trim_user=true")
	if err != nil {
		fmt.Println("ERROR!")
		return nil, err
	}

	defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)	

	dec := json.NewDecoder(resp.Body)
	list := new(TweetList)

	if err = dec.Decode(list); err != nil {
		return nil, err
	}
	return list, nil
}
