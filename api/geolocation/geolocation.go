package geolocation

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
)

const ipstackAPIKey = "YOUR_IPSTACK_API_KEY"

type IPStackResponse struct {
	IP          string  `json:"ip"`
	City        string  `json:"city"`
	RegionName  string  `json:"region_name"`
	CountryName string  `json:"country_name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

func GetLocationFromIP(ip string) (*IPStackResponse, error) {
	url := fmt.Sprintf("http://api.ipstack.com/%s?access_key=%s", ip, ipstackAPIKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make IPStack request: %w", err)
	}
	defer resp.Body.Close()

	var data IPStackResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode IPStack response: %w", err)
	}

	return &data, nil
}

func GetIP(r *http.Request) string {
	// Check for X-Forwarded-For header (set by proxy/load balancer)
	ip := r.Header.Get("X-Real-IP")
	fmt.Println("RemoteAddr:", r.RemoteAddr)
	fmt.Println("X-Forwarded-For:", r.Header.Get("X-Forwarded-For"))
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
	fmt.Println("Client IP is:", clientIP)
	// you can now pass clientIP to location APIs like ipstack
}
