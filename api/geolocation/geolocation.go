package geolocation

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

const ipstackAPIKey = "5afee7526c541569085a261da449374d"

type IPAPIResponse struct {
	IP          string `json:"ip"`
	City        string `json:"city"`
	Region      string `json:"region"`
	Country     string `json:"country_name"`
	CountryCode string `json:"country_code"`
	Currency    string `json:"currency"`
}

func GetLocationFromIP(ip string) (*IPAPIResponse, error) {
	url := fmt.Sprintf("https://ipapi.co/%s/json", ip)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make ipapi request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Print raw JSON response
	fmt.Println("Raw IPAPI JSON response:")
	fmt.Println(string(bodyBytes))

	// Decode JSON into struct
	var data IPAPIResponse
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, fmt.Errorf("failed to decode ipapi response: %w", err)
	}

	return &data, nil
}

func GetIP(r *http.Request) string {
	// Check for X-Forwarded-For header (set by proxy/load balancer)
	ip := r.Header.Get("X-Forwarded-For")
	fmt.Println("RemoteAddr:", r.RemoteAddr)
	fmt.Println("X-Forwarded-For:")
	fmt.Println("X-Real-IP:", r.Header.Get("X-Real-IP"))
	println("IP", ip)
	if ip != "" {
		// May contain multiple IPs — take the first one
		parts := strings.Split(ip, ",")
		return strings.TrimSpace(parts[0])
	}

	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // fallback if parsing fails
	}
	return ip
}

func IPHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := GetIP(r)

	data, err := GetLocationFromIP(clientIP)
	if err != nil {
		http.Error(w, "Failed to get location: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Optional: Log location data for debugging
	fmt.Printf("Full Location Data: %+v\n", *data)
	fmt.Println("Client IP is:", clientIP)

	// Set content type and respond with JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
	// you can now pass clientIP to location APIs like ipstack
}
