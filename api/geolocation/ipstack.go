package geolocation

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//const ipstackAPIKey = "5afee7526c541569085a261da449374d"
//const ipstackBaseURL = "http://api.ipstack.com/"

type GeoData struct {
	IP          string `json:"ip"`
	City        string `json:"city"`
	Region      string `json:"region"`
	Country     string `json:"country_name"`
	CountryCode string `json:"country_code"`
	Currency    string `json:"currency"`
}

func FetchGeoData(ip string) (GeoData, error) {

	url := fmt.Sprintf("https://ipapi.co/%s/json/", ip)
	println(url)
	resp, err := http.Get(url)
	print(resp)
	if err != nil {
		return GeoData{}, err
	}
	defer resp.Body.Close()

	var data GeoData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return GeoData{}, err
	}

	return data, nil
}
