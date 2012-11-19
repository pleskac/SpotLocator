package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func tripOutput(w http.ResponseWriter, r *http.Request) {
	output := GetCurrentTrip()

	enc := json.NewEncoder(w)
	enc.Encode(output)
}

func endpoint() {
	//http://localhost:8080/trip.json
	http.HandleFunc("/trip.json", tripOutput)

	//Blocks and runs in a loop
	log.Fatal(http.ListenAndServe(":8080", nil))
}
