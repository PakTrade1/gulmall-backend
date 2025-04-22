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
	"pak-trade-go/api/mammals"
	"pak-trade-go/api/signin"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OTPRequestInfo struct {
	LastRequested time.Time
	RequestCount  int
}

var otpRequests = make(map[string]*OTPRequestInfo)
var mu sync.Mutex

func canSendOTP(phoneNumber string) bool {
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: error loading .env file (ignored in production)")
		}
	}
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
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Phone string `json:"phone"`
		OTP   string `json:"otp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Phone == "" || req.OTP == "" {
		http.Error(w, "Phone or OTP is missing", http.StatusBadRequest)
		return
	}

	json.NewDecoder(r.Body).Decode(&req)

	otpStore.Lock()
	expectedOTP := otpStore.data[req.Phone]
	otpStore.Unlock()

	user, err := signin.FindUserByPhone(req.Phone)
	if err != nil {
		// Handle DB error (optional: log or report)
	}
	if user == nil {
		// User not found — create a new one
		// user = CreateUser(req.Phone)
	}

	if expectedOTP == req.OTP {
		json.NewEncoder(w).Encode(map[string]bool{"verified": true})
	} else {
		json.NewEncoder(w).Encode(map[string]bool{"verified": false})
	}
}

func CreateUser(phone string) (*mammals.User, error) {
	var user mammals.User
	user.ID = primitive.NewObjectID()
	user.Credit = 5
	user.AdsRemaining = 5
	user.ServerDate = primitive.NewDateTimeFromTime(time.Now())
	user.CreationDate = time.Now().Format(time.RFC3339)
	user.LastSignedIn = time.Now().Format(time.RFC3339)
	user.IsEmailVerified = false
	user.AccountStatus = true
	user.BusinessPhone = "N/A"
	user.IsBusiness = false
	user.PublicID = mammals.GetNextPublicID()

	planID := "64735fe18f737b74c13bd6d3"
	Planid, _ := primitive.ObjectIDFromHex(planID)
	user.PlanID = Planid
	// coll := docking.PakTradeDb.Collection("Mammalas_login")
	// _, err = coll.InsertOne(context.TODO(), user)
	// if err != nil {
	// 	return nil, err
	// }
	return &user, nil

}
