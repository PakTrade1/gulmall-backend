package items

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	docking "pak-trade-go/Docking"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemType struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Price int32              `json:"price,omitempty`

	Name struct {
		En string `json:"en,omitempty"`
		Ar string `json:"ar,omitempty"`
	} `json:"name,omitempty"`

	Feature []struct {
		//Low_quility  []string `json:"low_quility,omitempty"`
		//High_quility []string `json:"high_quility,omitempty"`
		Name struct {
			En string `json:"en,omitempty"`
			Ar string `json:"ar,omitempty"`
		} `json:"name,omitempty"`
	} `json:"feature,omitempty"`

	Available_size []string `json:"available_size,omitempty"`

	Images []struct {
		Low_quility  []string `json:"low_quility,omitempty"`
		High_quility []string `json:"high_quility,omitempty"`
	} `json:"images,omitempty"`

	Available_color []struct {
		ID     primitive.ObjectID `bson:"_id,omitempty"`
		CSSHex string             `json:"cssHex,omitempty"`
		Name   struct {
			En string `json:"en,omitempty"`
			Ar string `json:"ar,omitempty"`
		} `json:"name,omitempty"`
	} `json:"available_color,omitempty"`
}

func Items(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	coll := docking.PakTradeDb.Collection("cloths")

	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver
	ctx := context.TODO()

	// Open an aggregation cursor
	cursor, err := coll.Aggregate(ctx, bson.A{
		bson.D{
			{Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$available_color"},
					{Key: "includeArrayIndex", Value: "index_1"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "color"},
					{Key: "localField", Value: "available_color"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "result"},
				},
			},
		},
		bson.D{{Key: "$set", Value: bson.D{{Key: "result", Value: bson.D{{Key: "$first", Value: "$result"}}}}}},
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id",
						Value: bson.D{
							{Key: "name", Value: "$name"},
							{Key: "feature", Value: "$feature"},
							{Key: "sizes", Value: "$available_size"},
							{Key: "images", Value: "$images"},
							{Key: "_id", Value: "$_id"},
							{Key: "price", Value: "$price"},
						},
					},
					{Key: "colors", Value: bson.D{{Key: "$push", Value: "$result"}}},
				},
			},
		},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: "$_id._id"},
					{Key: "name", Value: "$_id.name"},
					{Key: "feature", Value: "$_id.feature"},
					{Key: "available_size", Value: "$_id.sizes"},
					{Key: "images", Value: "$_id.images"},
					{Key: "price", Value: "$_id.price"},
					{Key: "available_color", Value: "$colors"},
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	var results []ItemType
	for cursor.Next(context.TODO()) {
		var abc ItemType
		cursor.Decode(&abc)
		results = append(results, abc)

	}

	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)

	}
	fmt.Fprintf(w, "%s\n", output)

}
