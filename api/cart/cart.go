package cart

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	docking "pak-trade-go/Docking"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartMammals struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Quantity  int                `json:"quantity,omitempty"`
	Price     int                `json:"price,omitempty"`
	Mammal_id primitive.ObjectID `json:"mammal_id,omitempty"`
	Item_id   primitive.ObjectID `json:"item_id,omitempty"`
	Color_id  primitive.ObjectID `json:"color_id,omitempty"`
	Size_id   primitive.ObjectID `json:"size_id,omitempty"`
}

func Cart_getall(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	coll := docking.PakTradeDb.Collection("cart_mammals")

	cursor, err := coll.Find(context.Background(), coll)
	if err != nil {
		panic(err)
	}

	var results []CartMammals
	for cursor.Next(context.TODO()) {
		var abc CartMammals
		cursor.Decode(&abc)
		results = append(results, abc)

	}

	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)

	}
	fmt.Fprintf(w, "%s\n", output)

}
func Mammals_insertone(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	var cart_init CartMammals
	err := json.NewDecoder(req.Body).Decode(&cart_init)
	if err != nil {
		panic(err)
	}
	// mongo query
	mongo_query := bson.M{
		"mammal_id": cart_init.Mammal_id,
		"item_id":   cart_init.Item_id,
		"color_id":  cart_init.Color_id,
		"size_id":   cart_init.Size_id,
		"quantity":  cart_init.Quantity,
		"price":     cart_init.Price,
	}

	coll := docking.PakTradeDb.Collection("cart_mammals")

	// // // insert a user

	inset, err3 := coll.InsertOne(context.TODO(), mongo_query)
	if err3 != nil {
		fmt.Fprintf(w, "%s\n", err3)
	}

	fmt.Fprintf(w, "%s\n", inset)

}
