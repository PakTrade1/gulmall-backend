package cart

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	docking "pak-trade-go/Docking"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChangeSet struct {
	Old map[string]interface{}
	New map[string]interface{}
}

type Cart struct {
	ID              primitive.ObjectID `bson:"_id" json:"cart_id"`
	ItemID          primitive.ObjectID `bson:"item_id" json:"item_id"`
	ColorID         primitive.ObjectID `bson:"color_id" json:"color_id"`
	Quantity        int                `bson:"quantity" json:"quantity"`
	TotalPrice      float64            `bson:"total_price" json:"total_price"`
	Discount        string             `bson:"discount" json:"discount"`
	PaymentMethod   primitive.ObjectID `bson:"payment_method" json:"payment_method"`
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	SellerId        primitive.ObjectID `bson:"seller_id" json:"seller_id"`
	DeliveryStatus  string             `bson:"delivery_status" json:"delivery_status"`
	OrderDate       time.Time          `bson:"order_date" json:"order_date"`
	SizeID          primitive.ObjectID `bson:"size_id" json:"size_id"`
	Currency        string             `bson:"currency" json:"currency"`
	Category        string             `bson:"category" json:"category"`
	SubCategory     string             `bson:"sub_category" json:"sub_category"`
	CreatedAt       time.Time          `bson:"created_at"`
	IsModified      bool               `bson:"isModified" json:"isModified"`
	DeliveredOn     time.Time          `bson:"delivery_date" json:"delivery_date"`
	DeliveredBy     primitive.ObjectID `bson:"deliver_by" json:"deliver_by"`
	OrderNumber     int64              `bson:"order_number" json:"order_number"`
	OrderVerified   bool               `bson:"order_verified" json:"order_verified"`
	OrderVerifiedBy bool               `bson:"order_verified_by" json:"order_verified_by"`
}

// CartOutput represents the final structure of an order item in the response
type CartOutput struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	ItemID          primitive.ObjectID `bson:"item_id"`                     // From your original Cart struct
	Color           string             `bson:"color"`                       // From color lookup
	Quantity        int                `bson:"quantity"`                    // From your original Cart struct
	TotalPrice      float64            `bson:"total_price"`                 // From your original Cart struct
	Discount        string             `bson:"discount"`                    // From your original Cart struct
	PaymentMode     PaymentModeDoc     `bson:"payment_mode"`                // From payment_services lookup
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`      // From your original Cart struct
	BuyerInfo       BuyerSellerInfo    `bson:"buyer_info"`                  // From Mammalas_login lookup for user_id
	SellerId        primitive.ObjectID `bson:"seller_id"`                   // From your original Cart struct
	SellerInfo      BuyerSellerInfo    `bson:"seller_info,omitempty"`       // From Mammalas_login lookup for seller_id
	DeliveryStatus  string             `bson:"delivery_status"`             // From your original Cart struct
	OrderDate       time.Time          `bson:"order_date"`                  // From your original Cart struct
	Size            string             `bson:"size"`                        // From size lookup
	Currency        string             `bson:"currency"`                    // From your original Cart struct
	Category        string             `bson:"category"`                    // From your original Cart struct
	SubCategory     string             `bson:"sub_category"`                // From your original Cart struct
	IsModified      bool               `bson:"isModified"`                  // From your original Cart struct
	DeliveredOn     *time.Time         `bson:"delivery_date,omitempty"`     // Pointer to time.Time for potentially null values, mapping to your struct's delivered_on
	DeliveredBy     primitive.ObjectID `bson:"deliver_by"`                  // Original DeliverBy ID if needed, or map to a string name if lookup provides it
	OrderNumber     int64              `bson:"order_number"`                // From your original Cart struct
	OrderVerified   bool               `bson:"order_verified"`              // From your original Cart struct
	OrderVerifiedBy *OrderVerifiedBy   `bson:"order_verified_by,omitempty"` // From employee lookup
}

// OrderVerifiedBy represents the structure of the order_verified_by object from employee lookup
type OrderVerifiedBy struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	EmpID string             `bson:"emp_id,omitempty"`
	Name  string             `bson:"name,omitempty"`
}

// BuyerSellerInfo represents the structure for buyer_info and seller_info
type BuyerSellerInfo struct {
	CountryName  string `bson:"countryName,omitempty"`
	Currency     string `bson:"currency,omitempty"`
	PrimaryPhone string `bson:"primaryPhone,omitempty"`
}

// PaymentModeDoc represents the structure of the payment_mode object
type PaymentModeDoc struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Country string             `bson:"country,omitempty"`
	Name    interface{}        `bson:"name,omitempty"` // Can be string or nested PaymentName
}

// CartResponse structure to hold the orders and summary
type CartResponse struct {
	Orders  []bson.M    `bson:"orders"`
	Summary CartSummary `bson:"summary"`
}

// CartSummary for the aggregated data
type CartSummary struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	TotalOrders     int                `bson:"total_orders"`
	DeliveredOrders int                `bson:"delivered_orders"`
	PendingOrders   int                `bson:"pending_orders"`
}

type Change struct {
	Field string      `json:"field"`
	From  interface{} `json:"from"`
	To    interface{} `json:"to"`
}

type FlatChange struct {
	OrderID    string    `json:"orderId"`
	ModifiedBy string    `json:"modify_by"`
	ModifiedAt time.Time `json:"updated_at"`
	Changes    struct {
		Field string      `json:"field"`
		From  interface{} `json:"from"`
		To    interface{} `json:"to"`
	} `json:"changes"`
}

type UserInfo struct {
	PublicID int64  `bson:"publicId" json:"public_id"`
	Phone    string `bson:"primaryPhone" json:"phone"`
}

type AddToCartPayload struct {
	Orders []Cart `json:"orders"`
}

func GetDetailedCartItemsHandler(cartCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Get optional query params
		status := r.URL.Query().Get("status")
		orderLocation := r.URL.Query().Get("order_location")
		startDateStr := r.URL.Query().Get("start_date")
		endDateStr := r.URL.Query().Get("end_date")
		dateStr := r.URL.Query().Get("date")
		// Dynamic match stage
		matchStage := bson.D{}

		if dateStr != "" {
			orderDate, err := time.Parse("2006-01-02", dateStr)
			if err == nil {
				start := orderDate.UTC()
				end := start.Add(24 * time.Hour)
				matchStage = append(matchStage, bson.E{Key: "order_date", Value: bson.D{
					{"$gte", start},
					{"$lt", end},
				}})
			}
		}

		if startDateStr != "" || endDateStr != "" {
			dateFilter := bson.D{}
			if startDateStr != "" {
				startDate, err := time.Parse("2006-01-02", startDateStr)
				if err == nil {
					dateFilter = append(dateFilter, bson.E{Key: "$gte", Value: startDate})
				}
			}
			if endDateStr != "" {
				endDate, err := time.Parse("2006-01-02", endDateStr)
				if err == nil {
					// Add 1 day to make it inclusive of end date
					dateFilter = append(dateFilter, bson.E{Key: "$lte", Value: endDate.Add(24 * time.Hour)})
				}
			}
			if len(dateFilter) > 0 {
				matchStage = append(matchStage, bson.E{Key: "order_date", Value: dateFilter})
			}
		}

		if status != "" {
			matchStage = append(matchStage, bson.E{Key: "delivery_status", Value: status})
		}
		if orderLocation != "" {
			matchStage = append(matchStage, bson.E{Key: "buyer_info.countryName", Value: orderLocation})
		}

		pipeline := mongo.Pipeline{}

		// Add $match stage only if filters exist
		if len(matchStage) > 0 {
			pipeline = append(pipeline, bson.D{{"$match", matchStage}})
		}

		// Rest of your aggregation pipeline
		pipeline = append(pipeline,
			// Lookup buyer_info
			bson.D{{"$lookup", bson.D{
				{"from", "Mammalas_login"},
				{"localField", "user_id"},
				{"foreignField", "_id"},
				{"as", "buyer_info"},
			}}},
			bson.D{{"$lookup", bson.D{
				{"from", "Mammalas_login"},
				{"localField", "seller_id"},
				{"foreignField", "_id"},
				{"as", "seller_info"},
			}}},
			bson.D{{"$lookup", bson.D{
				{"from", "color"},
				{"localField", "color_id"},
				{"foreignField", "_id"},
				{"as", "color"},
			}}},
			bson.D{{"$addFields", bson.D{
				{"color", bson.D{
					{"$arrayElemAt", bson.A{"$color.name", 0}},
				}},
			}}},
			bson.D{{"$lookup", bson.D{
				{"from", "size"},
				{"localField", "size_id"},
				{"foreignField", "_id"},
				{"as", "size"},
			}}},
			bson.D{{"$addFields", bson.D{
				{"size", bson.D{
					{"$arrayElemAt", bson.A{"$size.name", 0}},
				}},
			}}},
			bson.D{{"$lookup", bson.D{
				{"from", "payment_services"},
				{"localField", "payment_method"},
				{"foreignField", "_id"},
				{"as", "payment_mode"},
			}}},
			bson.D{{"$addFields", bson.D{
				{"payment_mode", bson.D{
					{"$arrayElemAt", bson.A{"$payment_mode.name", 0}},
				}},
			}}},
			bson.D{{"$addFields", bson.D{
				{"buyer_info", bson.D{
					{"$arrayElemAt", bson.A{"$buyer_info", 0}},
				}},
				{"seller_info", bson.D{
					{"$arrayElemAt", bson.A{"$seller_info", 0}},
				}},
			}}},
			bson.D{{"$lookup", bson.D{
				{"from", "cloths"},
				{"localField", "item_id"},
				{"foreignField", "_id"},
				{"as", "item"},
			}}},
			bson.D{{"$addFields", bson.D{
				{"item_name", bson.D{
					{"$arrayElemAt", bson.A{"$item.name", 0}},
				}},
				{"item_fabric", bson.D{
					{"$arrayElemAt", bson.A{"$item.fabric", 0}},
				}},
			}}},
			bson.D{{"$lookup", bson.D{
				{"from", "fabric"},
				{"localField", "item_fabric"},
				{"foreignField", "_id"},
				{"as", "item_fabric"},
			}}},
			bson.D{{"$addFields", bson.D{
				{"item_fabric", bson.D{
					{"$arrayElemAt", bson.A{"$item_fabric.name", 0}},
				}},
			}}},
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "employee"},
						{"localField", "order_verified_by"},
						{"foreignField", "_id"},
						{"as", "order_verified_by"},
					},
				},
			},
			bson.D{
				{"$addFields",
					bson.D{
						{"order_verified_by",
							bson.D{
								{"$arrayElemAt",
									bson.A{
										"$order_verified_by",
										0,
									},
								},
							},
						},
					},
				},
			},

			bson.D{{"$project", bson.D{
				{"seller_info.primaryPhone", 1},
				{"seller_info.currency", 1},
				{"seller_info.countryName", 1},
				{"buyer_info.primaryPhone", 1},
				{"buyer_info.currency", 1},
				{"buyer_info.countryName", 1},
				{"payment_mode", 1},
				{"order_date", 1},
				{"color", 1},
				{"total_price", 1},
				{"delivery_status", 1},
				{"size", 1},
				{"quantity", 1},
				{"discount", 1},
				{"currency", 1},
				{"item_fabric", 1},
				{"item_name", 1},
				{"item_id", 1},
				{"delivery_date", 1},
				{"deliver_by", 1},
				{"category", 1},
				{"sub_category", 1},
				{"order_number", 1},
				{"order_verified", 1},
				{"isModified", 1},
				{"order_verified_by", 1},
			}}},
		)

		cursor, err := cartCollection.Aggregate(ctx, pipeline)
		if err != nil {
			http.Error(w, "Aggregation error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		var results []bson.M
		if err := cursor.All(ctx, &results); err != nil {
			http.Error(w, "Error decoding cart items: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}

func GetDetailedCartItemsHandler_v2(cartCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Parse query params
		status := r.URL.Query().Get("status")
		orderLocation := r.URL.Query().Get("order_location")
		startDateStr := r.URL.Query().Get("start_date")
		endDateStr := r.URL.Query().Get("end_date")
		dateStr := r.URL.Query().Get("date")

		matchStage := bson.D{}

		// Date filters
		if dateStr != "" {
			orderDate, err := time.Parse("2006-01-02", dateStr)
			if err == nil {
				start := orderDate.UTC()
				end := start.Add(24 * time.Hour)
				matchStage = append(matchStage, bson.E{Key: "order_date", Value: bson.D{
					{"$gte", start},
					{"$lt", end},
				}})
			}
		} else if startDateStr != "" || endDateStr != "" {
			dateFilter := bson.D{}
			if startDateStr != "" {
				startDate, err := time.Parse("2006-01-02", startDateStr)
				if err == nil {
					dateFilter = append(dateFilter, bson.E{Key: "$gte", Value: startDate})
				}
			}
			if endDateStr != "" {
				endDate, err := time.Parse("2006-01-02", endDateStr)
				if err == nil {
					dateFilter = append(dateFilter, bson.E{Key: "$lte", Value: endDate.Add(24 * time.Hour)})
				}
			}
			if len(dateFilter) > 0 {
				matchStage = append(matchStage, bson.E{Key: "order_date", Value: dateFilter})
			}
		}

		if status != "" {
			matchStage = append(matchStage, bson.E{Key: "delivery_status", Value: status})
		}
		if orderLocation != "" {
			matchStage = append(matchStage, bson.E{Key: "buyer_info.countryName", Value: orderLocation})
		}

		// Create shared stages
		lookupStages := mongo.Pipeline{
			// Apply match
			{{
				"$match", matchStage,
			}},
			{{"$lookup", bson.D{{"from", "Mammalas_login"}, {"localField", "user_id"}, {"foreignField", "_id"}, {"as", "buyer_info"}}}},
			{{"$lookup", bson.D{{"from", "Mammalas_login"}, {"localField", "seller_id"}, {"foreignField", "_id"}, {"as", "seller_info"}}}},
			{{"$lookup", bson.D{{"from", "color"}, {"localField", "color_id"}, {"foreignField", "_id"}, {"as", "color"}}}},
			{{"$addFields", bson.D{{"color", bson.D{{"$arrayElemAt", bson.A{"$color.name", 0}}}}}}},
			{{"$lookup", bson.D{{"from", "size"}, {"localField", "size_id"}, {"foreignField", "_id"}, {"as", "size"}}}},
			{{"$addFields", bson.D{{"size", bson.D{{"$arrayElemAt", bson.A{"$size.name", 0}}}}}}},
			{{"$lookup", bson.D{{"from", "payment_services"}, {"localField", "payment_method"}, {"foreignField", "_id"}, {"as", "payment_mode"}}}},
			{{"$addFields", bson.D{{"payment_mode", bson.D{{"$arrayElemAt", bson.A{"$payment_mode", 0}}}}}}},
			{{"$addFields", bson.D{
				{"buyer_info", bson.D{{"$arrayElemAt", bson.A{"$buyer_info", 0}}}},
				{"seller_info", bson.D{{"$arrayElemAt", bson.A{"$seller_info", 0}}}},
			}}},
			{{"$lookup", bson.D{{"from", "cloths"}, {"localField", "item_id"}, {"foreignField", "_id"}, {"as", "item"}}}},
			{{"$addFields", bson.D{
				{"item_name", bson.D{{"$arrayElemAt", bson.A{"$item.name", 0}}}},
				{"item_fabric", bson.D{{"$arrayElemAt", bson.A{"$item.fabric", 0}}}},
			}}},
			{{"$lookup", bson.D{{"from", "fabric"}, {"localField", "item_fabric"}, {"foreignField", "_id"}, {"as", "item_fabric"}}}},
			{{"$addFields", bson.D{{"item_fabric", bson.D{{"$arrayElemAt", bson.A{"$item_fabric.name", 0}}}}}}},
			{{"$lookup", bson.D{{"from", "employee"}, {"localField", "order_verified_by"}, {"foreignField", "_id"}, {"as", "order_verified_by"}}}},
			{{"$addFields", bson.D{{"order_verified_by", bson.D{{"$arrayElemAt", bson.A{"$order_verified_by", 0}}}}}}},
			bson.D{{"$sort", bson.D{{"order_date", -1}}}},
		}

		// Full pipeline with $facet
		pipeline := mongo.Pipeline{
			{{"$facet", bson.D{
				{"orders", append(lookupStages,
					bson.D{{"$project", bson.D{
						{"seller_info.primaryPhone", 1},
						{"seller_info.currency", 1},
						{"seller_info.countryName", 1},
						{"buyer_info.primaryPhone", 1},
						{"buyer_info.currency", 1},
						{"buyer_info.countryName", 1},
						{"payment_mode", 1},
						{"order_date", 1},
						{"color", 1},
						{"total_price", 1},
						{"delivery_status", 1},
						{"size", 1},
						{"quantity", 1},
						{"discount", 1},
						{"currency", 1},
						{"item_fabric", 1},
						{"item_name", 1},
						{"item_id", 1},
						{"delivery_date", 1},
						{"deliver_by", 1},
						{"category", 1},
						{"sub_category", 1},
						{"order_number", 1},
						{"order_verified", 1},
						{"isModified", 1},
						{"order_verified_by", 1},
					}}},
				)},
				{"summary", mongo.Pipeline{
					{{"$group", bson.D{
						{"_id", nil},
						{"total_orders", bson.D{{"$sum", 1}}},
						{"delivered_orders", bson.D{{"$sum", bson.D{{"$cond", bson.A{bson.D{{"$eq", bson.A{"$delivery_status", "Delivered"}}}, 1, 0}}}}}},
						{"pending_orders", bson.D{{"$sum", bson.D{{"$cond", bson.A{bson.D{{"$eq", bson.A{"$delivery_status", "Pending"}}}, 1, 0}}}}}},
					}}},
				}},
			}}},
			{{"$addFields", bson.D{
				{"summary", bson.D{{"$arrayElemAt", bson.A{"$summary", 0}}}},
			}}},
		}

		cursor, err := cartCollection.Aggregate(ctx, pipeline)
		if err != nil {
			http.Error(w, "Aggregation error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		var finalResult []bson.M
		if err := cursor.All(ctx, &finalResult); err != nil {
			http.Error(w, "Error decoding result: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return single object from array
		if len(finalResult) > 0 {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(finalResult[0])
		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(bson.M{
				"orders":  []bson.M{},
				"summary": bson.M{"total_orders": 0, "delivered_orders": 0, "pending_orders": 0},
			})
		}
	}
}

// GetDetailedCartItemsHandler fetches detailed cart items.
// It can fetch all orders if no 'user_id' is provided, or specific orders for a user if 'user_id' is provided.
// GetDetailedCartItemsHandler fetches detailed cart items.
func GetDetailedCartItemsHandler_v3(cartCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		userID := r.URL.Query().Get("user_id")
		status := r.URL.Query().Get("status")
		orderLocation := r.URL.Query().Get("order_location")
		startDateStr := r.URL.Query().Get("start_date")
		endDateStr := r.URL.Query().Get("end_date")
		dateStr := r.URL.Query().Get("date")

		matchStage := bson.D{}

		if userID != "" {
			objectID, err := primitive.ObjectIDFromHex(userID)
			if err != nil {
				http.Error(w, "Invalid user ID format. Please use a valid hex string.", http.StatusBadRequest)
				return
			}
			matchStage = append(matchStage, bson.E{Key: "user_id", Value: objectID})
		}

		// Date filters
		if dateStr != "" {
			orderDate, err := time.Parse("2006-01-02", dateStr)
			if err == nil {
				start := orderDate.UTC()
				end := start.Add(24 * time.Hour)
				matchStage = append(matchStage, bson.E{Key: "order_date", Value: bson.D{
					{"$gte", start},
					{"$lt", end},
				}})
			} else {
				http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
				return
			}
		} else if startDateStr != "" || endDateStr != "" {
			dateFilter := bson.D{}
			if startDateStr != "" {
				startDate, err := time.Parse("2006-01-02", startDateStr)
				if err != nil {
					http.Error(w, "Invalid start_date format. Use YYYY-MM-DD", http.StatusBadRequest)
					return
				}
				dateFilter = append(dateFilter, bson.E{Key: "$gte", Value: startDate.UTC()})
			}
			if endDateStr != "" {
				endDate, err := time.Parse("2006-01-02", endDateStr)
				if err != nil {
					http.Error(w, "Invalid end_date format. Use YYYY-MM-DD", http.StatusBadRequest)
					return
				}
				dateFilter = append(dateFilter, bson.E{Key: "$lte", Value: endDate.Add(24 * time.Hour).UTC()})
			}
			if len(dateFilter) > 0 {
				matchStage = append(matchStage, bson.E{Key: "order_date", Value: dateFilter})
			}
		}

		if status != "" {
			matchStage = append(matchStage, bson.E{Key: "delivery_status", Value: status})
		}
		if orderLocation != "" {
			matchStage = append(matchStage, bson.E{Key: "buyer_info.countryName", Value: orderLocation})
		}

		// --- Aggregation Pipeline ---
		pipeline := mongo.Pipeline{
			{{"$match", matchStage}},

			// 1. Lookups
			{{"$lookup", bson.D{{"from", "Mammalas_login"}, {"localField", "user_id"}, {"foreignField", "_id"}, {"as", "buyer_info"}}}},
			{{"$lookup", bson.D{{"from", "Mammalas_login"}, {"localField", "seller_id"}, {"foreignField", "_id"}, {"as", "seller_info"}}}},
			{{"$lookup", bson.D{{"from", "color"}, {"localField", "color_id"}, {"foreignField", "_id"}, {"as", "color_doc"}}}},
			{{"$lookup", bson.D{{"from", "size"}, {"localField", "size_id"}, {"foreignField", "_id"}, {"as", "size_doc"}}}},
			{{"$lookup", bson.D{{"from", "payment_services"}, {"localField", "payment_method"}, {"foreignField", "_id"}, {"as", "payment_mode_doc"}}}},
			{{"$lookup", bson.D{{"from", "cloths"}, {"localField", "item_id"}, {"foreignField", "_id"}, {"as", "item_doc"}}}},
			{{"$lookup", bson.D{{"from", "fabric"}, {"localField", "item_fabric"}, {"foreignField", "_id"}, {"as", "item_fabric_doc"}}}},
			{{"$lookup", bson.D{{"from", "employee"}, {"localField", "order_verified_by"}, {"foreignField", "_id"}, {"as", "order_verified_by_doc"}}}},

			// 2. Reshape and Flatten
			{
				{"$addFields", bson.D{
					{"buyer_info", bson.D{{"$arrayElemAt", bson.A{"$buyer_info", 0}}}},
					{"seller_info", bson.D{{"$arrayElemAt", bson.A{"$seller_info", 0}}}},
					{"color", bson.D{{"$arrayElemAt", bson.A{"$color_doc.name", 0}}}},
					{"size", bson.D{{"$arrayElemAt", bson.A{"$size_doc.name", 0}}}},
					{"payment_mode", bson.D{{"$arrayElemAt", bson.A{"$payment_mode_doc", 0}}}},
					{"item_name", bson.D{{"$arrayElemAt", bson.A{"$item_doc.name", 0}}}},
					{"item_fabric", bson.D{{"$arrayElemAt", bson.A{"$item_fabric_doc.name", 0}}}},
					{"order_verified_by", bson.D{{"$arrayElemAt", bson.A{"$order_verified_by_doc", 0}}}},

					// Map your Cart struct fields directly to output fields.
					// We are using the names that will be in the final $project stage.
					{"quantity", "$quantity"},
					{"total_price", "$total_price"},
					{"discount", "$discount"},
					{"delivery_status", "$delivery_status"},
					{"order_date", "$order_date"},
					{"currency", "$currency"},
					{"category", "$category"},
					{"sub_category", "$sub_category"},
					{"isModified", "$isModified"},
					{"order_number", "$order_number"},
					{"order_verified", "$order_verified"},

					// Handle delivery_date mapping from your Cart.DeliveredOn
					// Use a conditional to check if DeliveredOn is not the zero time.
					{"delivery_date", bson.D{{"$cond", bson.D{
						{"if", bson.D{{"$ne", bson.A{"$delivered_on", primitive.DateTime(0)}}}}, // Check if delivered_on is not the zero time.
						{"then", "$delivered_on"},
						{"else", nil}, // If it's zero time, use null (nil in Go)
					}}}},
					// You might also want to map DeliveredBy if it's meaningful as an ID, or perform another lookup for employee name.
					// For now, mapping the ID from your struct directly.
					{"deliver_by", "$delivered_by"},
				}},
			},

			// 3. Project the final fields to shape the output
			{
				{"$project", bson.D{
					{"_id", 1},
					{"item_id", 1}, // Map from your Cart struct's ItemID
					{"color", 1},
					{"quantity", 1},
					{"total_price", 1},
					{"discount", 1},
					{"delivery_status", 1},
					{"order_date", 1},
					{"order_number", 1},
					{"order_verified", 1},
					{"currency", 1},
					{"category", 1},
					{"sub_category", 1},
					{"isModified", 1},
					{"delivery_date", 1}, // Use the reshaped field
					{"deliver_by", 1},    // Use the reshaped field

					{"buyer_info", bson.D{
						{"countryName", "$buyer_info.countryName"},
						{"currency", "$buyer_info.currency"},
						{"primaryPhone", "$buyer_info.primaryPhone"},
					}},
					{"seller_info", bson.D{
						{"countryName", "$seller_info.countryName"},
						{"currency", "$seller_info.currency"},
						{"primaryPhone", "$seller_info.primaryPhone"},
					}},
					{"payment_mode", 1},
					{"size", 1},
					{"item_name", 1},
					{"item_fabric", 1},
					{"order_verified_by", bson.D{
						{"_id", "$order_verified_by._id"},
						{"emp_id", "$order_verified_by.emp_id"},
						{"name", "$order_verified_by.name"},
					}},
				}},
			},

			// 4. Sort
			{{"$sort", bson.D{{"order_date", -1}}}},

			// 5. Facet
			{
				{"$facet", bson.D{
					{"orders", mongo.Pipeline{}},
					{"summary", mongo.Pipeline{
						{{"$group", bson.D{
							{"_id", nil},
							{"total_orders", bson.D{{"$sum", 1}}},
							{"delivered_orders", bson.D{{"$sum", bson.D{{"$cond", bson.A{bson.D{{"$eq", bson.A{"$delivery_status", "Delivered"}}}, 1, 0}}}}}},
							{"pending_orders", bson.D{{"$sum", bson.D{{"$cond", bson.A{bson.D{{"$eq", bson.A{"$delivery_status", "Pending"}}}, 1, 0}}}}}},
						}}},
					}},
				}},
			},
			// 6. Flatten Summary
			{
				{"$addFields", bson.D{
					{"summary", bson.D{{"$arrayElemAt", bson.A{"$summary", 0}}}},
				}},
			},
		}

		cursor, err := cartCollection.Aggregate(ctx, pipeline)
		if err != nil {
			http.Error(w, "Error processing request: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		var resultWithFacet []bson.M
		if err := cursor.All(ctx, &resultWithFacet); err != nil {
			http.Error(w, "Error processing results: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if len(resultWithFacet) > 0 {
			var cartResponse CartResponse
			bsonBytes, err := bson.Marshal(resultWithFacet[0])
			if err != nil {
				http.Error(w, "Internal server error during response formatting", http.StatusInternalServerError)
				return
			}
			err = bson.Unmarshal(bsonBytes, &cartResponse)
			if err != nil {
				// log.Printf("Error unmarshalling BSON into CartResponse: %v, data: %+v", err, resultWithFacet[0])
				http.Error(w, "Internal server error during response formatting", http.StatusInternalServerError)
				return
			}

			if err := json.NewEncoder(w).Encode(cartResponse); err != nil {
				http.Error(w, "Error writing response", http.StatusInternalServerError)
			}
		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(bson.M{
				"orders":  []bson.M{},
				"summary": bson.M{"total_orders": 0, "delivered_orders": 0, "pending_orders": 0},
			})
		}
	}
}

// GetOrders handler to handle the API request
func GetCartFromDateToDate(w http.ResponseWriter, r *http.Request) {
	// Retrieve start and end date from query parameters
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	iso_start_Date, err := convertToISO8601(startDateStr)
	iso_end_date, err := convertToISO8601(endDateStr)
	startDate, err := time.Parse(time.RFC3339, iso_start_Date)
	if err != nil {
		log.Fatalf("Failed to parse date: %v", err)
	}

	endDate, err := time.Parse(time.RFC3339, iso_end_date)
	if err != nil {
		log.Fatalf("Failed to parse date: %v", err)
	}

	println("Start Date", iso_start_Date)
	println("End Date", iso_end_date)
	// MongoDB aggregation pipeline
	pipeline := bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"order_date",
						bson.D{
							{"$gte", startDate},
							{"$lt", endDate},
						},
					},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "Mammalas_login"},
					{"localField", "user_id"},
					{"foreignField", "_id"},
					{"as", "buyer_info"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "Mammalas_login"},
					{"localField", "seller_info"},
					{"foreignField", "_id"},
					{"as", "seller_info"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "color"},
					{"localField", "color_id"},
					{"foreignField", "_id"},
					{"as", "color"},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"color",
						bson.D{
							{"$arrayElemAt",
								bson.A{
									"$color.name",
									0,
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "size"},
					{"localField", "size_id"},
					{"foreignField", "_id"},
					{"as", "size"},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"size",
						bson.D{
							{"$arrayElemAt",
								bson.A{
									"$size.name",
									0,
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "payment_services"},
					{"localField", "payement_method"},
					{"foreignField", "_id"},
					{"as", "payment_mode"},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"payment_mode",
						bson.D{
							{"$arrayElemAt",
								bson.A{
									"$payment_mode.name",
									0,
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"buyer_info",
						bson.D{
							{"$arrayElemAt",
								bson.A{
									"$buyer_info",
									0,
								},
							},
						},
					},
					{"seller_info",
						bson.D{
							{"$arrayElemAt",
								bson.A{
									"$seller_info",
									0,
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "cloths"},
					{"localField", "item_id"},
					{"foreignField", "_id"},
					{"as", "item"},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"item_name",
						bson.D{
							{"$arrayElemAt",
								bson.A{
									"$item.name",
									0,
								},
							},
						},
					},
					{"item_fabric",
						bson.D{
							{"$arrayElemAt",
								bson.A{
									"$item.fabric",
									0,
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "fabric"},
					{"localField", "item_fabric"},
					{"foreignField", "_id"},
					{"as", "item_fabric"},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"item_fabric",
						bson.D{
							{"$arrayElemAt",
								bson.A{
									"$item_fabric.name",
									0,
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "employee"},
					{"localField", "order_verified_by"},
					{"foreignField", "_id"},
					{"as", "order_verified_by"},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"order_verified_by",
						bson.D{
							{"$arrayElemAt",
								bson.A{
									"$order_verified_by",
									0,
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"seller_info.primaryPhone", 1},
					{"seller_info.currency", 1},
					{"seller_info.countryName", 1},
					{"buyer_info.primaryPhone", 1},
					{"buyer_info.currency", 1},
					{"buyer_info.countryName", 1},
					{"payment_mode", 1},
					{"orderDate", 1},
					{"color", 1},
					{"total_price", 1},
					{"delivery_status", 1},
					{"size", 1},
					{"quantity", 1},
					{"discount", 1},
					{"currency", 1},
					{"item_fabric", 1},
					{"item_name", 1},
					{"item_id", 1},
					{"order_verified_by", 1},
					{"order_date", 1},
				},
			},
		},
		bson.D{
			{"$facet",
				bson.D{
					{"orders",
						bson.A{
							bson.D{
								{"$project",
									bson.D{
										{"seller_info.primaryPhone", 1},
										{"seller_info.currency", 1},
										{"seller_info.countryName", 1},
										{"buyer_info.primaryPhone", 1},
										{"buyer_info.currency", 1},
										{"buyer_info.countryName", 1},
										{"payment_mode", 1},
										{"orderDate", 1},
										{"color", 1},
										{"total_price", 1},
										{"delivery_status", 1},
										{"size", 1},
										{"quantity", 1},
										{"discount", 1},
										{"currency", 1},
										{"item_fabric", 1},
										{"item_name", 1},
										{"item_id", 1},
										{"order_verified_by", 1},
										{"order_date", 1},
									},
								},
							},
						},
					},
					{"summary",
						bson.A{
							bson.D{
								{"$group",
									bson.D{
										{"_id", primitive.Null{}},
										{"total_orders", bson.D{{"$sum", 1}}},
										{"delivered_orders",
											bson.D{
												{"$sum",
													bson.D{
														{"$cond",
															bson.A{
																bson.D{
																	{"$eq",
																		bson.A{
																			"$delivery_status",
																			"Delivered",
																		},
																	},
																},
																1,
																0,
															},
														},
													},
												},
											},
										},
										{"pending_orders",
											bson.D{
												{"$sum",
													bson.D{
														{"$cond",
															bson.A{
																bson.D{
																	{"$eq",
																		bson.A{
																			"$delivery_status",
																			"Pending",
																		},
																	},
																},
																1,
																0,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"summary",
						bson.D{
							{"$arrayElemAt",
								bson.A{
									"$summary",
									0,
								},
							},
						},
					},
				},
			},
		},
	}

	// Execute the aggregation query
	cursor, err := docking.PakTradeDb.Collection("cart_mammals").Aggregate(context.TODO(), pipeline)
	if err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	// Decode the result
	var orders []bson.M
	if err := cursor.All(context.TODO(), &orders); err != nil {
		http.Error(w, "Failed to decode orders", http.StatusInternalServerError)
		return
	}

	// Set the response header to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode the result to JSON and send it as the response
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func AddToCartHandler(cartCollection *mongo.Collection, itemCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload AddToCartPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, "Invalid JSON payload: "+err.Error(), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		orderNumber, err := getNextOrderNumber()
		for _, order := range payload.Orders {
			// ✅ Step 1: Get item price from itemCollection
			var item struct {
				Price float64 `bson:"price"`
			}

			itemObjID, err := primitive.ObjectIDFromHex(order.ItemID.Hex())
			if err != nil {
				http.Error(w, "Invalid item ID: "+err.Error(), http.StatusBadRequest)
				return
			}

			err = itemCollection.FindOne(ctx, bson.M{"_id": itemObjID}).Decode(&item)
			if err != nil {
				http.Error(w, "Item not found: "+err.Error(), http.StatusNotFound)
				return
			}

			// ✅ Step 2: Calculate total price
			totalPrice := item.Price * float64(order.Quantity)
			println("ORDER QTY: ", order.Quantity)
			// ✅ Step 3: Build cart document
			doc := bson.M{
				"user_id":           order.UserID,
				"item_id":           order.ItemID,
				"size_id":           order.SizeID,
				"color_id":          order.ColorID,
				"order_date":        order.OrderDate,
				"discount":          order.Discount,
				"currency":          order.Currency,
				"seller_id":         order.SellerId,
				"payment_method":    order.PaymentMethod,
				"total_price":       totalPrice,
				"category":          order.Category,
				"sub_category":      order.SubCategory,
				"quantity":          order.Quantity,
				"delivery_status":   "Pending",
				"order_number":      orderNumber,
				"isModified":        false,
				"delivery_date":     nil,
				"deliver_by":        nil,
				"order_verified_by": nil,
				"order_verified":    false,
			}

			_, err = cartCollection.InsertOne(ctx, doc)
			if err != nil {
				log.Println("Error inserting cart:", err)
				http.Error(w, "Failed to add item to cart", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Items added to cart successfully",
			"status":  "201",
		})
	}
}

func UpdateOrderPartial(ctx context.Context, orderID string, updates map[string]interface{}, modifiedBy string, orderCol, auditCol *mongo.Collection) error {
	oid, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return fmt.Errorf("invalid order ID: %w", err)
	}
	println("Order ID: ", orderID)
	// Get current order
	var current bson.M

	if err := orderCol.FindOne(ctx, bson.M{"_id": oid}).Decode(&current); err != nil {
		println("IN ERROR")
		return fmt.Errorf("order not found: %w", err)
	}

	// Prepare diff
	oldValues := make(map[string]interface{})
	newValues := make(map[string]interface{})

	for key, newVal := range updates {
		oldVal, exists := current[key]
		if !exists || !reflect.DeepEqual(oldVal, newVal) {
			oldValues[key] = oldVal
			newValues[key] = newVal
		}
	}

	if len(newValues) == 0 {
		return errors.New("no changes detected")
	}
	uid, _ := primitive.ObjectIDFromHex(modifiedBy)
	// Add audit log
	audit := bson.M{
		"orderId":    oid,
		"modify_by":  uid,
		"updated_at": time.Now(),
		"changes": bson.M{
			"old": oldValues,
			"new": newValues,
		},
	}
	if _, err := auditCol.InsertOne(ctx, audit); err != nil {
		return fmt.Errorf("failed to save audit log: %w", err)
	}

	// Set system fields
	updates["modify_by"] = uid
	updates["updated_at"] = time.Now()
	updates["isModified"] = true

	// Apply changes
	_, err = orderCol.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	return nil
}

// Handler
func UpdateOrderHandler(orderCollection, auditCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut && r.Method != http.MethodPatch {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 1. Get order ID from query params
		orderID := r.URL.Query().Get("cart_id")
		if orderID == "" {
			http.Error(w, "Missing order ID", http.StatusBadRequest)
			return
		}

		// 2. Get 'ModifiedBy' from header (you can change this logic)
		modifiedBy := r.Header.Get("emp_id")
		if modifiedBy == "" {
			http.Error(w, "Missing Employee-User-ID header", http.StatusBadRequest)
			return
		}

		// Decode only the changed fields
		var updatedFields map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updatedFields); err != nil {
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		// 4. Call the update logic
		err := UpdateOrderPartial(r.Context(), orderID, updatedFields, modifiedBy, orderCollection, auditCollection)
		if err != nil {
			http.Error(w, "Failed to update orde: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 5. Respond success
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Order updated and audit logged successfully",
		})
	}
}

func GetOrderSnapshotsHandler(auditCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := r.URL.Query().Get("id")
		if orderID == "" {
			http.Error(w, "Missing 'id' query param", http.StatusBadRequest)
			return
		}

		oid, err := primitive.ObjectIDFromHex(orderID)
		if err != nil {
			http.Error(w, "Invalid order ID", http.StatusBadRequest)
			return
		}

		findOptions := options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}})

		filter := bson.M{"orderId": oid}
		cursor, err := auditCollection.Find(r.Context(), filter, findOptions)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(r.Context())

		var snapshots []bson.M
		if err := cursor.All(r.Context(), &snapshots); err != nil {
			http.Error(w, "Failed to read data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var flattened []FlatChange
		for _, doc := range snapshots {
			flattened = append(flattened, transformAuditFlat(doc)...)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(flattened)
	}
}

func transformAuditFlat(doc bson.M) []FlatChange {
	orderID := doc["orderId"].(primitive.ObjectID).Hex()
	modifiedBy := doc["modify_by"].(primitive.ObjectID).Hex()
	modifiedAt := doc["updated_at"].(primitive.DateTime)

	oldMap, _ := doc["changes"].(bson.M)["old"].(bson.M)
	newMap, _ := doc["changes"].(bson.M)["new"].(bson.M)

	var result []FlatChange
	for key, oldVal := range oldMap {
		if newVal, ok := newMap[key]; ok {
			entry := FlatChange{
				OrderID:    orderID,
				ModifiedBy: modifiedBy,
				ModifiedAt: modifiedAt.Time(),
			}
			entry.Changes.Field = key
			entry.Changes.From = oldVal
			entry.Changes.To = newVal
			result = append(result, entry)
		}
	}

	return result
}

func getNextOrderNumber() (int, error) {
	filter := bson.M{"_id": "order_number_1"}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result struct {
		Seq int `bson:"seq"`
	}

	err := docking.PakTradeDb.Collection("cart_counter").FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		opts,
	).Decode(&result)

	if err != nil {
		return 0, err
	}

	return result.Seq, nil
}

func convertToISO8601(dateStr string) (string, error) {
	// Parse the date from "YYYY-MM-DD" format
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", err
	}

	// Convert to ISO 8601 format (RFC3339)
	return parsedDate.Format(time.RFC3339), nil
}
