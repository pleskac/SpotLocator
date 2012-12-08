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
	r := mux.NewRouter()
	r.HandleFunc("/trips/{"+tripId+"}", TripHandler)
	r.HandleFunc("/currentTrip", tripOutput)
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
