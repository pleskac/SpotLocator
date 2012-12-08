package main

import (
	"code.google.com/p/gorilla/mux"
	"encoding/json"
	"fmt"
	_ "log"
	"net/http"
)

const tripId = "tripId"

func tripOutput(w http.ResponseWriter, r *http.Request) {
	output := GetCurrentTrip()

	enc := json.NewEncoder(w)
	enc.Encode(output)
}

func endpoint() {
	/*
		//http://localhost:8080/trip.json
		http.HandleFunc("/trip.json", tripOutput)

		//Blocks and runs in a loop
		log.Fatal(http.ListenAndServe(":8080", nil))
	*/
	router := mux.NewRouter()
	r := router.Host("{domain:pleskac.org:8080|localhost}").Subrouter()
	r.HandleFunc("/api/trips/{"+tripId+"}", TripHandler)
	r.HandleFunc("/api/currentTrip", tripOutput)
	//legacy. remove later.
	r.HandleFunc("/trip.json", tripOutput)
	//can add other JSON handlers here
	http.Handle("/", r)
}

func TripHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars[tripId]
	fmt.Println("Id:", id)
}
