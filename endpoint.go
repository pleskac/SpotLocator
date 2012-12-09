package main

import (
	//change to github.com/pleskac/SpotLocator/dblayer
	"./dblayer"
	"code.google.com/p/gorilla/mux"
	"encoding/json"
	"fmt"
	_ "log"
	"net/http"
	"strconv"
)

const tripId = "tripId"
const tripName = "tripName"
const longitude = "longitude"
const latitude = "latitude"
const gpsType = "gpsType"
const gpsMessage = "gpsMessage"

//JSON endpoints:
//	/api/trip/id/{ID}		looks up by trip id
//	/api/trip/name/{NAME}	searches tips by name and returns the first trip matching that string
//	/api/trip/currentTrip	returns the current trip
//	/api/trip/list			returns a list of all trips
func endpoint() {
	password := dblayer.GetPassword()

	fmt.Println("PASSWORD:", password)

	router := mux.NewRouter()
	r := router.Host("{domain:pleskac.org|api.pleskac.org|localhost}").Subrouter()

	r.HandleFunc("/api/trip/id/{"+tripId+":[0-9]+}", TripIdHandler)
	r.HandleFunc("/api/trip/name/{"+tripName+"}", TripNameHandler)
	r.HandleFunc("/api/trip/currentTrip", CurrentTripHandler)
	r.HandleFunc("/api/trip/list", TripListHandler)
	r.HandleFunc("/api/trip/add/"+password+"/{"+tripName+"}", AddTripHandler)
	r.HandleFunc("api/gps/add/"+password+"/{"+longitude+"/{"+latitude+"}/{"+gpsType+"}/{"+gpsMessage+"}", AddGPSHandler)

	http.ListenAndServe(":8080", r)
}

func AddGPSHandler(w http.ResponseWriter, r *http.Request) {
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "http://pleskac.org")

	vars := mux.Vars(r)
	longStr := vars[longitude]
	latStr := vars[latitude]

	fmt.Println("Adding GPS location via webservice")

	longFlt, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		fmt.Println("Error parsing", longStr, "\n", err)
		return
	}

	latFlt, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		fmt.Println("Error parsing", latStr, "\n", err)
		return
	}

	dblayer.AddGPSNow(longFlt, latFlt, vars[gpsMessage], vars[gpsType])

	enc := json.NewEncoder(w)
	enc.Encode(vars[gpsType])
}

func AddTripHandler(w http.ResponseWriter, r *http.Request) {
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "http://pleskac.org")

	vars := mux.Vars(r)
	name := vars[tripName]

	if name != "" {
		dblayer.CreateTrip(name)
	}

	enc := json.NewEncoder(w)
	enc.Encode(name)
}

func TripListHandler(w http.ResponseWriter, r *http.Request) {
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "http://pleskac.org")

	output := dblayer.GetTripList()

	enc := json.NewEncoder(w)
	enc.Encode(output)
}

func CurrentTripHandler(w http.ResponseWriter, r *http.Request) {
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "http://pleskac.org")

	currentTripId := dblayer.GetCurrentTripId()
	output := dblayer.GetTrip(currentTripId)

	enc := json.NewEncoder(w)
	enc.Encode(output)
}

func TripIdHandler(w http.ResponseWriter, r *http.Request) {
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "http://pleskac.org")

	vars := mux.Vars(r)
	idStr := vars[tripId]

	//convert the string id in the URI to an int
	//should not be an error, as the idStr must be digits based on the regex
	id, err := strconv.Atoi(idStr)
	if id < 0 || err != nil {
		fmt.Println("Error parsing", id, "\n", err)
		return
	}

	output := dblayer.GetTrip(id)

	enc := json.NewEncoder(w)
	enc.Encode(output)
}

func TripNameHandler(w http.ResponseWriter, r *http.Request) {
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "http://pleskac.org")

	vars := mux.Vars(r)
	name := vars[tripName]

	id := dblayer.FindTrip(name)
	if id < 0 {
		return
	}

	output := dblayer.GetTrip(id)

	enc := json.NewEncoder(w)
	enc.Encode(output)
}
