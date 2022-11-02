package colors

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	docking "pak-trade-go/Docking"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Name struct {
	Ar string `json:"ar,omitempty"`
	En string `json:"en,omitempty"`
}

type Color1 struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	CSSHex string             `json:"cssHex,omitempty"`
	Name   `json:"name,omitempty"`
}

func Color(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	coll := docking.PakTradeDb.Collection("color")
	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}

	var results []Color1
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	for cursor.Next(context.TODO()) {
		var abc Color1
		cursor.Decode(&abc)
		results = append(results, abc)

	}
	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "%s\n", output)

}
