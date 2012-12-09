# SpotLocator #

## Problem ##
The Spot GPS is a device for "checking in" in the wilderness. It works where cellphones don't. It is used primarily as a safety device, alerting local search and rescue as well as anybody else that you wish to notify. It even has the option to share locations to Twitter, Facebook, or a "shared map" (basically an embedded Google Map that doesn't work in Firefox nor Chrome).

## Solution ##
Since the Spot website is horrendous, doesn't support my favorite browsers, and only is able to display the last 7 days of data, I decided to rewrite my own version which I am able to embed in a blog. I ping the public API of my GPS locations, looking for new data. When a new location is found, I save it locally. I also serve up all of this information to my own API, which I consume in javascript, creating a customized Google Map of my locations. I am able to associate GPS locations with individual trips and am able to save more than 7 days(!).

Originally, I scraped the Twitter API. Then Twitter implemented mandatory OAuth, causing much frustration on my end. I decided to scrap Twitter completely and use the Spot API. It is easier, despite some weird naming conventions. I'm also able to access all types of messages, rather than just the "OK" messages that Twitter was privy to.

As a disclaimer, Spot started allowing users to save longer "trips" through a proprietary "Adventures" website (which only works in IE and Safari). It still is less customizable and less accessable than my solution.

## API ##
The base URL is ```pleskac.org:8080``` for my instance.

``` /api/trip/list ``` Returns a list of Trip objects.

``` /api/trip/currentTrip ``` Returns a Trip object of the current trip. Returns an empty trip if there are no current trips.

``` /api/trip/id/{int} ``` Returns a trip object matching the key. Returns an empty trip if it is not found in the database.

``` /api/trip/name/{string} ```Returns a single trip object of the same name. It does string matching and returns the first trip found if multiple trips match the string. Not sure what this will be used for. The list of trip objects could be used for searching and selecting. The id can be used to get for a specific trip. This is basically for fun and/or testing.

``` /api/trip/add/{PASSWORD}/{string} ``` Adds a new trip and sets it to the current trip. The password is checked against the database. It's not very secure, more a deterrant to people abusing this API. All checkins after creating the new trip will go to this new current trip.

``` /api/gps/add/{PASSWORD}/{double}/{double}/{string}/{string} ``` TODO. Adds a new GPS location with the passed in parameters. The order is: longitude, latitude, type, message. Automatically set to the timestamp of the call to the API and added to the current trip.

## Configuration ##
Still need to code, then documentation will come
* Account config
* Server setup
* HTML usage

## Files ##

### Package main ###

``` main.go ```
Main function. Calls spot.go to get new locations. Sends them to mysql.go to save them. Keeps track of the latest location, persisted in MySQL.

``` spot.go ```
Deals with the Spot API. Returns a list of messages. Gets around a weird case of having 1 message vs multiple messages in json.

``` endpoint.go ```
Serves my custom API. See the API section for URIs this handles

### Package dblayer ###

``` dblayer/*.go ```
Contacts the database. Saving and retrieving information supported. Formats the outputs to nice objects (maybe too much formatting).

### Other Files ###
``` loadMap.json ```
Consumes my custom API, creating a Google Map which is easily put into any ```<div>``` named ```map_canvas```. Also includes a select box that allows selection of all trips. The current trip is denoted by ```(Current)``` if it exists. Example page at ```pleskac.org/map.html```. I then embed it in an iframe in my blog. If I don't do this, the WordPress theme's CSS will make Google Map's CSS all funky.

``` map.html ```
This displays the map and select box wrapped up in html.

## TODO ##
* Timestamp is in wrong timezone. Adjust to reflect the browser's timezone?
* TESTS!!!
* Move spot.go to new package to allow for reuse
* Add Foursquare integration
* Make configureable with config files