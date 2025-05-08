package cart

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	docking "pak-trade-go/Docking"
	"reflect"
	"time"

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
	ID             primitive.ObjectID `bson:"_id" json:"cart_id"`
	ItemID         primitive.ObjectID `bson:"item_id" json:"item_id"`
	ColorID        primitive.ObjectID `bson:"color_id" json:"color_id"`
	Quantity       int                `bson:"quantity" json:"quantity"`
	TotalPrice     float64            `bson:"total_price" json:"total_price"`
	Discount       string             `bson:"discount" json:"discount"`
	PaymentMethod  primitive.ObjectID `bson:"payment_method" json:"payment_method"`
	UserID         primitive.ObjectID `bson:"user_id" json:"user_id"`
	SellerId       primitive.ObjectID `bson:"seller_id" json:"seller_id"`
	DeliveryStatus string             `bson:"delivery_status" json:"delivery_status"`
	OrderDate      time.Time          `bson:"orderDate" json:"order_placed_on"`
	SizeID         primitive.ObjectID `bson:"size_id" json:"size_id"`
	Currency       string             `bson:"currency" json:"currency"`
	Category       string             `bson:"category" json:"category"`
	SubCategory    string             `bson:"sub_category" json:"sub_category"`
	CreatedAt      time.Time          `bson:"created_at"`
	IsModified     primitive.ObjectID `bson:"isModified" json:"isModified"`
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

func GetAllCartHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collection := docking.PakTradeDb.Collection("cart_mammals") // replace with your actual cart collection

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "orderDate", Value: -1}})
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		http.Error(w, "Error fetching carts", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var carts []Cart
	if err := cursor.All(ctx, &carts); err != nil {
		http.Error(w, "Error decoding cart data", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(carts)
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
				"user_id":         order.UserID,
				"item_id":         order.ItemID,
				"size_id":         order.SizeID,
				"color_id":        order.ColorID,
				"order_date":      time.Now(),
				"discount":        order.Discount,
				"currency":        order.Currency,
				"seller_id":       order.SellerId,
				"payment_method":  order.PaymentMethod,
				"total_price":     totalPrice,
				"category":        order.Category,
				"sub_category":    order.SubCategory,
				"quantity":        order.Quantity,
				"delivery_status": "PENDING",
				"isModified":      false,
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
