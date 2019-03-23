 
const mapStyle = [
    {
      "featureType": "administrative",
      "elementType": "all",
      "stylers": [
        {
          "visibility": "on"
        },
        {
          "lightness": 33
        }
      ]
    },
    {
      "featureType": "landscape",
      "elementType": "all",
      "stylers": [
        {
          "color": "#f2e5d4"
        }
      ]
    },
    {
      "featureType": "poi.park",
      "elementType": "geometry",
      "stylers": [
        {
          "color": "#c5dac6"
        }
      ]
    },
    {
      "featureType": "poi.park",
      "elementType": "labels",
      "stylers": [
        {
          "visibility": "on"
        },
        {
          "lightness": 20
        }
      ]
    },
    {
      "featureType": "road",
      "elementType": "all",
      "stylers": [
        {
          "lightness": 20
        }
      ]
    },
    {
      "featureType": "road.highway",
      "elementType": "geometry",
      "stylers": [
        {
          "color": "#c5c6c6"
        }
      ]
    },
    {
      "featureType": "road.arterial",
      "elementType": "geometry",
      "stylers": [
        {
          "color": "#e4d7c6"
        }
      ]
    },
    {
      "featureType": "road.local",
      "elementType": "geometry",
      "stylers": [
        {
          "color": "#fbfaf7"
        }
      ]
    },
    {
      "featureType": "water",
      "elementType": "all",
      "stylers": [
        {
          "visibility": "on"
        },
        {
          "color": "#acbcc9"
        }
      ]
    }
  ]; 

  function readTextFile(file, callback) {
    var rawFile = new XMLHttpRequest();
    rawFile.overrideMimeType("application/json");
    rawFile.open("GET", file, true);
    rawFile.onreadystatechange = function() {
        if (rawFile.readyState === 4 && rawFile.status == "200") {
            callback(rawFile.responseText);
        }
    }
    rawFile.send(null);
  }
 

  // Escapes HTML characters in a template literal string, to prevent XSS.
  // See https://www.owasp.org/index.php/XSS_%28Cross_Site_Scripting%29_Prevention_Cheat_Sheet#RULE_.231_-_HTML_Escape_Before_Inserting_Untrusted_Data_into_HTML_Element_Content
  function sanitizeHTML(strings) {
    const entities = {'&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;'};
    let result = strings[0];
    for (let i = 1; i < arguments.length; i++) {
      result += String(arguments[i]).replace(/[&<>'"]/g, (char) => {
        return entities[char];
      });
      result += strings[i];
    }
    return result;
  }
   
 
  initMap = async (consumerAccount, latitude, lontitude) => { 
     
    const map = new google.maps.Map(document.getElementById('container'), {
      zoom: 13,
      center: {lat: latitude, lng: lontitude}, 
      styles: mapStyle
    });
   
    //FILE READER - map.data.loadGeoJson('/public/geoData.json');
    
    // call api
    await map.data.loadGeoJson('/geodata/consumer/' + consumerAccount);
   
    // Define the custom marker icons, using the store's "category".
    map.data.setStyle(feature => {
      const label = `${feature.getProperty('name')}`
      return {
        title: `${feature.getProperty('name')}`,
        label: label, 
        icon: {
          url: `/public/images/icon_${feature.getProperty('icon')}.png`,
          scaledSize: new google.maps.Size(48, 48), 
        }
      };
    });
   

    const apiKey = 'AIzaSyB_2WtkzS6oe6y07LXAtC1Dwd5dI_L3tj0';
    const infoWindow = new google.maps.InfoWindow();
    infoWindow.setOptions({pixelOffset: new google.maps.Size(0, -30)});
  
    // Show the information for a store when its marker is clicked.
    map.data.addListener('click', event => {
  
        const account = event.feature.getProperty('account');
        const name = event.feature.getProperty('name');
        const category = event.feature.getProperty('category');
        const startDate = event.feature.getProperty('startdate');
        const endDate = event.feature.getProperty('enddate');
        const icon = event.feature.getProperty('icon');
        const position = event.feature.getGeometry().get();
        var content;
        
        if (icon === 'farm') {
            content = sanitizeHTML`
            <img style="float:left; width:200px; margin-top:30px" src="/public/images/icon_${icon}.png">
            <div style="margin-left:220px; margin-bottom:20px;">
              <h2>${name}</h2><p><b>Account:</b> ${account}</p><b>Category:</b> ${category} producer</p>
              <b>Production Timeline:</b> 
              <p>&nbsp;&nbsp;<b>Start Date:</b> ${startDate}<br/>&nbsp;&nbsp;<b>End Date:</b> ${endDate}</p> 
            </div>
          `;
        } else {
            content = sanitizeHTML`
            <img style="float:left; width:200px; margin-top:30px" src="/public/images/icon_${icon}.png">
            <div style="margin-left:220px; margin-bottom:20px;">
              <h2>${name}</h2><p><b>Account:</b> ${account}</p><b>Category:</b> ${category} consumer</p>
              <b>Required Timeline</b>  
              <p>&nbsp;&nbsp;<b>Start Date:</b> ${startDate}<br/>&nbsp;&nbsp;<b>End Date:</b> ${endDate}</p> 
            </div>
          `;
        }
          
        infoWindow.setContent(content);
        infoWindow.setPosition(position);
        infoWindow.open(map);
    });
  
  }