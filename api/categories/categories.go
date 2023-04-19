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
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Gender_flag bool               `json:"gender_flag"`
	Icon        string             `json:"icon"`
	Name        struct {
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
	Icon   string             `json:"icon"`
	Name   struct {
		En string `json:"en"`
		Ar string `json:"ar"`
	} `json:"name"`
}

// // this functhion  is for sub_category collection that serch data w.r.t category id
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
	Chaild_cat_id string `json:"sub_category_id"`
}
type Chaild_sub_Categoies_selected struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Sub_category_id string             `json:"sub_category_id"`
	Name            struct {
		En string `json:"en"`
		Ar string `json:"ar"`
	} `json:"name"`
	Gender_flag bool `json:"gneder_flag"`
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
	cursor, err := coll.Aggregate(context.TODO(), bson.A{
		bson.D{{"$match", bson.D{{"sub_category_id", objectIDS}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "sub_category"},
					{"localField", "sub_category_id"},
					{"foreignField", "_id"},
					{"as", "result"},
				},
			},
		},
		bson.D{{"$set", bson.D{{"cat_id", bson.D{{"$first", "$result.cat_id"}}}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "categories"},
					{"localField", "cat_id"},
					{"foreignField", "_id"},
					{"as", "gender"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$gender"},
					{"includeArrayIndex", "index"},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"sub_category_id", "$sub_category_id"},
					{"name", "$name"},
					{"gender_flag", "$gender.gender_flag"},
				},
			},
		},
	})
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

// / function for add singel category
type respone_add_category struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    status_result `json:"data"`
}
type status_result struct {
	Status    string      `json:"status"`
	Insert_id interface{} `json:"category_id"`
}

func Add_category(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	var strcutinit Categories
	err := json.NewDecoder(req.Body).Decode(&strcutinit)
	if err != nil {
		panic(err)
	}

	insertdat := bson.M{
		"name": bson.M{
			"en": strcutinit.Name.En,
			"ar": strcutinit.Name.Ar,
		},
		"gender_flag": strcutinit.Gender_flag,
		"icon":        strcutinit.Icon,
	}

	//fmt.Print(body)
	coll := docking.PakTradeDb.Collection("categories")

	// // // insert a category
	var results respone_add_category

	inset, err3 := coll.InsertOne(context.TODO(), insertdat)
	if err3 != nil {
		fmt.Fprintf(w, "%s\n", err3)
	}

	//	fmt.Fprintf(w, "%s\n", inset)

	if inset != nil && inset.InsertedID != "" {
		results.Status = http.StatusOK
		results.Message = "success"
		results.Data.Status = "add cetegory successfully"
		results.Data.Insert_id = inset.InsertedID
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

// update category
type cat_update struct {
	Cat_id string `json:"cat_id"`
	Name   struct {
		En string `json:"en"`
		Ar string `json:"ar"`
	} `json:"name"`
	Gender_flag bool `json:"gender_flag"`
}

func Update_Category(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var cat_updt cat_update
	err := json.NewDecoder(req.Body).Decode(&cat_updt)
	if err != nil {
		panic(err)
	}
	coll := docking.PakTradeDb.Collection("categories")
	objectIDS, _ := primitive.ObjectIDFromHex(cat_updt.Cat_id)
	// fmt.Print(objectIDS)

	result1, err := coll.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectIDS},
		bson.D{
			{Key: "$set", Value: bson.M{
				"name": bson.M{
					"en": cat_updt.Name.En,
					"ar": cat_updt.Name.Ar,
				},

				"gender_flag": cat_updt.Gender_flag,
			}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	//end update
	var results respone_add_category
	if result1 != nil && result1.ModifiedCount == 1 {
		results.Status = http.StatusOK
		results.Message = "success"
		results.Data.Status = " update successfully"
	} else {
		results.Message = "decline"
		results.Data.Status = "no data update "

	}
	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)

	}

	fmt.Fprintf(w, "%s\n", output)
}

// //////// delte item by id
type delete_cat struct {
	Cat_id string `json:"cat_id"`
}

func Delete_category(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var strcutinit delete_cat
	err := json.NewDecoder(req.Body).Decode(&strcutinit)
	if err != nil {
		panic(err)
	}
	coll := docking.PakTradeDb.Collection("categories")
	objectIDS, _ := primitive.ObjectIDFromHex(strcutinit.Cat_id)
	// fmt.Print(objectIDS)

	res, err := coll.DeleteOne(context.TODO(), bson.D{{"_id", objectIDS}})
	if err != nil {
		log.Fatal(err)
	}

	var results respone_add_category
	if res != nil && res.DeletedCount == 1 {
		results.Status = http.StatusOK
		results.Message = "success"
		results.Data.Status = "delete successfully"

	} else {
		results.Message = "decline"
		results.Data.Status = "no data delete"

	}
	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)

	}

	fmt.Fprintf(w, "%s\n", output)

}

// ////////////// Insert Sub cetogery and responce id and name
type Sub_cat_requset struct {
	Cat_id primitive.ObjectID `json:"cat_id"`
	Icon   string             `json:"icon"`
	Name   struct {
		En string `json:"en"`
		Ar string `json:"ar"`
	} `json:"name"`
}

type respone_add_sub_category struct {
	Status  int                   `json:"status"`
	Message string                `json:"message"`
	Data    status_result_sub_cat `json:"data"`
}
type status_result_sub_cat struct {
	Status    string      `json:"status"`
	Insert_id interface{} `json:"category_id"`
}

func Add_sub_category(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	var strcutinit Sub_cat_requset
	err := json.NewDecoder(req.Body).Decode(&strcutinit)
	if err != nil {
		panic(err)
	}

	insertdat := bson.M{
		"cat_id": strcutinit.Cat_id,
		"icon":   strcutinit.Icon,
		"name": bson.M{
			"en": strcutinit.Name.En,
			"ar": strcutinit.Name.Ar,
		},
	}

	//fmt.Print(body)
	coll := docking.PakTradeDb.Collection("sub_category")

	// // // insert a category
	var results respone_add_category

	inset, err3 := coll.InsertOne(context.TODO(), insertdat)
	if err3 != nil {
		fmt.Fprintf(w, "%s\n", err3)
	}

	//	fmt.Fprintf(w, "%s\n", inset)

	if inset != nil && inset.InsertedID != "" {
		results.Status = http.StatusOK
		results.Message = "success"
		results.Data.Status = "add sub-cetegory successfully"
		results.Data.Insert_id = inset.InsertedID
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

// //////////// add sub child  cetegory
type Sub_Child_cat struct {
	// ID          primitive.ObjectID `bson:"_id,omitempty"`

	Sub_category_id primitive.ObjectID `json:"sub_category_id"`
	Icon            string             `json:"icon"`
	Name            struct {
		En string `json:"en"`
		Ar string `json:"ar"`
	} `json:"name"`
	Available_size []string `json:"available_size"`
}
type respone_add_sub_category_child struct {
	Status  int                   `json:"status"`
	Message string                `json:"message"`
	Data    status_result_sub_cat `json:"data"`
}
type status_result_sub_cat_child struct {
	Status    string      `json:"status"`
	Insert_id interface{} `json:"category_id"`
}

func Add_sub_child_category(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	var strcutinit Sub_Child_cat
	err := json.NewDecoder(req.Body).Decode(&strcutinit)
	if err != nil {
		panic(err)
	}

	insertdat := bson.M{
		"sub_category_id": strcutinit.Sub_category_id,
		"icon":            strcutinit.Icon,
		"name": bson.M{
			"en": strcutinit.Name.En,
			"ar": strcutinit.Name.Ar,
		},
		"available_size": strcutinit.Available_size,
	}

	//fmt.Print(body)
	coll := docking.PakTradeDb.Collection("sub_category_child")

	// // // insert a category
	var results respone_add_sub_category_child

	inset, err3 := coll.InsertOne(context.TODO(), insertdat)
	if err3 != nil {
		fmt.Fprintf(w, "%s\n", err3)
	}

	//	fmt.Fprintf(w, "%s\n", inset)

	if inset != nil && inset.InsertedID != "" {
		results.Status = http.StatusOK
		results.Message = "success"
		results.Data.Status = "add sub-cetegory-child successfully"
		results.Data.Insert_id = inset.InsertedID
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
