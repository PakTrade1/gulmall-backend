package weight

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	docking "pak-trade-go/Docking"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type weight_get struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name struct {
		Short struct {
			En string `json:"en"`
			Ar string `json:"ar"`
		} `json:"short"`
		Full struct {
			En string `json:"en"`
			Ar string `json:"ar"`
		} `json:"full"`
	} `json:"name"`
}

type respone_struct struct {
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Data    []weight_get `json:"data"`
}

func Weight(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	coll := docking.PakTradeDb.Collection("weight")
	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}

	var results []weight_get
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	for cursor.Next(context.TODO()) {
		var abc weight_get
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
