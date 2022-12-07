package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	docking "pak-trade-go/Docking"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type payment_type struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	Name struct {
		En string `json:"en"`
		Ar string `json:"ar"`
	}
}

type respone_struct struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Data    []payment_type `json:"data"`
}

func Get_payment_method(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	coll := docking.PakTradeDb.Collection("payment_services")

	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver

	cursor, err := coll.Find(context.Background(), coll)
	if err != nil {
		panic(err)
	}

	results := new(respone_struct)
	var resp1 []payment_type
	for cursor.Next(context.TODO()) {
		var xy payment_type
		cursor.Decode(&xy)
		resp1 = append(resp1, xy)

	}
	if resp1 != nil {
		results.Status = http.StatusOK
		results.Message = "success"

	} else {
		results.Message = "decline"

	}
	results.Data = resp1
	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)

	}

	fmt.Fprintf(w, "%s\n", output)

}
