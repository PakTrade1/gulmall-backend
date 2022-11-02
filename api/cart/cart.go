package cart

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
func Cart_insertone(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var data []CartMammals
	json.NewDecoder(req.Body).Decode(&data)
	//fmt.Fprintf(w, "%s\n", data)

	coll := docking.PakTradeDb.Collection("cart_mammals")

	// insert a user

	_, err3 := coll.InsertOne(context.TODO(), data)
	if err3 != nil {
		fmt.Fprintf(w, "%s\n", err3)
	}
	output, err2 := json.Marshal(data)
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Fprintf(w, "%s\n", output)
	// output, err := json.MarshalIndent(insertResult.InsertedID, "", "    ")
	// if err != nil {
	// 	panic(err)

	// }
	//fmt.Fprintf(w, "%s\n",)

}
