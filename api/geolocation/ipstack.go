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
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
}

func FetchGeoData(ip string) (GeoData, error) {
	println("IP: ", ip)
	url := fmt.Sprintf("https://api.ipinfo.io/lite/%s/?token=6794b10129b8b5", ip)
	println(url)
	resp, err := http.Get(url)
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
