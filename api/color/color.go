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
	Name   string             `json:"name,omitempty"`
}
type respone_struct struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    []Color1 `json:"data"`
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
	var responce respone_struct
	if results != nil {
		responce.Status = http.StatusOK
		responce.Message = "success"
		responce.Data = results
	} else {
		responce.Status = http.StatusBadRequest
		responce.Message = "declined"
	}
	output, err := json.MarshalIndent(responce, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "%s\n", output)

}
