package geolocation

import (
	"encoding/json"
	"net"
	"net/http"
)

func GetIPAddress(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	println("RemoteAddr:", r.RemoteAddr)
	println("X-Forwarded-For:")
	println("X-Real-IP:", r.Header.Get("X-Real-IP"))
	println("IP", ip)
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return ip
}

func IsValidUserAgent(ua string) bool {
	return ua != "" && len(ua) > 10 // basic filter for bots
}

func RespondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
