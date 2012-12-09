var infowindow;
var map;

//jQuery(document).ready(firstLoad);
jQuery(document).ready(loadSelectBox);

function loadSelectBox(){
	jQuery.ajax("http://pleskac.org:8080/api/trip/list").done(
		function(data){
			tripList = JSON.parse(data);

			jQuery.each(tripList, function(){
				var id = this.TripId;
				var name = this.TripName;
				var current = this.IsCurrent;

				if(current == 0){
					jQuery("#TripSelectBox").append("<option value='" + id + "'>" + name + "</option>");
				}else{
					jQuery("#TripSelectBox").append("<option value='" + id + "'>(Current) " + name + "</option>");
				}
			});
			jQuery("#TripSelectBox").change(changeMap);
		}
	);

	jQuery("#TripSelectBox").change(changeMap);
}

function changeMap(){
	  var id = jQuery("#TripSelectBox").val();
	  loadMap(id);
}

function firstLoad(){
	loadMap(-1);
}

function loadMap(loadId){
	infowindow = new google.maps.InfoWindow();

	var apiUrl = ""
	if(loadId == -1){
		apiUrl = "http://pleskac.org:8080/api/trip/currentTrip";
	}else{
		apiUrl = "http://pleskac.org:8080/api/trip/id/" + loadId;
	}

	jQuery.ajax(apiUrl).done(
		function(data){
			trips = JSON.parse(data);
				
			var flightPlanCoordinates = [];

			if(trips.Coordinates != null){
				jQuery.each(trips.Coordinates, function(){
					flightPlanCoordinates.push(new google.maps.LatLng(this.Latitude, this.Longitude));
				});
			}

			var mapOptions = {
				zoom: 13,
				center: new google.maps.LatLng(0, 0),
				mapTypeId: google.maps.MapTypeId.TERRAIN,
				streetViewControl: false,
				scrollwheel: false
			};

			map = new google.maps.Map(document.getElementById('map_canvas'), mapOptions);

			var bounds = new google.maps.LatLngBounds();

			if(trips.Coordinates != null){
				jQuery.each(trips.Coordinates, function(){
					var position = new google.maps.LatLng(this.Latitude, this.Longitude);
					createMarker(position, this.Details, this.Color);
					bounds.extend(position);
				});
			}

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