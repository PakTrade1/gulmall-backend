package size

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

type Size1 struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	CSSHex string             `json:"cssHex,omitempty"`
	Name   `json:"name,omitempty"`
}
type respone_struct struct {
	Status  int     `json:"status"`
	Message string  `json:"message"`
	Data    []Size1 `json:"data"`
}

func Size(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	coll := docking.PakTradeDb.Collection("size")
	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}

	var results []Size1
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	for cursor.Next(context.TODO()) {
		var abc Size1
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

type size_chart_search struct {
	Child_cat_id string `json:"child_cat_id"`
	Type1        string `json:"type"` // this for sizeing /i.e man , woman,junior/ / food i.e letter, weight,dozen/
}
type size_chart struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Size []struct {
		ID   primitive.ObjectID `bson:"_id,omitempty"`
		Name struct {
			En string `json:"en"`
			Ar string `json:"ar"`
		} `json:"name"`
	} `json:"size"`
}

type respone_struct_child struct {
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Data    []size_chart `json:"data"`
}

func Size_select_by_child_id(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var search1 size_chart_search
	err := json.NewDecoder(req.Body).Decode(&search1)
	if err != nil {
		panic(err)
	}
	coll := docking.PakTradeDb.Collection("size_chart")
	objectIDS, _ := primitive.ObjectIDFromHex(search1.Child_cat_id)

	//var result size_chart
	//filter := bson.M{"_id": objectIDS}
	//bson.D{{"$match", bson.D{{"sub_cat_child_id", objectIDS}}}},

	// mongoqury := bson.A{
	// 	bson.D{{"$match", bson.D{{"sub_cat_child_id", objectIDS}}}},
	// 	bson.D{{"$project", bson.D{{"chart", "$chart." + search1.Type1 + ".size"}}}},
	// 	bson.D{
	// 		{"$lookup",
	// 			bson.D{
	// 				{"from", "size"},
	// 				{"localField", "chart"},
	// 				{"foreignField", "_id"},
	// 				{"as", "result"},
	// 			},
	// 		},
	// 	},
	// 	bson.D{{"$project", bson.D{{"size", "$result"}}}},
	// }
	type_Feild := "chart." + search1.Type1 + ".size"
	type_Feild1 := "$chart." + search1.Type1 + ".size"

	mongoqury :=
		bson.A{
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "size"},
						{"localField", type_Feild},
						{"foreignField", "_id"},
						{"as", "result"},
					},
				},
			},
			bson.D{{"$match", bson.D{{"sub_cat_child_id", objectIDS}}}},
			bson.D{{"$project", bson.D{{"size", type_Feild1}}}},
		}
	cursor, err1 := coll.Aggregate(context.TODO(), mongoqury)
	if err1 != nil {
		fmt.Println("errror retrieving user userid : " + objectIDS.Hex())
	}

	// end findOne
	var results []size_chart
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	for cursor.Next(context.TODO()) {
		var abc size_chart
		cursor.Decode(&abc)
		results = append(results, abc)
	}
	var responce respone_struct_child
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

// ////////// Add size
type size_add struct {
	Name struct {
		Ar string `json:"ar"`
		En string `json:"en"`
	} `json:"name"`
}

type respone_add_category struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    status_result `json:"data"`
}
type status_result struct {
	Status string `json:"status"`
}

func Add_size(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	var strcutinit size_add
	err := json.NewDecoder(req.Body).Decode(&strcutinit)
	if err != nil {
		panic(err)
	}

	insertdat := bson.M{
		"name": bson.M{
			"en": strcutinit.Name.En,
			"ar": strcutinit.Name.Ar,
		},
	}

	//fmt.Print(body)
	coll := docking.PakTradeDb.Collection("size")

	// // // insert a user
	var results respone_add_category

	inset, err3 := coll.InsertOne(context.TODO(), insertdat)
	if err3 != nil {
		fmt.Fprintf(w, "%s\n", err3)
	}

	//	fmt.Fprintf(w, "%s\n", inset)

	if inset != nil && inset.InsertedID != "" {
		results.Status = http.StatusOK
		results.Message = "success"
		results.Data.Status = "add Size successfully"
	} else {
		results.Message = "decline"
		results.Data.Status = "not added "

	}
	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)

	}

	fmt.Fprintf(w, "%s\n", output)

}
