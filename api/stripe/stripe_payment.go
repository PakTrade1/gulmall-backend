package stripe

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/paymentlink"
	"github.com/stripe/stripe-go/v82/price"
	"github.com/stripe/stripe-go/v82/product"
)

// --- Structs for Clear API Contracts ---

// PaymentIntentRequest defines the structure for the /create-payment-intent request body.
type PaymentIntentRequest struct {
	Amount int64 `json:"amount"`
}

// PaymentIntentResponse defines the structure for the /create-payment-intent response body.

// PaymentLinkResponse defines the structure for the /create-payment-link response body.
type PaymentLinkResponse struct {
	URL string `json:"url"`
}

// ErrorResponse defines the structure for a generic JSON error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// --- Main Application Entry Point ---

func MainInitiate() {
	// Best Practice: Load your secret key from an environment variable.
	stripe.Key = "sk_test_51RYGbZ2cxpojAOVyXRxfyjwW3xO5R36gkwgirB4wbdod9XCt5Mqdb93hc1jAwi43t77BqXqICecJT4Wk3l6hkVhD00aZM7peAm"
	if stripe.Key == "" {
		// Fallback for development. Replace with your actual test secret key.
		// IMPORTANT: Do not commit your secret key to version control.
		stripe.Key = os.Getenv("STRIPE_SECRET_KEY_TEST")
		log.Println("WARNING:Using a key for development.")
	}

}

// --- Route Handlers ---

// createPaymentLinkHandler creates a product, a price, and a one-time payment link.
func CreatePaymentLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	productParams := &stripe.ProductParams{Name: stripe.String("Ionic & Go Special")}
	newProduct, err := product.New(productParams)
	if err != nil {
		writeJSONError(w, "Failed to create product", err, http.StatusInternalServerError)
		return
	}

	priceParams := &stripe.PriceParams{
		Currency:   stripe.String(string(stripe.CurrencyUSD)),
		UnitAmount: stripe.Int64(2500), // Example: $25.00
		Product:    stripe.String(newProduct.ID),
	}
	newPrice, err := price.New(priceParams)
	if err != nil {
		writeJSONError(w, "Failed to create price", err, http.StatusInternalServerError)
		return
	}

	paymentLinkParams := &stripe.PaymentLinkParams{
		LineItems: []*stripe.PaymentLinkLineItemParams{{
			Price:    stripe.String(newPrice.ID),
			Quantity: stripe.Int64(1),
		}},
	}
	pl, err := paymentlink.New(paymentLinkParams)
	if err != nil {
		writeJSONError(w, "Failed to create payment link", err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, PaymentLinkResponse{URL: pl.URL}, http.StatusOK)
}

// createPaymentIntentHandler creates a Payment Intent to prepare for a card or wallet payment.
func CreatePaymentIntentHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PaymentIntentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Invalid request body", err, http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		writeJSONError(w, "Invalid amount", nil, http.StatusBadRequest)
		return
	}

	params := &stripe.PaymentIntentParams{

		Amount:   stripe.Int64(req.Amount),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}
	pi, err := paymentintent.New(params)
	if err != nil {
		writeJSONError(w, "Failed to create payment intent", err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, PaymentIntentResponse{ClientSecret: pi.ClientSecret}, http.StatusOK)
}

// writeJSON is a helper to write a JSON response with a status code.
func writeJSON(w http.ResponseWriter, v interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}

// writeJSONError is a helper to log an error and write a JSON error response.
func writeJSONError(w http.ResponseWriter, message string, err error, statusCode int) {
	if err != nil {
		log.Printf("❌ ERROR: %s: %v", message, err)
	} else {
		log.Printf("❌ ERROR: %s", message)
	}
	resp := ErrorResponse{Error: message}
	writeJSON(w, resp, statusCode)
}
