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
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		Country     string `json:"country"`
		State       string `json:"state"`
	} `json:"Address"`
	Mammal_id string `json:"mammal_id"`
}
type Address struct {
	Name        string `json:"name"`
	Line_1      string `json:"address_line_1"`
	Line_2      string `json:"address_line_2"`
	Subrub      string `json:"subrub"`
	City        string `json:"city"`
	Postal_code string `json:"postal_code"`
	Uid         string `json:"uid"`
	Country     string `json:"country"`
	State       string `json:"state"`
}
type Address1 struct {
	Name        string `json:"name"`
	Line_1      string `json:"address_line_1"`
	Line_2      string `json:"address_line_2"`
	Subrub      string `json:"subrub"`
	City        string `json:"city"`
	Postal_code string `json:"postal_code"`
	Uid         string `json:"uid"`
	Country     string `json:"country"`
	State       string `json:"state"`
}

type respone_struct struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	// Data    []shipping_table `json:"data"`
}

type respone_struct1 struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    []Address1 `json:"data"`
}

func Get_shipping_address_with_mammal_id(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	userId := req.URL.Query().Get("userId")
	objID, err1 := primitive.ObjectIDFromHex(userId)
	if err1 != nil {
		log.Fatal(err1)
	}

	coll := docking.PakTradeDb.Collection("Mammals_address")

	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver
	ctx := context.TODO()

	mongo_query := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "mammal_id", Value: objID}}}},
		bson.D{{"$unwind", bson.D{{"path", "$Address"}}}},
		bson.D{
			{"$project",
				bson.D{
					{"line_1", "$Address.line_1"},
					{"line_2", "$address.line2"},
					{"postal_code", "$Address.postal_code"},
					{"state", "$Address.state"},
					{"country", "$Address.country"},
					{"name", "$Address.name"},
					{"subrub", "$Address.subrub"},
					{"city", "$Address.city"},
					{"uid", "$Address.uid"},
				},
			},
		},
	}
	// Open an aggregation cursor
	cursor, err := coll.Aggregate(ctx, mongo_query)
	if err != nil {
		log.Fatal(err)
	}
	var results respone_struct1
	var resp1 []Address1
	for cursor.Next(context.TODO()) {
		var xy Address1
		cursor.Decode(&xy)
		resp1 = append(resp1, xy)
	}
	if cursor != nil {
		results.Status = http.StatusOK
		results.Message = "success"
		results.Data = resp1

	} else {
		results.Message = "decline"

	}
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

	userId := req.URL.Query().Get("userId")
	objID, err1 := primitive.ObjectIDFromHex(userId)
	if err1 != nil {
		log.Fatal(err1)
	}

	var shipping_req Address
	err := json.NewDecoder(req.Body).Decode(&shipping_req)
	if err != nil {
		panic(err)
	}

	coll := docking.PakTradeDb.Collection("Mammals_address")
	var result shipping_table
	filter := bson.M{"mammal_id": objID}
	err2 := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err2 != nil {
		fmt.Println("errror retrieving user userid : " + result.Mammal_id)
	}
	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver
	if result.Mammal_id != "" {
		ctx := context.TODO()
		// Open an aggregation cursor
		cursor, err := coll.UpdateOne(ctx,
			bson.M{"mammal_id": objID},
			bson.D{
				{Key: "$push", Value: bson.M{
					"Address": bson.M{
						"name":        shipping_req.Name,
						"line_1":      shipping_req.Line_1,
						"line_2":      shipping_req.Line_2,
						"subrub":      shipping_req.Subrub,
						"city":        shipping_req.City,
						"postal_code": shipping_req.Postal_code,
						"uid":         newUUID,
						"country":     shipping_req.Country,
						"state":       shipping_req.State,
					},
				}},
			},
		)
		if err != nil {
			fmt.Fprintf(w, "%s\n", err)
		}
		var results respone_struct
		// var resp1 []shipping_table
		if cursor != nil {
			results.Status = http.StatusOK
			results.Message = "success"

		} else {
			results.Message = "decline"

		}

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
					"name":        shipping_req.Name,
					"line_1":      shipping_req.Line_1,
					"line_2":      shipping_req.Line_2,
					"subrub":      shipping_req.Subrub,
					"city":        shipping_req.City,
					"postal_code": shipping_req.Postal_code,
					"uid":         newUUID,
					"country":     shipping_req.Country,
					"state":       shipping_req.State,
				},
			},
			"mammal_id": objID,
		}
		// Open an aggregation cursor
		cursor, err := coll.InsertOne(ctx, mongo_query)
		if err != nil {
			fmt.Fprintf(w, "%s\n", err)
		}
		var results respone_struct
		// var resp1 []shipping_table
		if cursor != nil {
			results.Status = http.StatusOK
			results.Message = "success"

		} else {
			results.Message = "decline"

		}
		// results.Data = resp1
		output, err := json.MarshalIndent(results, "", "    ")
		if err != nil {
			panic(err)

		}

		fmt.Fprintf(w, "%s\n", output)
	}

}

// delete adress

func Delete_shipping_address(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	UserUuId := req.URL.Query().Get("userUuId")

	coll := docking.PakTradeDb.Collection("Mammals_address")

	ctx := context.TODO()
	filter := bson.D{{"Address", bson.D{{"$elemMatch", bson.D{{"uid", UserUuId}}}}}}

	cursor, err := coll.UpdateOne(ctx,
		filter,
		bson.D{
			{Key: "$pull", Value: bson.M{
				"Address": bson.M{
					"uid": UserUuId,
				},
			}},
		},
	)
	if err != nil {
		fmt.Fprintf(w, "%s\n", err)
	}
	var results respone_struct
	// var resp1 []shipping_table
	if cursor.ModifiedCount > 0 {
		results.Status = http.StatusOK
		results.Message = "success"

	} else {
		results.Message = "decline"

	}
	// results.Data = resp1
	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)

	}
	//..a.a.
	fmt.Fprintf(w, "%s\n", output)

}
