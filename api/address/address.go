package address

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	docking "pak-trade-go/Docking"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type shipping_mammal_id struct {
	Mammal_id string `json:"mammal_id"`
}

type shipping_table struct {
	Address []struct {
		Name        string `json:"name"`
		Line_1      string `json:"address_line_1"`
		Line_2      string `json:"address_line_2"`
		Subrub      string `json:"subrub"`
		City        string `json:"city"`
		Postal_code string `json:"postal_code"`
		Uid         string `json:"uid"`
	} `json:"Address"`
	Mammal_id string `json:"mammal_id"`
}

type respone_struct struct {
	Status  int              `json:"status"`
	Message string           `json:"message"`
	Data    []shipping_table `json:"data"`
}

func Get_shipping_address_with_mammal_id(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var shipping_req shipping_mammal_id
	err := json.NewDecoder(req.Body).Decode(&shipping_req)
	if err != nil {
		panic(err)
	}

	coll := docking.PakTradeDb.Collection("Mammals_address")

	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver
	ctx := context.TODO()

	mongo_query := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "mammal_id", Value: shipping_req.Mammal_id}}}},
	}
	// Open an aggregation cursor
	cursor, err := coll.Aggregate(ctx, mongo_query)
	if err != nil {
		log.Fatal(err)
	}
	var results respone_struct
	var resp1 []shipping_table
	for cursor.Next(context.TODO()) {

		var xy shipping_table
		cursor.Decode(&xy)
		resp1 = append(resp1, xy)

	}
	if cursor != nil {
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

func Add_shipping_address(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	newUUID := (uuid.New()).String()

	var shipping_req shipping_table
	err := json.NewDecoder(req.Body).Decode(&shipping_req)
	if err != nil {
		panic(err)
	}

	coll := docking.PakTradeDb.Collection("Mammals_address")
	var result shipping_table
	filter := bson.M{"mammal_id": shipping_req.Mammal_id}
	err1 := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err1 != nil {
		fmt.Println("errror retrieving user userid : " + result.Mammal_id)
	}
	fmt.Print(shipping_req.Mammal_id)
	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver
	if result.Mammal_id != "" {
		ctx := context.TODO()
		// Open an aggregation cursor
		cursor, err := coll.UpdateOne(ctx,
			bson.M{"mammal_id": shipping_req.Mammal_id},
			bson.D{
				{Key: "$push", Value: bson.M{
					"Address": bson.M{
						"name":        shipping_req.Address[0].Name,
						"line_1":      shipping_req.Address[0].Line_1,
						"line_2":      shipping_req.Address[0].Line_2,
						"subrub":      shipping_req.Address[0].Subrub,
						"city":        shipping_req.Address[0].City,
						"postal_code": shipping_req.Address[0].Postal_code,
						"uid":         newUUID,
					},
				}},
			},
		)
		if err != nil {
			fmt.Fprintf(w, "%s\n", err)
		}
		var results respone_struct
		var resp1 []shipping_table
		if cursor != nil {
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
	} else {
		ctx := context.TODO()

		mongo_query := bson.M{
			"Address": bson.A{
				bson.M{
					"name":        shipping_req.Address[0].Name,
					"line_1":      shipping_req.Address[0].Line_1,
					"line_2":      shipping_req.Address[0].Line_2,
					"subrub":      shipping_req.Address[0].Subrub,
					"city":        shipping_req.Address[0].City,
					"postal_code": shipping_req.Address[0].Postal_code,
					"uid":         newUUID,
				},
			},
			"mammal_id": shipping_req.Mammal_id,
		}
		// Open an aggregation cursor
		cursor, err := coll.InsertOne(ctx, mongo_query)
		if err != nil {
			fmt.Fprintf(w, "%s\n", err)
		}
		var results respone_struct
		var resp1 []shipping_table
		if cursor != nil {
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

}

// delete adress
type delete_address_uid struct {
	Mammal_id string `json:"mammal_id"`
	Uid       string `json:"uid"`
}

func Delete_shipping_address(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var shipping_req delete_address_uid
	err := json.NewDecoder(req.Body).Decode(&shipping_req)
	if err != nil {
		panic(err)
	}

	coll := docking.PakTradeDb.Collection("Mammals_address")
	var result shipping_table
	filter := bson.M{"mammal_id": shipping_req.Mammal_id}
	err1 := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err1 != nil {
		fmt.Println("errror retrieving user userid : " + result.Mammal_id)
	}

	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver
	if result.Mammal_id != "" {
		ctx := context.TODO()
		// Open an aggregation cursor
		cursor, err := coll.UpdateOne(ctx,
			bson.M{"mammal_id": shipping_req.Mammal_id},
			bson.D{
				{Key: "$pull", Value: bson.M{
					"Address": bson.M{
						"uid": shipping_req.Uid,
					},
				}},
			},
		)
		if err != nil {
			fmt.Fprintf(w, "%s\n", err)
		}
		var results respone_struct
		var resp1 []shipping_table
		if cursor != nil {
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
	} else {

		// Open an aggregation cursor

		var results respone_struct
		var resp1 []shipping_table

		results.Message = "decline"
		results.Data = resp1
		output, err := json.MarshalIndent(results, "", "    ")
		if err != nil {
			panic(err)

		}

		fmt.Fprintf(w, "%s\n", output)
	}

}
