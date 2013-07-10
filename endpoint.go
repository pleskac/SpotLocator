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
const email = "userEmail"
const username = "userName"
const displayname = "displayName"
const userpassword = "userPassword"

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
	r.HandleFunc("/api/gps/add/"+password+"/{"+longitude+"}/{"+latitude+"}", AddGPSHandler)
	r.HandleFunc("/api/gps/add/"+password+"/{"+longitude+"}/{"+latitude+"}/{"+gpsType+"}", AddGPSHandler)
	r.HandleFunc("/api/user/add/"+password+"/{"+email+"}/{"+username+"}/{"+displayname+"}/{"+userpassword+"}", AddUserHandler)
	http.ListenAndServe(":8080", r)
}

func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "http://pleskac.org")
	vars := mux.Vars(r)
	useremail := vars[email]
	fmt.Println(useremail)

	dblayer.AddUser(vars[email], vars[username], vars[displayname], vars[userpassword])
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

	standardType := vars[gpsType]
	if standardType != "OK" && standardType != "TRACK" {
		standardType = "TRACK"
	}

	//TODO: FIX THIS!!
	dblayer.AddGPSNow(longFlt, latFlt, "This was sent via iPhone, not SPOT.", standardType, "markpleskac@gmail.com")

	enc := json.NewEncoder(w)
	enc.Encode(standardType)
}

func AddTripHandler(w http.ResponseWriter, r *http.Request) {
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "http://pleskac.org")

	vars := mux.Vars(r)
	name := vars[tripName]

	if name != "" {
		//TODO: this is hardcoded!
		dblayer.CreateTrip(name, "", "markpleskac@gmail.com")
	}

	enc := json.NewEncoder(w)
	enc.Encode(name)
}

func TripListHandler(w http.ResponseWriter, r *http.Request) {
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "http://pleskac.org")

	output := dblayer.GetTripList("markpleskac@gmail.com")

	enc := json.NewEncoder(w)
	enc.Encode(output)
}

func CurrentTripHandler(w http.ResponseWriter, r *http.Request) {
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "http://pleskac.org")

	currentTripId := dblayer.GetCurrentTripId("markpleskac@gmail.com")
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

	id := dblayer.FindTrip(name, "markpleskac@gmail.com")
	if id < 0 {
		return
	}

	output := dblayer.GetTrip(id)

	enc := json.NewEncoder(w)
	enc.Encode(output)
}
