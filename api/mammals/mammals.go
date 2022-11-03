package mammals

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	docking "pak-trade-go/Docking"

	"go.mongodb.org/mongo-driver/bson"
)

type Mammals_user struct {
	//ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name struct {
		Firt_name string `json:"firt_name" bson:"firt_name"`
		Last_name string `json:"last_name" bson:"last_name"`
	} `json:"name" bson:"name"`
	Address []struct {
		Home_address struct {
			Address  string `json:"address" bson:"address"`
			Country  string `json:"country"  bson:"country"`
			City     string `json:"city"  bson:"city"`
			Province string `json:"province" bson:"province"`
			Zip_code int    `json:"zip_code" bson:"zip_code"`
		} `json:"home_address,omitempty" bson:"home_address,omitempty"`
	} `json:"address"`
	Email   string `json:"email" bson:"email"`
	Phone   string `json:"phone" bson:"phone"`
	Profile string `json:"profile" bson:"profile"`
}

func Mammals_getall(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	coll := docking.PakTradeDb.Collection("mammals")

	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver
	// get all user in  cursor
	cursor, err := coll.Find(context.Background(), coll)
	if err != nil {
		panic(err)
	}

	var results []Mammals_user
	for cursor.Next(context.TODO()) {
		var abc Mammals_user
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

	var strcutinit Mammals_user
	err := json.NewDecoder(req.Body).Decode(&strcutinit)
	if err != nil {
		panic(err)
	}

	insertdat := bson.M{
		"name": bson.M{
			"firt_name": strcutinit.Name.Firt_name,
			"last_name": strcutinit.Name.Last_name,
		},
		"address": bson.A{
			bson.M{
				"home_address": bson.M{
					"address":  strcutinit.Address[0].Home_address.Address,
					"country":  strcutinit.Address[0].Home_address.Country,
					"city":     strcutinit.Address[0].Home_address.City,
					"province": strcutinit.Address[0].Home_address.Province,
					"zip_code": strcutinit.Address[0].Home_address.Zip_code,
				},
			},
		},
		"email":   strcutinit.Email,
		"phone":   strcutinit.Phone,
		"profile": strcutinit.Profile,
	}

	//fmt.Print(body)
	coll := docking.PakTradeDb.Collection("mammals")

	// // // insert a user

	inset, err3 := coll.InsertOne(context.TODO(), insertdat)
	if err3 != nil {
		fmt.Fprintf(w, "%s\n", err3)
	}

	fmt.Fprintf(w, "%s\n", inset)

}
