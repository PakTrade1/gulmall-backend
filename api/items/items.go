package items

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

type ItemType struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Price         int32              `json:"price,omitempty`
	Status        string             `json:"status"`
	Category      string             `json:"category"`
	Gender        string             `json:"gender"`
	Sub_category  string             `json:"sub_category"`
	Country       string             `json:"country"`
	Qty           int                `json:"qty"`
	Remaining_qty int                `json:"remaining_qty"`
	Name          struct {
		En string `json:"en,omitempty"`
		Ar string `json:"ar,omitempty"`
	} `json:"name,omitempty"`

	Feature []struct {
		//Low_quility  []string `json:"low_quility,omitempty"`
		//High_quility []string `json:"high_quility,omitempty"`
		Name struct {
			En string `json:"en,omitempty"`
			Ar string `json:"ar,omitempty"`
		} `json:"name,omitempty"`
	} `json:"feature,omitempty"`

	Available_size []string `json:"available_size,omitempty"`

	Images []struct {
		Low_quility  []string `json:"low_quility,omitempty"`
		High_quility []string `json:"high_quility,omitempty"`
	} `json:"images,omitempty"`

	Available_color []struct {
		ID     primitive.ObjectID `bson:"_id,omitempty"`
		CSSHex string             `json:"cssHex,omitempty"`
		Name   struct {
			En string `json:"en,omitempty"`
			Ar string `json:"ar,omitempty"`
		} `json:"name,omitempty"`
	} `json:"available_color,omitempty"`
}

type update_item struct {
	ID string `json:"item_id"`

	Name struct {
		En string `json:"en"`
		Ar string `json:"ar"`
	} `json:"name"`
	Feature []struct {
		Name struct {
			En string `json:"en"`
			Ar string `json:"ar"`
		} `json:"name"`
	} `json:"feature"`
	Available_color []primitive.ObjectID `json:"available_color,omitempty"`
	// } `json:"available_color"`
	Size struct {
		Available_size []primitive.ObjectID `json:"available_size,omitempty"`
		Size_chart     string               `json:"size_chart"`
	} `json:"size"`

	// AvailableSize []struct {
	// 	ID primitive.ObjectID `bson:"_id,omitempty"`
	// } `json:"available_size"`
	Images struct {
		Highquility []string `json:"highquility"`
		Lowquility  []string `json:"lowquility"`
	} `json:"images"`
	Price        int    `json:"price"`
	Gender       string `json:"gender"`
	Category     string `json:"category"`
	Sub_category string `json:"sub_category"`
}
type status_req struct {
	Status string `json:"status"`
}
type respone_struct struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    []ItemType `json:"data"`
}

func Items(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var Status_req status_req
	err := json.NewDecoder(req.Body).Decode(&Status_req)
	if err != nil {
		panic(err)
	}

	coll := docking.PakTradeDb.Collection("cloths")

	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver
	ctx := context.TODO()

	// Open an aggregation cursor
	cursor, err := coll.Aggregate(ctx, bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: Status_req.Status}}}},
		bson.D{

			{Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$available_color"},
					{Key: "includeArrayIndex", Value: "index_1"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "color"},
					{Key: "localField", Value: "available_color"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "result"},
				},
			},
		},
		bson.D{{Key: "$set", Value: bson.D{{Key: "result", Value: bson.D{{Key: "$first", Value: "$result"}}}}}},
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id",
						Value: bson.D{
							{Key: "name", Value: "$name"},
							{Key: "feature", Value: "$feature"},
							{Key: "sizes", Value: "$available_size"},
							{Key: "images", Value: "$images"},
							{Key: "_id", Value: "$_id"},
							{Key: "price", Value: "$price"},
							{Key: "status", Value: "$status"},
						},
					},
					{Key: "colors", Value: bson.D{{Key: "$push", Value: "$result"}}},
				},
			},
		},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: "$_id._id"},
					{Key: "name", Value: "$_id.name"},
					{Key: "feature", Value: "$_id.feature"},
					{Key: "available_size", Value: "$_id.sizes"},
					{Key: "images", Value: "$_id.images"},
					{Key: "price", Value: "$_id.price"},
					{Key: "status", Value: "$_id.status"},
					{Key: "available_color", Value: "$colors"},
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	var results respone_struct
	var resp1 []ItemType
	for cursor.Next(context.TODO()) {
		var xy ItemType
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

// func ItemInsertone(w http.ResponseWriter, req *http.Request) {

// 	w.Header().Set("Content-Type", "application/json")
// 	// w.Header().Set("Access-Control-Allow-Origin", "*")

// 	var strcutinit inset_item
// 	err := json.NewDecoder(req.Body).Decode(&strcutinit)
// 	if err != nil {
// 		panic(err)
// 	}

// 	insertdat := bson.M{"name": bson.M{
// 		"en": strcutinit.Name.En,
// 		"ar": strcutinit.Name.Ar,
// 	},
// 		"feature": bson.A{
// 			strcutinit.Feature,
// 		},
// 		"available_color": bson.A{
// 			strcutinit.AvailableColor,
// 		},
// 		"available_size": bson.A{
// 			strcutinit.AvailableSize,
// 		},
// 		"images": bson.A{
// 			strcutinit.Images,
// 		},
// 		"price": strcutinit.Price,
// 	}

// 	//fmt.Print(body)
// 	coll := docking.PakTradeDb.Collection("cloths")

// 	// // // insert a user

// 	inset, err3 := coll.InsertOne(context.TODO(), insertdat)
// 	if err3 != nil {
// 		fmt.Fprintf(w, "%s\n", err3)
// 	}

// 	fmt.Fprintf(w, "%s\n", inset.InsertedID)

// }
func Item_update_one(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var strcutinit update_item
	err := json.NewDecoder(req.Body).Decode(&strcutinit)
	if err != nil {
		panic(err)
	}
	coll := docking.PakTradeDb.Collection("cloths")
	objectIDS, _ := primitive.ObjectIDFromHex(strcutinit.ID)
	gender_objectIDS, _ := primitive.ObjectIDFromHex(strcutinit.Gender)
	sub_category_objectIDS, _ := primitive.ObjectIDFromHex(strcutinit.Sub_category)
	category_objectIDS, _ := primitive.ObjectIDFromHex(strcutinit.Category)
	size_chart_objectIDS, _ := primitive.ObjectIDFromHex(strcutinit.Size.Size_chart)

	// fmt.Print(objectIDS)

	result1, err := coll.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectIDS},
		bson.D{
			{Key: "$set", Value: bson.M{
				"name": bson.M{
					"en": strcutinit.Name.En,
					"ar": strcutinit.Name.Ar,
				},
				"feature":         strcutinit.Feature,
				"available_color": strcutinit.Available_color,

				"size": bson.M{
					"available_size": strcutinit.Size.Available_size,

					"size_chart": size_chart_objectIDS,
				},

				"price":        strcutinit.Price,
				"status":       "pending",
				"gender":       gender_objectIDS,
				"category":     category_objectIDS,
				"sub_category": sub_category_objectIDS,
			}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	//end update
	var record update_resp
	if result1.ModifiedCount >= 1 {
		record.Status = http.StatusOK
		record.Message = "success"
		record.Update_record = int(result1.ModifiedCount)
	} else {
		record.Status = http.StatusBadRequest
		record.Message = "decline"
		record.Update_record = 0
	}
	output, err2 := json.MarshalIndent(record, "", "    ")
	if err2 != nil {
		panic(err2)
	}

	fmt.Fprintf(w, "%s\n", output)
}

type delete_id struct {
	Item_id string `json:"item_id"`
}
type serch_itme_struct struct {
	Status  string `json:"status"`
	Message string `json:"message "`
	Data    ItemType
}

func Serch_item_by_id(w http.ResponseWriter, r *http.Request) {
	coll := docking.PakTradeDb.Collection("cloths")

	var id_get delete_id
	err := json.NewDecoder(r.Body).Decode(&id_get)
	handleError(err)
	objectId, err := primitive.ObjectIDFromHex(id_get.Item_id)
	handleError(err)

	cursor, err := coll.Aggregate(context.TODO(), bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: objectId}}}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "color"},
					{Key: "localField", Value: "available_color"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "color_result"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "size"},
					{Key: "localField", Value: "size.available_size"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "size_result"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "gender"},
					{Key: "localField", Value: "gender"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "gender_result"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "categories"},
					{Key: "localField", Value: "category"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "categories"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "sub_category"},
					{Key: "localField", Value: "sub-category"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "sub_categories"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "sub_category_child"},
					{Key: "localField", Value: "sub-category"},
					{Key: "foreignField", Value: "sub_category_id"},
					{Key: "as", Value: "sub-category-child"},
				},
			},
		},
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id",
						Value: bson.D{
							{Key: "name", Value: "$name"},
							{Key: "feature", Value: "$feature"},
							{Key: "images", Value: "$images"},
							{Key: "_id", Value: "$_id"},
							{Key: "category", Value: "$categories"},
							{Key: "sub_category", Value: "$sub_categories"},
							{Key: "gender", Value: "$gender_result"},
							{Key: "price", Value: "$price"},
							{Key: "qty", Value: "$qty"},
							{Key: "remaining_qty", Value: "$remaining_qty"},
							{Key: "status", Value: "$status"},
							{Key: "country", Value: "$country"},
							{Key: "color", Value: "$color_result"},
							{Key: "size", Value: "$size_result"},
						},
					},
				},
			},
		},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "name", Value: "$_id.name"},
					{Key: "feature", Value: "$_id.feature"},
					{Key: "images", Value: "$_id.images"},
					{Key: "_id", Value: "$_id._id"},
					{Key: "category", Value: "$_id.category"},
					{Key: "sub_category", Value: "$_id.sub_category"},
					{Key: "gender", Value: "$_id.gender"},
					{Key: "price", Value: "$_id.price"},
					{Key: "qty", Value: "$_id.qty"},
					{Key: "remaining_qty", Value: "$_id.remaining_qty"},
					{Key: "status", Value: "$_id.status"},
					{Key: "country", Value: "$_id.country"},
					{Key: "color", Value: "$_id.color"},
					{Key: "size", Value: "$_id.size"},
				},
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	results := new(respone_struct1)
	var resp1 []AutoGenerated
	for cursor.Next(context.TODO()) {
		var xy AutoGenerated
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

// ////////////////// delte item by id
type delte_status struct {
	Item_id string `json:"item_id"`
	Status  string `josn:"status"`
}
type update_resp struct {
	Status        int    `json:"status"`
	Message       string `json:"message"`
	Update_record int    `json:"update_record"`
}

func Item_delete_by_id(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var strcutinit delte_status
	err := json.NewDecoder(req.Body).Decode(&strcutinit)
	if err != nil {
		panic(err)
	}
	coll := docking.PakTradeDb.Collection("cloths")
	objectIDS, _ := primitive.ObjectIDFromHex(strcutinit.Item_id)
	// fmt.Print(objectIDS)

	result1, err := coll.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectIDS},
		bson.D{
			{Key: "$set", Value: bson.M{
				"status": "inactive",
			}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	//end update
	var record update_resp
	if result1.ModifiedCount == 1 {
		record.Status = http.StatusOK
		record.Message = "success"
		record.Update_record = int(result1.ModifiedCount)
	} else {
		record.Status = http.StatusBadRequest
		record.Message = "decline"
		record.Update_record = 0
	}
	output, err2 := json.MarshalIndent(record, "", "    ")
	if err2 != nil {
		panic(err2)
	}

	fmt.Fprintf(w, "%s\n", output)

}

//////////////end of delte item
//code not use it delete item form azure
/*
type AutoGenerated struct {
	Name struct {
		En string `json:"en"`
		Ar string `json:"ar"`
	} `json:"name"`
	Feature []struct {
		Name struct {
			En string `json:"en"`
			Ar string `json:"ar"`
		} `json:"name"`
	} `json:"feature"`
	AvailableColor []struct {
		Id string `json:"$id"`
	} `json:"available_color"`
	AvailableSize []string `json:"available_size"`
	Images        []struct {
		Low_quility  []string `json:"low_quility,omitempty"`
		High_quility []string `json:"high_quility,omitempty"`
	} `json:"images,omitempty"`
	Price int `json:"price"`
}

type delete_id struct {
	Delete_id string `json:"delete_id"`
}
type resp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Id      string `json:"item_id"`
}

var wg sync.WaitGroup

type blob_path struct {
	Blobpath []string
}

func Delte_item(w http.ResponseWriter, r *http.Request) {
	coll := docking.PakTradeDb.Collection("cloths")
	azblob.ParseURL("")

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var id_get delete_id
	err := json.NewDecoder(r.Body).Decode(&id_get)
	handleError(err)
	objectId, err := primitive.ObjectIDFromHex(id_get.Delete_id)
	handleError(err)
	// opts := options.Delete().SetHint(bson.D{{"_id", objectId}})
	// result, err := coll.DeleteMany(context.TODO(), bson.M{"_id": objectId}, nil)
	// if err != nil {
	// 	panic(err)
	// }

	cursor := coll.FindOne(context.Background(), bson.M{"_id": objectId})
	var results AutoGenerated
	var imagearry []string
	cursor.Decode(&results)
	// if results.Images[0].Low_quility[0] != "no image" {
	for _, a := range results.Images {

		imagearry = append(imagearry, a.Low_quility...)

		// fmt.Print(x)
	}
	// } else {
	// 	for _, a := range results.Images {
	// 		imagearry = append(imagearry, a.High_quility...)
	// 	}
	// }
	for i, a := range imagearry {
		res1 := strings.Split(a, "https://paktradegallery.blob.core.windows.net/gallerycontainer/")
		justString := fmt.Sprint(res1)
		t := strings.Replace(justString, "[", "", -1)
		m := strings.Replace(t, "]", "", -1)
		go Deltefile(w, m)
		wg.Add(1)
		// c := Deltefile(m)
		// fmt.Fprintf(w, string(c))
		fmt.Print(i)
	}
	wg.Wait()
}

func Deltefile(w http.ResponseWriter, path string) {
	fmt.Print("call func")
	var mesage1 resp
	_, err1 := client.DeleteBlob(context.TODO(), "gallerycontainer", "\""+path+"\"", nil)
	fmt.Println("fucnc call 1 time of    ", "\""+path+"\"")

	if err1 != nil {
		mesage1 = resp{
			Status:  http.StatusBadGateway,
			Message: "image not found",
		}
		output, err2 := json.MarshalIndent(mesage1, "", "    ")
		if err2 != nil {
			panic(err2)
		}

		fmt.Fprintf(w, "%s\n", output)
		wg.Done()

	} else {
		mesage1 = resp{
			Status:  http.StatusOK,
			Message: "delete Image successful",
		}

		output, err2 := json.MarshalIndent(mesage1, "", "    ")
		if err2 != nil {
			panic(err2)
		}

		fmt.Fprintf(w, "%s\n", output)
		wg.Done()
	}

}
*/

type respone_struct1 struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    []AutoGenerated `json:"data"`
}
type AutoGenerated struct {
	Name struct {
		En string `json:"en"`
		Ar string `json:"ar"`
	} `json:"name"`
	Feature []struct {
		Name struct {
			En string `json:"en"`
			Ar string `json:"ar"`
		} `json:"name"`
	} `json:"feature"`
	Images []struct {
		Low_quility  []string `json:"low_quility,omitempty"`
		High_Quility []string `json:"high_quility,omitempty"`
	} `json:"images"`
	ID primitive.ObjectID `bson:"_id,omitempty"`

	Category []struct {
		ID   primitive.ObjectID `bson:"_id,omitempty"`
		Name struct {
			En string `json:"en"`
			Ar string `json:"ar"`
		} `json:"name"`
		Gender_flag bool   `json:"gender_flag"`
		Icon        string `json:"icon"`
	} `json:"category"`
	Sub_category []struct {
		ID primitive.ObjectID `bson:"_id,omitempty"`

		Cat_id primitive.ObjectID `bson:"cat_id,omitempty"`

		Name struct {
			En string `json:"en"`
			Ar string `json:"ar"`
		} `json:"name"`
	} `json:"sub_category"`
	Gender []struct {
		ID primitive.ObjectID `bson:"_id,omitempty"`

		Name struct {
			En string `json:"en"`
			Ar string `json:"ar"`
		} `json:"name"`
	} `json:"gender"`
	Price         int    `json:"price"`
	Qty           int    `json:"qty"`
	Remaining_qty int    `json:"remaining_qty"`
	Status        string `json:"status"`
	Country       string `json:"country"`
	Color         []struct {
		ID primitive.ObjectID `bson:"_id,omitempty"`

		CSSHex string `json:"cssHex"`
		Name   struct {
			En string `json:"en"`
			Ar string `json:"ar"`
		} `json:"name"`
	} `json:"available_color"`
	Size []struct {
		ID primitive.ObjectID `bson:"_id,omitempty"`

		Name struct {
			En string `json:"en"`
			Ar string `json:"ar"`
		} `json:"name"`
	} `json:"available_size"`
}

func Get_all_items(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	coll := docking.PakTradeDb.Collection("cloths")

	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver
	ctx := context.TODO()

	mongoquery := bson.A{
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "color"},
					{Key: "localField", Value: "available_color"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "color_result"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "size"},
					{Key: "localField", Value: "size.available_size"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "size_result"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "gender"},
					{Key: "localField", Value: "gender"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "gender_result"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "categories"},
					{Key: "localField", Value: "category"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "categories"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "sub_category"},
					{Key: "localField", Value: "sub-category"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "sub_categories"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "sub_category_child"},
					{Key: "localField", Value: "sub-category"},
					{Key: "foreignField", Value: "sub_category_id"},
					{Key: "as", Value: "sub-category-child"},
				},
			},
		},
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id",
						Value: bson.D{
							{Key: "name", Value: "$name"},
							{Key: "feature", Value: "$feature"},
							{Key: "images", Value: "$images"},
							{Key: "_id", Value: "$_id"},
							{Key: "category", Value: "$categories"},
							{Key: "sub_category", Value: "$sub_categories"},
							{Key: "gender", Value: "$gender_result"},
							{Key: "price", Value: "$price"},
							{Key: "qty", Value: "$qty"},
							{Key: "remaining_qty", Value: "$remaining_qty"},
							{Key: "status", Value: "$status"},
							{Key: "country", Value: "$country"},
							{Key: "color", Value: "$color_result"},
							{Key: "size", Value: "$size_result"},
						},
					},
				},
			},
		},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "name", Value: "$_id.name"},
					{Key: "feature", Value: "$_id.feature"},
					{Key: "images", Value: "$_id.images"},
					{Key: "_id", Value: "$_id._id"},
					{Key: "category", Value: "$_id.category"},
					{Key: "sub_category", Value: "$_id.sub_category"},
					{Key: "gender", Value: "$_id.gender"},
					{Key: "price", Value: "$_id.price"},
					{Key: "qty", Value: "$_id.qty"},
					{Key: "remaining_qty", Value: "$_id.remaining_qty"},
					{Key: "status", Value: "$_id.status"},
					{Key: "country", Value: "$_id.country"},
					{Key: "color", Value: "$_id.color"},
					{Key: "size", Value: "$_id.size"},
				},
			},
		},
	}

	// Open an aggregation cursor
	cursor, err := coll.Aggregate(ctx, mongoquery)
	if err != nil {
		log.Fatal(err)
	}
	results := new(respone_struct1)
	var resp1 []AutoGenerated
	for cursor.Next(context.TODO()) {
		var xy AutoGenerated
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
