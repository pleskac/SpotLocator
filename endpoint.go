package main

import (
	"code.google.com/p/gorilla/mux"
	"encoding/json"
	"fmt"
	_ "log"
	"net/http"
	"strconv"
)

const tripId = "tripId"
const tripName = "tripName"

func endpoint() {
	router := mux.NewRouter()
	r := router.Host("{domain:pleskac.org|api.pleskac.org|localhost}").Subrouter()

	r.HandleFunc("/api/trip/id/{"+tripId+":[0-9]+}", TripIdHandler)
	r.HandleFunc("/api/trip/name/{"+tripName+"}", TripNameHandler)
	r.HandleFunc("/api/trip/currentTrip", CurrentTripHandler)

	http.ListenAndServe(":8080", r)
}

func CurrentTripHandler(w http.ResponseWriter, r *http.Request) {
	currentTripId := GetCurrentTripId()
	output := GetTrip(currentTripId)

	fmt.Println(output)

	enc := json.NewEncoder(w)
	enc.Encode(output)
}

func TripIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars[tripId]

	//convert the string id in the URI to an int
	id, err := strconv.Atoi(idStr)
	if id < 0 || err != nil {
		fmt.Println("Error parsing", id, "\n", err)
		return
	}

	output := GetTrip(id)

	enc := json.NewEncoder(w)
	enc.Encode(output)
}

func TripNameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars[tripName]

	fmt.Println("Name:", name)

}
