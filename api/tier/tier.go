package tier

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

type tier_struct struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `json:"name"`
	Order   int                `json:"order"`
	IconURL string             `json:"iconUrl"`
}

type respone_struct struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    []tier_struct `json:"data"`
}

func Tier_get(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	coll := docking.PakTradeDb.Collection("tier")
	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}

	var results []tier_struct
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	for cursor.Next(context.TODO()) {
		var abc tier_struct
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

type Tire_serch struct {
	Tier_id string `json:"tierId"`
}

type Plan_sturct_resp struct {
	ID              string  `json:"_id"`
	Week            int     `json:"week"`
	AdDuration      int     `json:"ad_duration"`
	SpecialDuration int     `json:"special_duration"`
	Price           float64 `json:"price"`
	TierID          string  `json:"tierId"`
	Discount        string  `json:"discount"`
}

type serch_itme_result struct {
	Status  int                `json:"status"`
	Message string             `json:"message "`
	Data    []Plan_sturct_resp `json:"data"`
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
func Plan_select(w http.ResponseWriter, r *http.Request) {
	coll := docking.PakTradeDb.Collection("plans")

	var id_get Tire_serch
	err := json.NewDecoder(r.Body).Decode(&id_get)
	handleError(err)
	objectId, err := primitive.ObjectIDFromHex(id_get.Tier_id)
	handleError(err)

	cursor, err := coll.Find(context.TODO(), bson.D{{Key: "tierId", Value: objectId}})
	if err != nil {
		log.Fatal(err)
	}
	var results serch_itme_result

	// results := new(Plan_sturserch_itme_structct_resp)
	var resp1 []Plan_sturct_resp
	for cursor.Next(context.TODO()) {
		var xy Plan_sturct_resp
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
