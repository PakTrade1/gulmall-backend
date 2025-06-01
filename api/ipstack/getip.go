package ipstack

import (
	"net/http"
	"pak-trade-go/api/geolocation"
)

func GetUserLocation(w http.ResponseWriter, r *http.Request) {
	ip := geolocation.GetIPAddress(r)
	ua := r.UserAgent()

	if !geolocation.IsValidUserAgent(ua) {
		http.Error(w, "Invalid User-Agent", http.StatusForbidden)
		return
	}

	// Check cache first
	if data, found := geolocation.GetFromCache(ip); found {
		geolocation.RespondJSON(w, data)
		return
	}

	// Fetch from IPStack
	geoData, err := geolocation.FetchGeoData(ip)
	if err != nil {
		http.Error(w, "Unable to get location", http.StatusInternalServerError)
		return
	}

	geolocation.SetToCache(ip, geoData)
	geolocation.RespondJSON(w, geoData)
}
