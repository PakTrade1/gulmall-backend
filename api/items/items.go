package items

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	docking "pak-trade-go/Docking"
	"strconv"
	"time"

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
				"sub-category": sub_category_objectIDS,
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
			{Key: "$set",
				Value: bson.D{
					{Key: "cat", Value: bson.D{{Key: "$first", Value: "$categories"}}},
					{Key: "subcat", Value: bson.D{{Key: "$first", Value: "$sub_categories"}}},
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
							{Key: "category", Value: "$cat"},
							{Key: "sub_category", Value: "$subcat"},
							{Key: "gender", Value: "$gender_result"},
							{Key: "price", Value: "$price"},
							{Key: "qty", Value: "$qty"},
							{Key: "remaining_qty", Value: "$remaining_qty"},
							{Key: "status", Value: "$status"},
							{Key: "country", Value: "$country"},
							{Key: "color", Value: "$color_result"},
							{Key: "size", Value: "$size_result"},
							{Key: "currency", Value: "$currency"},
							{Key: "rating", Value: "$rating"},
							{Key: "number_ratings", Value: "$number_ratings"},
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
					{Key: "currency", Value: "$_id.currency"},
					{Key: "rating", Value: "$_id.rating"},
					{Key: "number_ratings", Value: "$_id.number_ratings"},
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
	Images struct {
		Low_quility []struct {
			Image string `json:"image"`
			Color string `json:"color"`
		} `json:"low_quility,omitempty"`
		High_Quility []string `json:"high_quility,omitempty"`
	} `json:"images"`
	ID primitive.ObjectID `bson:"_id,omitempty"`

	Category struct {
		ID   primitive.ObjectID `bson:"_id,omitempty"`
		Name struct {
			En string `json:"en"`
			Ar string `json:"ar"`
		} `json:"name"`
		Gender_flag bool   `json:"gender_flag"`
		Icon        string `json:"icon"`
	} `json:"category"`
	Sub_category struct {
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
	Currency       string  `json:"currency"`
	Rating         float64 `json:"rating"`
	Number_ratings int32   `json:"number_ratings"`
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
			{Key: "$set",
				Value: bson.D{
					{Key: "cat", Value: bson.D{{Key: "$first", Value: "$categories"}}},
					{Key: "subcat", Value: bson.D{{Key: "$first", Value: "$sub_categories"}}},
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
							{Key: "category", Value: "$cat"},
							{Key: "sub_category", Value: "$subcat"},
							{Key: "gender", Value: "$gender_result"},
							{Key: "price", Value: "$price"},
							{Key: "qty", Value: "$qty"},
							{Key: "remaining_qty", Value: "$remaining_qty"},
							{Key: "status", Value: "$status"},
							{Key: "country", Value: "$country"},
							{Key: "color", Value: "$color_result"},
							{Key: "size", Value: "$size_result"},
							{Key: "currency", Value: "$currency"},
							{Key: "rating", Value: "$rating"},
							{Key: "number_ratings", Value: "$number_ratings"},
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
					{Key: "currency", Value: "$_id.currency"},
					{Key: "rating", Value: "$_id.rating"},
					{Key: "number_ratings", Value: "$_id.number_ratings"},
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

type Item_struct_insert struct {
	ItemId string "itemId"

	Price   int    `json:"price"`
	Name    string `json:"name"`
	Feature []struct {
		Name string `json:"name"`
	} `json:"feature"`
	Available_color []primitive.ObjectID `json:"availableColor"`
	Images          struct {
		LowQuility []struct {
			Image string `json:"image"`
			Color string `json:"color"`
		} `json:"lowQuility"`
		HighQuility []string `json:"highQuility"`
	} `json:"images"`
	Status   string `json:"status"`
	Gender   string `json:"gender"`
	Category string `json:"category"`
	Size     struct {
		Available_size []primitive.ObjectID `json:"availableSize"`
		SizeChart      primitive.ObjectID   `json:"sizeChart"`
	} `json:"size"`
	Country       string  `json:"country"`
	Qty           int     `json:"qty"`
	Currency      string  `json:"currency"`
	Rating        float64 `json:"rating"`
	Title         string  `json:"title"`
	OwnerID       string  `json:"ownerId"`
	NumberRatings int     `json:"numberRatings"`
	RemainingQty  int     `json:"remainingQty"`
	SubCategory   string  `json:"subCategory"`
	// CreationTimestamp time.Time `json:"creationTimestamp"`
	HasDimension bool   `json:"hasDimension"`
	ParentID     string `json:"parentId"`
}
type Prant_struct struct {
	Price  int `json:"price"`
	Images struct {
		HighQuility []string `json:"highQuility"`
		LowQuility  []struct {
			Image string `json:"image"`
			Color string `json:"color"`
		} `json:"lowQuility"`
	} `json:"images"`
	Status        string             `json:"status"`
	Category      primitive.ObjectID `json:"category"`
	Country       string             `json:"country"`
	Qty           int                `json:"qty"`
	Currency      string             `json:"currency"`
	Rating        float64            `json:"rating"`
	Title         string             `json:"title"`
	OwnerID       primitive.ObjectID `json:"ownerId"`
	NumberRatings int                `json:"numberRatings"`
	RemainingQty  int                `json:"remainingQty"`
	SubCategory   primitive.ObjectID `json:"subCategory"`
	// CreationTimestamp time.Time `json:"creationTimestamp"`
}
type update_parent_item_cat struct {
	Status  int        `json:"status"`
	Message string     `json:"message "`
	Data    update_res `json:"data"`
}
type update_res struct {
	StatusParent string `json:"statusParent"`
	StatusChild  string `json:"statusChild"`
}

func Add_item_update(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var strcutinit Item_struct_insert
	err := json.NewDecoder(req.Body).Decode(&strcutinit)
	if err != nil {
		panic(err)
	}
	item_id, err := primitive.ObjectIDFromHex(strcutinit.ItemId)
	handleError(err)

	Cat, err := primitive.ObjectIDFromHex(strcutinit.Category)
	handleError(err)
	subCat, err := primitive.ObjectIDFromHex(strcutinit.SubCategory)
	handleError(err)
	OwID, err := primitive.ObjectIDFromHex(strcutinit.OwnerID)
	handleError(err)

	parentstruct_in := new(Prant_struct)
	parentstruct_in.Price = strcutinit.Price
	parentstruct_in.Status = strcutinit.Status
	parentstruct_in.Category = Cat
	parentstruct_in.Country = strcutinit.Country
	parentstruct_in.Qty = strcutinit.Qty
	parentstruct_in.Currency = strcutinit.Currency
	parentstruct_in.Rating = strcutinit.Rating
	parentstruct_in.Title = strcutinit.Title
	parentstruct_in.OwnerID = OwID
	parentstruct_in.NumberRatings = strcutinit.NumberRatings
	parentstruct_in.RemainingQty = strcutinit.RemainingQty
	parentstruct_in.SubCategory = subCat
	// parentstruct_in.CreationTimestamp = strcutinit.CreationTimestamp

	coll := docking.PakTradeDb.Collection("items-parent")
	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver
	ctx := context.TODO()

	result1, err := coll.UpdateOne(
		ctx,
		bson.M{"_id": item_id},
		bson.D{
			{Key: "$set", Value: bson.M{
				"price":             strcutinit.Price,
				"status":            "pending",
				"category":          Cat,
				"subCategory":       subCat,
				"country":           parentstruct_in.Country,
				"qty":               parentstruct_in.Qty,
				"currency":          parentstruct_in.Currency,
				"rating":            parentstruct_in.Rating,
				"title":             parentstruct_in.Title,
				"ownerId":           OwID,
				"numberRatings":     parentstruct_in.NumberRatings,
				"remainingQty":      parentstruct_in.RemainingQty,
				"creationTimestamp": primitive.NewDateTimeFromTime(time.Now()),
			}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	var clot_update_status = 0
	a := new(update_res)
	coll1 := docking.PakTradeDb.Collection("cloths")

	if strcutinit.Category == "63a9a76fd38789473ba919e6" {
		result1, err := coll1.UpdateOne(
			ctx,
			bson.M{"parentId": item_id},
			bson.D{
				{Key: "$set", Value: bson.M{
					"feature": strcutinit.Feature,
					"name":    "Cloth_number_5",
					"size": bson.M{
						"availableSize": strcutinit.Size.Available_size,
						"sizeChart":     strcutinit.Size.SizeChart,
					},
					"gender":         strcutinit.Gender,
					"hasDimension":   strcutinit.HasDimension,
					"availableColor": strcutinit.Available_color,
				}},
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		clot_update_status = int(result1.ModifiedCount)
	}
	var results update_parent_item_cat
	a.StatusParent = strconv.Itoa(int(result1.ModifiedCount))
	a.StatusChild = strconv.Itoa(clot_update_status)

	if result1.ModifiedCount >= 1 {
		results.Status = http.StatusOK
		results.Message = "success"
		results.Data = *a

	} else {
		results.Message = "decline"

	}
	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)

	}
	fmt.Fprintf(w, "%s\n", output)

}

// /////////// Add Image to Parint item with respt to thier cetegory like mart and cloth
type img_add_struct struct {
	CatId  string `json:"catId"`
	Images struct {
		HighQuility []string `json:"highQuility"`
		LowQuility  []struct {
			Image string `json:"image"`
			Color string `json:"color"`
		} `json:"lowQuility"`
	} `json:"images"`
}
type add_img_itme_result struct {
	Status  int    `json:"status"`
	Message string `json:"message "`
	Data    PID    `json:"data"`
}
type PID struct {
	ParentID interface{} `json:"parentId"`
}

func Add_item_img_wrt_category(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var strcutinit img_add_struct
	err := json.NewDecoder(req.Body).Decode(&strcutinit)
	if err != nil {
		panic(err)
	}
	dumy_array := [1]string{"no image"}
	insertdat := bson.M{
		"images": bson.M{
			"highQuility": dumy_array,

			"lowQuility": strcutinit.Images.LowQuility,
		},
	}

	//fmt.Print(body)
	coll := docking.PakTradeDb.Collection("items-parent")

	// // // insert a user

	responceid, err3 := coll.InsertOne(context.TODO(), insertdat)
	if err3 != nil {
		fmt.Print(err3)
	}
	inster_parent_id := bson.M{
		"parentId": responceid.InsertedID,
	}
	if strcutinit.CatId == "63a9a76fd38789473ba919e6" {
		coll1 := docking.PakTradeDb.Collection("cloths")
		_, err4 := coll1.InsertOne(context.TODO(), inster_parent_id)
		if err4 != nil {
			fmt.Print(err4)
		}
	} else if strcutinit.CatId == "63bdb52116cccb9bb8b48388" {
		coll1 := docking.PakTradeDb.Collection("item-mart")
		_, err4 := coll1.InsertOne(context.TODO(), inster_parent_id)
		if err4 != nil {
			fmt.Print(err4)
		}

	}
	////////// Result
	ObjectID_parent := new(PID)
	ObjectID_parent.ParentID = responceid.InsertedID

	var results add_img_itme_result

	if responceid.InsertedID != nil {
		results.Status = http.StatusOK
		results.Message = "success"
		results.Data = *ObjectID_parent

	} else {
		results.Message = "decline"

	}
	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)

	}
	fmt.Fprintf(w, "%s\n", output)
	//////////////// End of image upload_insert

}
