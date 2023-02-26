package size

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
type type_selection struct {
	Name string `json:"name"`
}

type women_size_show struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Name             string             `json:"name"`
	Sub_cat_child_id primitive.ObjectID `json:"sub_cat_child_id"`
	Chart            struct {
		Round_waist_width           []string `json:"round_waist_width"`
		Round_inner_hip_width       []string `json:"round_inner_hip_width"`
		Round_inner_thigh_at_crotch []string `json:"round_inner_thigh_at_crotch"`
		Outseam_length              []string `json:"outseam_length"`
	} `json:"chart"`
	Size []struct {
		ID   primitive.ObjectID `bson:"_id,omitempty"`
		Name struct {
			En string `json:"en"`
			Ar string `json:"ar"`
		} `json:"name"`
	} `json:"size"`
}
type women_resp struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Data    []women_size_show `json:"data"`
}

// // tShirt women
type women_size_tshirt struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Name             string             `json:"name"`
	Sub_cat_child_id primitive.ObjectID `json:"sub_cat_child_id"`
	Chart            struct {
		Length        []string `json:"length"`
		Chest         []string `json:"chest"`
		Shoulder      []string `json:"shoulder"`
		Sleeve_length []string `json:"sleeve_length"`
	} `json:"chart"`
	Size []struct {
		ID   primitive.ObjectID `bson:"_id,omitempty"`
		Name struct {
			En string `json:"en"`
			Ar string `json:"ar"`
		} `json:"name"`
	} `json:"size"`
}
type women_tShirt_Resp struct {
	Status  int                 `json:"status"`
	Message string              `json:"message"`
	Data    []women_size_tshirt `json:"data"`
}

type Men_women_shoze struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Name             string             `json:"name"`
	Sub_cat_child_id primitive.ObjectID `json:"sub_cat_child_id"`
	Chart            struct {
		Length  []string `json:"length"`
		Eu_size []string `json:"eu_size"`
		Us_size []string `json:"us_size"`
		Uk_size []string `json:"uk_size"`
	} `json:"chart"`
	Size []struct {
		ID   primitive.ObjectID `bson:"_id,omitempty"`
		Name struct {
			En string `json:"en"`
			Ar string `json:"ar"`
		} `json:"name"`
	} `json:"size"`
}
type men_wonen_shoze_resp struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Data    []Men_women_shoze `json:"data"`
}
type men_pent_size struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Name             string             `json:"name"`
	Sub_cat_child_id primitive.ObjectID `json:"sub_cat_child_id"`
	Chart            struct {
		Round_waist_width           []string `json:"round_waist_width"`
		Round_inner_hip_width       []string `json:"round_inner_hip_width"`
		Round_inner_thigh_at_crotch []string `json:"round_inner_thigh_at_crotch"`
		Outseam_length              []string `json:"outseam_length"`
	} `json:"chart"`
	Size []struct {
		ID   primitive.ObjectID `bson:"_id,omitempty"`
		Name struct {
			En string `json:"en"`
			Ar string `json:"ar"`
		} `json:"name"`
	} `json:"size"`
}
type men_pent_size_resp struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    []men_pent_size `json:"data"`
}
type men_shirt_size struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Name             string             `json:"name"`
	Sub_cat_child_id primitive.ObjectID `json:"sub_cat_child_id"`
	Chart            struct {
		Chest         []string `json:"chest"`
		Shoulder      []string `json:"shoulder"`
		Sleeve_length []string `json:"sleeve_length"`
		Body_length   []string `json:"body_length"`
	} `json:"chart"`
	Size []struct {
		ID   primitive.ObjectID `bson:"_id,omitempty"`
		Name struct {
			En string `json:"en"`
			Ar string `json:"ar"`
		} `json:"name"`
	} `json:"size"`
}
type men_shirt_size_resp struct {
	Status  int              `json:"status"`
	Message string           `json:"message"`
	Data    []men_shirt_size `json:"data"`
}
type child_pent_size struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Name             string             `json:"name"`
	Sub_cat_child_id primitive.ObjectID `json:"sub_cat_child_id"`
	Chart            struct {
		Length                    []string `json:"length"`
		Round_waist_width_relaxed []string `json:"round_waist_width_relaxed"`
		Round_hip_width           []string `json:"round_hip_width"`
		Round_thigh_at_crotch     []string `json:"round_thigh_at_crotch"`
	} `json:"chart"`
	Size []struct {
		ID   primitive.ObjectID `bson:"_id,omitempty"`
		Name struct {
			En string `json:"en"`
			Ar string `json:"ar"`
		} `json:"name"`
	} `json:"size"`
}
type child_pent_size_resp struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Data    []child_pent_size `json:"data"`
}
type child_shirt_size struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Name             string             `json:"name"`
	Sub_cat_child_id primitive.ObjectID `json:"sub_cat_child_id"`
	Chart            struct {
		Length              []string `json:"length"`
		Body_length_cropped []string `json:"body_length_cropped"`
		Chest               []string `json:"chest"`
		Shoulder            []string `json:"shoulder"`
		Sleeve_length       []string `json:"sleeve_length"`
		Half_sleeve_length  []string `json:"half_sleeve_length"`
	} `json:"chart"`
	Size []struct {
		ID   primitive.ObjectID `bson:"_id,omitempty"`
		Name struct {
			En string `json:"en"`
			Ar string `json:"ar"`
		} `json:"name"`
	} `json:"size"`
}
type child_shirt_size_resp struct {
	Status  int                `json:"status"`
	Message string             `json:"message"`
	Data    []child_shirt_size `json:"data"`
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
	type_Feild1 := "$chart." + search1.Type1
	var objectIDFromHex = func(hex string) primitive.ObjectID {
		objectID, err := primitive.ObjectIDFromHex(hex)
		if err != nil {
			log.Fatal(err)
		}
		return objectID
	}

	result_selecttion_type, err := coll.Aggregate(context.TODO(), bson.A{
		bson.D{{"$match", bson.D{{"sub_cat_child_id", objectIDFromHex(search1.Child_cat_id)}}}},
		bson.D{{"$project", bson.D{{"name", "$name"}}}},
	})

	var abc type_selection

	for result_selecttion_type.Next(context.TODO()) {
		result_selecttion_type.Decode(&abc)
	}

	mongo_Qury := bson.A{
		bson.D{{"$match", bson.D{{"sub_cat_child_id", objectIDS}}}},
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
		bson.D{
			{"$project",
				bson.D{
					{"name", "$name"},
					{"sub_cat_child_id", "$sub_cat_child_id"},
					{"chart", type_Feild1},
					{"size", "$result"},
				},
			},
		},
	}

	cursor, err1 := coll.Aggregate(context.TODO(), mongo_Qury)
	if err1 != nil {
		fmt.Println("errror retrieving user userid : " + objectIDS.Hex())
	}
	if search1.Type1 == "women" {

		if abc.Name == "Cloth" {
			var results []women_size_show
			if err = cursor.All(context.TODO(), &results); err != nil {
				panic(err)
			}

			for cursor.Next(context.TODO()) {
				var abc women_size_show
				cursor.Decode(&abc)
				results = append(results, abc)
			}
			var responce women_resp
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
		} else if abc.Name == "T-Shirt" {
			var results []women_size_tshirt
			if err = cursor.All(context.TODO(), &results); err != nil {
				panic(err)
			}

			for cursor.Next(context.TODO()) {
				var abc women_size_tshirt
				cursor.Decode(&abc)
				results = append(results, abc)
			}
			var responce women_tShirt_Resp
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
		} else if abc.Name == "Shose" {
			var results []Men_women_shoze
			if err = cursor.All(context.TODO(), &results); err != nil {
				panic(err)
			}

			for cursor.Next(context.TODO()) {
				var abc Men_women_shoze
				cursor.Decode(&abc)
				results = append(results, abc)
			}
			var responce men_wonen_shoze_resp
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
	} else if search1.Type1 == "men" {
		if abc.Name == "Cloth" {
			var results []men_pent_size
			if err = cursor.All(context.TODO(), &results); err != nil {
				panic(err)
			}

			for cursor.Next(context.TODO()) {
				var abc men_pent_size
				cursor.Decode(&abc)
				results = append(results, abc)
			}
			var responce men_pent_size_resp
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
		} else if abc.Name == "T-Shirt" {
			var results []men_shirt_size
			if err = cursor.All(context.TODO(), &results); err != nil {
				panic(err)
			}

			for cursor.Next(context.TODO()) {
				var abc men_shirt_size
				cursor.Decode(&abc)
				results = append(results, abc)
			}
			var responce men_shirt_size_resp
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
		} else if abc.Name == "Shose" {
			var results []Men_women_shoze
			if err = cursor.All(context.TODO(), &results); err != nil {
				panic(err)
			}

			for cursor.Next(context.TODO()) {
				var abc Men_women_shoze
				cursor.Decode(&abc)
				results = append(results, abc)
			}
			var responce men_wonen_shoze_resp
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
	} else if search1.Type1 == "junior" {
		if abc.Name == "Cloth" {
			var results []child_pent_size
			if err = cursor.All(context.TODO(), &results); err != nil {
				panic(err)
			}

			for cursor.Next(context.TODO()) {
				var abc child_pent_size
				cursor.Decode(&abc)
				results = append(results, abc)
			}
			var responce child_pent_size_resp
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
		} else if abc.Name == "T-Shirt" {
			var results []child_shirt_size
			if err = cursor.All(context.TODO(), &results); err != nil {
				panic(err)
			}

			for cursor.Next(context.TODO()) {
				var abc child_shirt_size
				cursor.Decode(&abc)
				results = append(results, abc)
			}
			var responce child_shirt_size_resp
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
		} else if abc.Name == "Shose" {
			var results []Men_women_shoze
			if err = cursor.All(context.TODO(), &results); err != nil {
				panic(err)
			}

			for cursor.Next(context.TODO()) {
				var abc Men_women_shoze
				cursor.Decode(&abc)
				results = append(results, abc)
			}
			var responce men_wonen_shoze_resp
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
	} else {
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

		var responce respone_struct_child
		if result_selecttion_type != nil {
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

func objectIDFromHex(s string) {
	panic("unimplemented")
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
