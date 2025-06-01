package geolocation

import (
	"encoding/json"
	"fmt"
	"io"
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
	println(resp.Body)
	println(err)
	if err != nil {
		return GeoData{}, fmt.Errorf("failed to make ipapi request: %w", err)
	}
	if err != nil {
		return GeoData{}, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return GeoData{}, fmt.Errorf("failed to read response body: %w", err)
	}
	fmt.Println("IP API JSON response:")
	fmt.Println(string(bodyBytes))
	var data GeoData

	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return GeoData{}, fmt.Errorf("failed to decode ipapi response: %w", err)
	}

	return data, nil
}
