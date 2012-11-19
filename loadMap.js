var infowindow;
var map;
function initialize() {
//call rest service here... get the center
//http://pleskac.org/api/SouthAmerica
infowindow = new google.maps.InfoWindow();
jQuery.ajax("http://pleskac.org/api/SouthAmerica").done(
	function(data){
		trips = JSON.parse(data);
        
	var flightPlanCoordinates = [];
	jQuery.each(trips.Coordinates, function(){
		flightPlanCoordinates.push(new google.maps.LatLng(this.Latitude, this.Longitude));
	});
	console.log(flightPlanCoordinates);


	var centerCoordinate = new google.maps.LatLng(trips.CenterLat, trips.CenterLong);
        console.log(centerCoordinate);
	var mapOptions = {
          zoom: trips.Zoom,
          center: centerCoordinate,
          mapTypeId: google.maps.MapTypeId.TERRAIN,
	  streetViewControl: false,
	  scrollwheel: false
        };

	console.log(mapOptions);

        map = new google.maps.Map(document.getElementById('map_canvas'), mapOptions);

	jQuery.each(trips.Coordinates, function(){
		console.log(this);
		createMarker(this.Latitude, this.Longitude, this.Details, this.Color);
	/*	var position = new google.maps.LatLng(this.Latitude, this.Longitude);
        	var marker = new google.maps.Marker({
			position: position, 
			map: map,
			title: this.Title
 		});

          	//marker.setTitle(this.Title);

		//var infowindow = new google.maps.InfoWindow({
          	//	content: this.Details
        	//});

        	google.maps.event.addListener(marker, 'click', function() {
          		infowindow.setContent(this.Title);
			//infowindow.open(marker.get('map'), marker);
			infowindow.open(map, this);
        	});*/
	});
       
	//var bounds = new google.maps.LatLngBounds(southWest, northEast);
        //map.fitBounds(bounds);

        var flightPath = new google.maps.Polyline({
          path: flightPlanCoordinates,
          strokeColor: "#556b2f",
          strokeOpacity: 1.0,
          strokeWeight: 2
        });

        flightPath.setMap(map);		
		
	}
);
}

function createMarker(latitude, longitude, name, color) {
	console.log(name);
	var position = new google.maps.LatLng(latitude, longitude);
	
        var marker = new google.maps.Marker({
                        icon: {	path: google.maps.SymbolPath.CIRCLE,
 				scale: 5,
				fillColor: color,
				fillOpacity: 1,
				strokeWeight: 1
			},
			position: position,
                        map: map,
                        title: this.Title
                });
        google.maps.event.addListener(marker, 'click', function() {
          infowindow.setContent(name);
          infowindow.open(map, this);
        });
}
