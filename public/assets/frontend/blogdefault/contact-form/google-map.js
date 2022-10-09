
function initialize() {
    var styles = {
        'creadom': [{
                "featureType": "administrative",
                "stylers": [
                    {"visibility": "on"}
                ]
            },
            {
                "featureType": "road",
                "stylers": [
                    {"visibility": "on"},
                    {"hue": "#6990cb"}
                ]
            },
            {
                "stylers": [
                    {"visibility": "on"},
                    {"hue": "#6990cb"},
                    {"saturation": -50}
                ]
            }
        ]};

    var myLatlng = new google.maps.LatLng(-34.397, 150.644);
    var myOptions = {
        zoom: 10,
        center: myLatlng,
        mapTypeId: google.maps.MapTypeId.ROADMAP,
        //disableDefaultUI: true,
        mapTypeId: 'creadom',
                draggable: true,
        scrollwheel: false,
    }
    var map = new google.maps.Map(document.getElementById("map_creadom"), myOptions);
    var styledMapType = new google.maps.StyledMapType(styles['creadom'], {name: 'creadom'});
    map.mapTypes.set('creadom', styledMapType);

    var marker = new google.maps.Marker({
        position: myLatlng,
        map: map,
        title: ""
    });
    }

google.maps.event.addDomListener(window, 'load', initialize);
google.maps.event.addDomListener(window, 'resize', initialize);

