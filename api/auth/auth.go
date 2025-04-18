package authWhatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"
	"time"
)

type OTPRequestInfo struct {
	LastRequested time.Time
	RequestCount  int
}

var otpRequests = make(map[string]*OTPRequestInfo)
var mu sync.Mutex

func canSendOTP(phoneNumber string) bool {

	mu.Lock()
	defer mu.Unlock()

	info, exists := otpRequests[phoneNumber]
	now := time.Now()

	if !exists {
		otpRequests[phoneNumber] = &OTPRequestInfo{LastRequested: now, RequestCount: 1}
		return true
	}

	// Cooldown: 1 minute
	if now.Sub(info.LastRequested) < time.Minute {
		return false
	}

	// Reset count every 10 minutes
	if now.Sub(info.LastRequested) > 10*time.Minute {
		info.RequestCount = 0
	}

	if info.RequestCount >= 3 {
		return false
	}

	info.RequestCount++
	info.LastRequested = now
	return true
}

var otpStore = struct {
	sync.Mutex
	data map[string]string
}{data: make(map[string]string)}

func generateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// THIS FUNCTION GET A PHONE NUMBER WITHOUT + SIGN LIKE 447918841539
func sendWhatsAppOTP(phone, otp string) error {
	accessToken := os.Getenv("WHATSAPP_ACCESS_TOKEN")
	phoneNumberID := os.Getenv("WHATSAPP_PHONE_NUMBER_ID")
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/messages", phoneNumberID)

	// JSON payload with new structure
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                phone,
		"type":              "template",
		"template": map[string]interface{}{
			"name": "authenticate",
			"language": map[string]string{
				"code": "en_US",
			},
			"components": []map[string]interface{}{
				{
					"type": "body",
					"parameters": []map[string]interface{}{
						{
							"type": "text",
							"text": otp,
						},
					},
				},
				{
					"type":     "button",
					"sub_type": "url",
					"index":    "0",
					"parameters": []map[string]interface{}{
						{
							"type": "text",
							"text": otp,
						},
					},
				},
			},
		},
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload: %v", err)
		return err
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	// Dump the request for debugging
	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		fmt.Println("Error dumping request:", err)
	} else {
		fmt.Println("Full HTTP request:\n", string(dump))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		return err
	}
	defer resp.Body.Close()

	log.Println("WhatsApp send status:", resp.Status)

	// Optionally log response body
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Response body:", string(respBody))

	return nil
}

func SendOTPHandler(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Phone string `json:"phone"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	otp := generateOTP()

	otpStore.Lock()
	otpStore.data[req.Phone] = otp
	otpStore.Unlock()
	w.Header().Set("Content-Type", "application/json")
	if !canSendOTP(req.Phone) {

		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "Bad",
			"message": "Too many OTP requests. Please wait.",
		})
		return
	}
	go func() {
		err := sendWhatsAppOTP(req.Phone, otp)
		if err != nil {
			log.Println("Error sending OTP:", err)
		}
	}()

	json.NewEncoder(w).Encode(map[string]string{"status": "OK", "message": "OTP sent."})
}

func VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone string `json:"phone"`
		OTP   string `json:"otp"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	w.Header().Set("Content-Type", "application/json")
	otpStore.Lock()
	expectedOTP := otpStore.data[req.Phone]
	otpStore.Unlock()

	if expectedOTP == req.OTP {
		json.NewEncoder(w).Encode(map[string]bool{"verified": true})
	} else {
		json.NewEncoder(w).Encode(map[string]bool{"verified": false})
	}
}
