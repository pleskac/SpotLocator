var infowindow;
var map;
function initialize() {

	//http://pleskac.org/api/SouthAmerica
	infowindow = new google.maps.InfoWindow();

	jQuery.ajax("http://pleskac.org:8080/api/trip/currentTrip").done(
		function(data){
			trips = JSON.parse(data);
				
			var flightPlanCoordinates = [];
			jQuery.each(trips.Coordinates, function(){
				flightPlanCoordinates.push(new google.maps.LatLng(this.Latitude, this.Longitude));
			});

			var mapOptions = {
				zoom: 15,
				center: new google.maps.LatLng(0, 0),
				mapTypeId: google.maps.MapTypeId.TERRAIN,
				streetViewControl: false,
				scrollwheel: false
			};

			map = new google.maps.Map(document.getElementById('map_canvas'), mapOptions);

			var bounds = new google.maps.LatLngBounds();
			jQuery.each(trips.Coordinates, function(){
				var position = new google.maps.LatLng(this.Latitude, this.Longitude);
				createMarker(position, this.Details, this.Color);
				bounds.extend(position);
			});

			map.fitBounds(bounds);

			var flightPath = new google.maps.Polyline({
				path: flightPlanCoordinates,
				strokeColor: "Black",
				strokeOpacity: 1.0,
				strokeWeight: 1.2
			});

			flightPath.setMap(map);
		}
	);
}

function createMarker(_position, name, color) {
	var marker = new google.maps.Marker({
		icon: 	{	path: google.maps.SymbolPath.CIRCLE,
					scale: 5,
					fillColor: color,
					fillOpacity: 1,
					strokeWeight: 1
				},
		position: _position,
		map: map,
		title: this.Title
	});
	
	google.maps.event.addListener(marker, 'click', function() {
		infowindow.setContent(name);
		infowindow.open(map, this);
	});
}