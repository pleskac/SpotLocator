var infowindow;
var map;
function initialize() {
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

	});
       
        var flightPath = new google.maps.Polyline({
          path: flightPlanCoordinates,
          strokeColor: "red",
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
