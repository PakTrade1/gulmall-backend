package items

import (
	"context"
	"encoding/json"
	"net/http"
	docking "pak-trade-go/Docking"
	"time"

	// if you're using Gin for routing
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetItemsByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Extract category ID from query param like /api/items?categoryId=xxx
	categoryIDStr := r.URL.Query().Get("categoryId")
	if categoryIDStr == "" {
		http.Error(w, "categoryId is required", http.StatusBadRequest)
		return
	}

	categoryID, err := primitive.ObjectIDFromHex(categoryIDStr)
	if err != nil {
		http.Error(w, "Invalid categoryId", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	itemParentColl := docking.PakTradeDb.Collection("product")

	pipeline := mongo.Pipeline{
		// Match category if needed
		{{"$match", bson.D{{"category", categoryID}}}},

		// Lookup category name
		{{"$lookup", bson.D{
			{"from", "categories"},
			{"localField", "category"},
			{"foreignField", "_id"},
			{"as", "category_name"},
		}}},

		// Lookup subCategory name
		{{"$lookup", bson.D{
			{"from", "sub_category"},
			{"localField", "subCategory"},
			{"foreignField", "_id"},
			{"as", "subCategory_name"},
		}}},

		// Lookup tier/plan name
		{{"$lookup", bson.D{
			{"from", "tier"},
			{"localField", "planId"},
			{"foreignField", "_id"},
			{"as", "tier_name"},
		}}},

		// Lookup variants from cloths
		{{"$lookup", bson.D{
			{"from", "cloths"},
			{"localField", "_id"},
			{"foreignField", "parentId"},
			{"as", "variants"},
		}}},

		// Add names using $arrayElemAt
		{{"$addFields", bson.D{
			{"category_name", bson.D{
				{"$arrayElemAt", bson.A{"$category_name.name", 0}},
			}},
			{"subCategory_name", bson.D{
				{"$arrayElemAt", bson.A{"$subCategory_name.name", 0}},
			}},
			{"tier_name", bson.D{
				{"$arrayElemAt", bson.A{"$tier_name.name", 0}},
			}},
		}}},

		// Optional: sort newest first
		{{"$sort", bson.D{{"creationTimestamp", -1}}}},
	}

	cursor, err := itemParentColl.Aggregate(ctx, pipeline)
	if err != nil {
		http.Error(w, "Database aggregation error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		http.Error(w, "Error decoding results: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"items": results,
	})
}
