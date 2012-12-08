package main

import (
	"code.google.com/p/gorilla/mux"
	"encoding/json"
	"fmt"
	_ "log"
	"net/http"
)

const tripId = "tripId"
const tripName = "tripName"

func endpoint() {
	router := mux.NewRouter()
	r := router.Host("{domain:pleskac.org|api.pleskac.org|localhost}").Subrouter()

	r.HandleFunc("/api/trip/id/{"+tripId+"}", TripIdHandler)
	r.HandleFunc("/api/trip/name/{"+tripName+"}", TripNameHandler)
	r.HandleFunc("/api/currentTrip", CurrentTripHandler)

	http.ListenAndServe(":8080", r)
}

func CurrentTripHandler(w http.ResponseWriter, r *http.Request) {
	output := GetCurrentTrip()

	enc := json.NewEncoder(w)
	enc.Encode(output)
}

func TripIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars[tripId]

	output := GetTrip(id)

	enc := json.NewEncoder(w)
	enc.Encode(output)

	fmt.Println("Id:", id)
}

func TripNameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars[tripName]

	fmt.Println("Name:", name)
}
