package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func tripOutput(w http.ResponseWriter, r *http.Request) {
	output := GetCurrentTrip() //[]Trip{} //{Trip{}, Trip{}}

	enc := json.NewEncoder(w)
	enc.Encode(output)
}

func endpoint() {
	// http://localhost:8080/trip.json

	http.HandleFunc("/trip.json", tripOutput)
	log.Fatal(http.ListenAndServe(":8080", nil)) // this blocks and runs in a loop for you.

}
