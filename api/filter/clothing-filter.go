package clothingFilter

import (
	"context"
	"encoding/json"
	"net/http"

	docking "pak-trade-go/Docking"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type filter_list struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}

type Filter struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty"`
	FilterName string               `bson:"filter_name" json:"filter_name"`
	FilterList []primitive.ObjectID `bson:"filter_list" json: "filter_list"`
}

func getFilters(client *mongo.Client) ([]Filter, error) {
	collection := docking.PakTradeDb.Collection("clothing_filter")
	ctx := context.TODO()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var filters []Filter

	for cursor.Next(ctx) {
		var filter Filter
		if err := cursor.Decode(&filter); err != nil {
			return nil, err
		}
		filters = append(filters, filter)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return filters, nil

}

func FiltersHandler(w http.ResponseWriter, r *http.Request) {
	clientOptions := options.Client().ApplyURI("mongodb+srv://developer001:IAmMuslim@cluster0.qeqntol.mongodb.net/test")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.TODO())

	filters, err := getFilters(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filters)

}
