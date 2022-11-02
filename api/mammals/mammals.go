package mammals

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	docking "pak-trade-go/Docking"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type ItemType struct {
// 	ID    primitive.ObjectID `bson:"_id,omitempty"`
// 	Price int32              `json:"price,omitempty`

// 	Name struct {
// 		En string `json:"en,omitempty"`
// 		Ar string `json:"ar,omitempty"`
// 	} `json:"name,omitempty"`

// 	Feature []struct {
// 		//Low_quility  []string `json:"low_quility,omitempty"`
// 		//High_quility []string `json:"high_quility,omitempty"`
// 		Name struct {
// 			En string `json:"en,omitempty"`
// 			Ar string `json:"ar,omitempty"`
// 		} `json:"name,omitempty"`
// 	} `json:"feature,omitempty"`

// 	Available_size []string `json:"available_size,omitempty"`

// 	Images []struct {
// 		Low_quility  []string `json:"low_quility,omitempty"`
// 		High_quility []string `json:"high_quility,omitempty"`
// 	} `json:"images,omitempty"`

//		Available_color []struct {
//			ID     primitive.ObjectID `bson:"_id,omitempty"`
//			CSSHex string             `json:"cssHex,omitempty"`
//			Name   struct {
//				En string `json:"en,omitempty"`
//				Ar string `json:"ar,omitempty"`
//			} `json:"name,omitempty"`
//		} `json:"available_color,omitempty"`
//	}

type Mammals_user struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// body, err := ioutil.ReadAll(req.Body)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	var data []Mammals_user
	json.NewDecoder(req.Body).Decode(&data)
	//fmt.Fprintf(w, "%s\n", data)

	coll := docking.PakTradeDb.Collection("mammals")

	// insert a user

	_, err3 := coll.InsertOne(context.TODO(), data)
	if err3 != nil {
		fmt.Fprintf(w, "%s\n", err3)
	}
	// output, err2 := json.Marshal(data)
	// if err2 != nil {
	// 	log.Fatal(err2)
	// }
	// fmt.Fprintf(w, "%s\n", output)
	// output, err := json.MarshalIndent(insertResult.InsertedID, "", "    ")
	// if err != nil {
	// 	panic(err)

	// }
	//fmt.Fprintf(w, "%s\n",)

}
