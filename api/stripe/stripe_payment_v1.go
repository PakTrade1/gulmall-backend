package stripe

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/setupintent"
)

// --- Configuration ---
// Replace with your actual Stripe Secret Key.
// It's highly recommended to use environment variables for sensitive keys.

// --- Request/Response Structures ---

// CreateCustomerRequest is the payload for creating a customer.
type CreateCustomerRequest struct {
	Email    string            `json:"email"`
	Phone    string            `json:"phone"`
	Name     string            `json:"name"`
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata"`
}

// CreateSetupIntentRequest is the payload for creating a SetupIntent.
type CreateSetupIntentRequest struct {
	CustomerID        string   `json:"customer_id"`         // If you already have a customer ID
	PaymentMethodType []string `json:"payment_method_type"` // e.g., ["card","abc"]
}

// CreatePaymentIntentRequest is the payload for creating a PaymentIntent.
type CreatePaymentIntentRequest struct {
	CustomerID        string `json:"customer_id"`
	Amount            int64  `json:"amount"` // Amount in the smallest currency unit (e.g., cents for USD)
	Currency          string `json:"currency"`
	PaymentMethodType string `json:"payment_method_type"` // e.g., "card"
	OffSession        bool   `json:"off_session"`         // For 3DS, if the customer is not present
}

// PaymentIntentResponse contains the client secret for the PaymentIntent.
type PaymentIntentResponse struct {
	ClientSecret string `json:"client_secret"`
}

// SetupIntentResponse contains the client secret for the SetupIntent.
type SetupIntentResponse struct {
	ClientSecret  string `json:"client_secret"`
	SetupIntentID string `json:"setup_intent_id"`
}

func CreateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateCustomerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	params := &stripe.CustomerParams{
		Email: stripe.String(req.Email),
		Phone: stripe.String(req.Phone),
		Name:  stripe.String(req.Name),
		Address: &stripe.AddressParams{
			Line1: stripe.String(req.Address),
		},
		Metadata: req.Metadata,
		// Name: stripe.String("John Doe"),
	}

	c, err := customer.New(params)
	if err != nil {
		log.Printf("Error creating Stripe customer: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create customer: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"customer_id": c.ID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreateSetupIntentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateSetupIntentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Error creating Stripe customer: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate CustomerID if provided
	if req.CustomerID == "" {
		http.Error(w, "CustomerID is required", http.StatusBadRequest)
		return
	}

	params := &stripe.SetupIntentParams{
		Customer: stripe.String(req.CustomerID),
		// Specify the payment method types you want to allow for setup.
		// 'card' is typical for Express Checkout.
		PaymentMethodTypes: stripe.StringSlice(req.PaymentMethodType),
	}

	si, err := setupintent.New(params)
	if err != nil {
		log.Printf("Error creating Stripe SetupIntent: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create setup intent: %v", err), http.StatusInternalServerError)
		return
	}

	response := SetupIntentResponse{
		ClientSecret:  si.ClientSecret,
		SetupIntentID: si.ID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreatePaymentIntentHandler_v1(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreatePaymentIntentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate CustomerID
	if req.CustomerID == "" {
		http.Error(w, "CustomerID is required", http.StatusBadRequest)
		return
	}

	// Ensure amount is positive
	if req.Amount <= 0 {
		http.Error(w, "Amount must be positive", http.StatusBadRequest)
		return
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(req.Amount),
		Currency: stripe.String(req.Currency),
		Customer: stripe.String(req.CustomerID),
		// For Express Checkout, you might not explicitly set a PaymentMethod here initially,
		// as the customer will provide it via the Element.
		// However, you can set `PaymentMethodTypes` to guide the Element.
		PaymentMethodTypes: stripe.StringSlice([]string{req.PaymentMethodType}),
		// Set `OffSession` if you expect the customer might not be present for authentication (e.g., recurring payments).
		// For initial checkout, `off_session` is usually false or not specified.
		OffSession: stripe.Bool(req.OffSession),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		log.Printf("Error creating Stripe PaymentIntent: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create payment intent: %v", err), http.StatusInternalServerError)
		return
	}

	response := PaymentIntentResponse{
		ClientSecret: pi.ClientSecret,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
