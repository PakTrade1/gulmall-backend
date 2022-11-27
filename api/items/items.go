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
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Price  int32              `json:"price,omitempty`
	Status string             `json:"status"`

	Name struct {
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
	AvailableColor []primitive.ObjectID `json:"available_color,omitempty"`
	// } `json:"available_color"`
	AvailableSize []primitive.ObjectID `json:"available_size,omitempty"`

	// AvailableSize []struct {
	// 	ID primitive.ObjectID `bson:"_id,omitempty"`
	// } `json:"available_size"`
	Images struct {
		Highquility []string `json:"highquility"`
		Lowquility  []string `json:"lowquility"`
	} `json:"images"`
	Price int `json:"price"`
}
type status_req struct {
	Status string `json:"status"`
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
	var results []ItemType
	for cursor.Next(context.TODO()) {
		var abc ItemType
		cursor.Decode(&abc)
		results = append(results, abc)

	}

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
				"available_color": strcutinit.AvailableColor,

				"available_size": strcutinit.AvailableSize,

				"price": strcutinit.Price,
			}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	//end update

	output, err2 := json.MarshalIndent(result1, "", "    ")
	if err2 != nil {
		panic(err2)
	}

	fmt.Fprintf(w, "%s\n", output)

}

type delete_id struct {
	Item_id string `json:"item_id"`
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
	var results []ItemType
	for cursor.Next(context.TODO()) {
		var abc ItemType
		cursor.Decode(&abc)
		results = append(results, abc)

	}

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
				"status": strcutinit.Status,
			}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	//end update

	output, err2 := json.MarshalIndent(result1, "", "    ")
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

func Get_all_items(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	coll := docking.PakTradeDb.Collection("cloths")

	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver
	ctx := context.TODO()

	// Open an aggregation cursor
	cursor, err := coll.Aggregate(ctx, bson.A{
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
	var results []ItemType
	for cursor.Next(context.TODO()) {
		var abc ItemType
		cursor.Decode(&abc)
		results = append(results, abc)

	}

	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)

	}
	fmt.Fprintf(w, "%s\n", output)

}
