package categories

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

var client = docking.AzureBloblogs()

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

type Categories struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name struct {
		En string `json:"en"`
		Ar string `json:"ar"`
	} `json:"name"`
}
type status_req struct {
	Status string `json:"status"`
}
type respone_struct struct {
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Data    []Categories `json:"data"`
}

func Get_all_categories(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	coll := docking.PakTradeDb.Collection("categories")

	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver

	// Open an aggregation cursor
	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	var results respone_struct
	var resp1 []Categories
	for cursor.Next(context.TODO()) {
		var xy Categories
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

type Cat_id_serch struct {
	Cat_id string `json:"cat_id"`
}
type respone_struct_sub_Cat struct {
	Status  int                      `json:"status"`
	Message string                   `json:"message"`
	Data    []sub_Categoies_selected `json:"data"`
}

type sub_Categoies_selected struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Cat_id string             `json:"cat_id"`
	Name   struct {
		En string `json:"en"`
		Ar string `json:"ar"`
	} `json:"name"`
}

func Sub_Categories_select_by_Cat_id(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var search1 Cat_id_serch
	err := json.NewDecoder(req.Body).Decode(&search1)
	if err != nil {
		panic(err)
	}
	coll := docking.PakTradeDb.Collection("sub_category")
	objectIDS, _ := primitive.ObjectIDFromHex(search1.Cat_id)

	//	var result sub_Categoies_selected
	//	filter := bson.M{"cat_id": objectIDS}

	cursor, err := coll.Find(context.Background(), bson.M{"cat_id": objectIDS})
	if err != nil {
		panic(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	var results respone_struct_sub_Cat
	var resp1 []sub_Categoies_selected
	for cursor.Next(context.TODO()) {
		var xy sub_Categoies_selected
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

type Chaild_cat_id_serch struct {
	Chaild_cat_id string `json:"chaild_cat_id"`
}
type Chaild_sub_Categoies_selected struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Sub_category_id string             `json:"sub_category_id"`
	Name            struct {
		En string `json:"en"`
		Ar string `json:"ar"`
	} `json:"name"`
}
type respone_struct_child_cat struct {
	Status  int                             `json:"status"`
	Message string                          `json:"message"`
	Data    []Chaild_sub_Categoies_selected `json:"data"`
}

func Child_Categories_select_by__sub_Cat_id(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var search1 Chaild_cat_id_serch
	err := json.NewDecoder(req.Body).Decode(&search1)
	if err != nil {
		panic(err)
	}
	coll := docking.PakTradeDb.Collection("sub_category_child")
	objectIDS, _ := primitive.ObjectIDFromHex(search1.Chaild_cat_id)

	//	var result sub_Categoies_selected
	//	filter := bson.M{"cat_id": objectIDS}

	cursor, err := coll.Find(context.Background(), bson.M{"sub_category_id": objectIDS})
	if err != nil {
		panic(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	var results respone_struct_child_cat
	var resp1 []Chaild_sub_Categoies_selected
	for cursor.Next(context.TODO()) {
		var xy Chaild_sub_Categoies_selected
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