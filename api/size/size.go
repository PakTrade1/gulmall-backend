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

type women_fashion_size_cahrt struct {
	ID    string `json:"_id"`
	Name  string `json:"name"`
	Chart struct {
		Round_waist_width           []string `json:"round_waist_width"`
		Round_inner_hip_width       []string `json:"round_inner_hip_width"`
		Round_inner_thigh_at_crotch []string `json:"round_inner_thigh_at_crotch"`
		Outseam_length              []string `json:"outseam_length"`
	} `json:"chart"`
}
type men_fashion_size_cahrt struct {
	ID    status_result `json:"_id"`
	Name  string        `json:"name"`
	Chart struct {
		// Size          []string `json:"size"`
		Body_length   []string `json:"body_length"`
		Chest         []string `json:"chest"`
		Shoulder      []string `json:"shoulder"`
		Sleeve_length []string `json:"sleeve_length"`
	} `json:"chart"`
}
type child_fashion_size_cahrt struct {
	ID    string `json:"_id"`
	Name  string `json:"name"`
	Chart struct {
		Length                    []string `json:"length"`
		Round_waist_width_relaxed []string `json:"round_waist_width_relaxed"`
		Round_hip_width           []string `json:"round_hip_width"`
		Round_thigh_at_crotch     []string `json:"round_thigh_at_crotch"`
	} `json:"chart"`
}

type men_respone_struct_child struct {
	Status     int                      `json:"status"`
	Message    string                   `json:"message"`
	Data       []size_chart             `json:"data"`
	Size_Chart []men_fashion_size_cahrt `json:"size_chart"`
}
type women_respone_struct_child struct {
	Status     int                        `json:"status"`
	Message    string                     `json:"message"`
	Data       []size_chart               `json:"data"`
	Size_Chart []women_fashion_size_cahrt `json:"size_chart"`
}
type child_respone_struct_child struct {
	Status     int                        `json:"status"`
	Message    string                     `json:"message"`
	Data       []size_chart               `json:"data"`
	Size_Chart []child_fashion_size_cahrt `json:"size_chart"`
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
	type_Feild := "chart." + search1.Type1 + ".size"
	type_Feild1 := "$chart." + search1.Type1 + ".size"

	mongoqury := bson.A{
		bson.D{{"$match", bson.D{{"sub_cat_child_id", objectIDS}}}},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", type_Feild1},
					{"includeArrayIndex", "index"},
				},
			},
		},
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
		bson.D{{"$set", bson.D{{"result", bson.D{{"$first", "$result"}}}}}},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$_id"},
					{"size", bson.D{{"$push", "$result"}}},
				},
			},
		},
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
	//// size chart with respt to men women and junior
	var men_results []men_fashion_size_cahrt
	var women_results []women_fashion_size_cahrt
	var child_results []child_fashion_size_cahrt

	if search1.Type1 == "men" {
		men_query := bson.A{
			bson.D{{"$match", bson.D{{"sub_cat_child_id", objectIDS}}}},
			bson.D{
				{"$project",
					bson.D{
						{"name", "$name"},
						{"chart", "$chart.men"},
					},
				},
			},
		}
		gender_result, err1 := coll.Aggregate(context.TODO(), men_query)
		if err1 != nil {
			fmt.Println("errror retrieving user userid : " + objectIDS.Hex())
		}
		for gender_result.Next(context.TODO()) {
			var abc men_fashion_size_cahrt
			gender_result.Decode(&abc)
			men_results = append(men_results, abc)

		}

		var responce men_respone_struct_child
		if results != nil {
			responce.Status = http.StatusOK
			responce.Message = "success"
			responce.Data = results
			responce.Size_Chart = men_results
		} else {
			responce.Status = http.StatusBadRequest
			responce.Message = "declined"
		}
		output, err := json.MarshalIndent(responce, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "%s\n", output)

	} else if search1.Type1 == "women" {

		women_query := bson.A{
			bson.D{{"$match", bson.D{{"sub_cat_child_id", objectIDS}}}},
			bson.D{
				{"$project",
					bson.D{
						{"name", "$name"},
						{"chart", "$chart.women"},
					},
				},
			},
		}

		gender_result, err1 := coll.Aggregate(context.TODO(), women_query)
		if err1 != nil {
			fmt.Println("errror retrieving user userid : " + objectIDS.Hex())
		}
		for gender_result.Next(context.TODO()) {
			var abc women_fashion_size_cahrt
			gender_result.Decode(&abc)
			women_results = append(women_results, abc)
		}
		var responce women_respone_struct_child
		if results != nil {
			responce.Status = http.StatusOK
			responce.Message = "success"
			responce.Data = results
			responce.Size_Chart = women_results
		} else {
			responce.Status = http.StatusBadRequest
			responce.Message = "declined"
		}
		output, err := json.MarshalIndent(responce, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "%s\n", output)

	} else if search1.Type1 == "junior" {

		women_query := bson.A{
			bson.D{{"$match", bson.D{{"sub_cat_child_id", objectIDS}}}},
			bson.D{
				{"$project",
					bson.D{
						{"name", "$name"},
						{"chart", "$chart.junior"},
					},
				},
			},
		}

		gender_result, err1 := coll.Aggregate(context.TODO(), women_query)
		if err1 != nil {
			fmt.Println("errror retrieving user userid : " + objectIDS.Hex())
		}
		for gender_result.Next(context.TODO()) {
			var abc child_fashion_size_cahrt
			gender_result.Decode(&abc)
			child_results = append(child_results, abc)
		}
		var responce child_respone_struct_child
		if results != nil {
			responce.Status = http.StatusOK
			responce.Message = "success"
			responce.Data = results
			responce.Size_Chart = child_results
		} else {
			responce.Status = http.StatusBadRequest
			responce.Message = "declined"
		}
		output, err := json.MarshalIndent(responce, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "%s\n", output)

	} else {
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

// func Size_chart(w http.ResponseWriter, req *http.Request) {

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")

// 	var search1 size_chart_search
// 	err := json.NewDecoder(req.Body).Decode(&search1)
// 	if err != nil {
// 		panic(err)
// 	}
// 	coll := docking.PakTradeDb.Collection("size_chart")
// 	objectIDS, _ := primitive.ObjectIDFromHex(search1.Child_cat_id)

// 	//// size chart with respt to men women and junior
// 	var men_results []men_fashion_size_cahrt
// 	//	var women_results []women_fashion_size_cahrt

// 	//	if search1.Type1 == "men" {

// 	men_query := bson.A{
// 		bson.D{{"$match", bson.D{{"sub_cat_child_id", objectIDS}}}},
// 		bson.D{
// 			{"$project",
// 				bson.D{
// 					{"name", "$name"},
// 					{"chart", "$chart.men"},
// 				},
// 			},
// 		},
// 	}
// 	gender_result, err1 := coll.Aggregate(context.TODO(), men_query)
// 	if err1 != nil {
// 		fmt.Println("errror retrieving user userid : " + objectIDS.Hex())
// 	}
// 	for gender_result.Next(context.TODO()) {
// 		var abc men_fashion_size_cahrt
// 		gender_result.Decode(&abc)
// 		men_results = append(men_results, abc)

// 	}
// 	// } else if search1.Type1 == "women" {

// 	// women_query := bson.A{
// 	// 	bson.D{{"$match", bson.D{{"sub_cat_child_id", objectIDS}}}},
// 	// 	bson.D{
// 	// 		{"$group",
// 	// 			bson.D{
// 	// 				{"_id",
// 	// 					bson.D{
// 	// 						{"name", "$name"},
// 	// 						{"chart", "$chart.women"},
// 	// 					},
// 	// 				},
// 	// 			},
// 	// 		},
// 	// 	},
// 	// }

// 	// gender_result, err1 := coll.Aggregate(context.TODO(), women_query)
// 	// if err1 != nil {
// 	// 	fmt.Println("errror retrieving user userid : " + objectIDS.Hex())
// 	// }
// 	// for gender_result.Next(context.TODO()) {
// 	// 	var abc women_fashion_size_cahrt
// 	// 	gender_result.Decode(&abc)
// 	// 	women_results = append(women_results, abc)
// 	// }

// 	// }
// 	////////end of if else men women
// 	//var responce respone_struct_child
// 	// if results != nil {
// 	// 	responce.Status = http.StatusOK
// 	// 	responce.Message = "success"
// 	// 	responce.Data = results
// 	// } else {
// 	// 	responce.Status = http.StatusBadRequest
// 	// 	responce.Message = "declined"
// 	// }
// 	output, err := json.MarshalIndent(men_results, "", "    ")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Fprintf(w, "%s\n", output)

//}
