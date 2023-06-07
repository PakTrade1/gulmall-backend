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
	"go.mongodb.org/mongo-driver/mongo/options"
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
		Highquility []string `json:"highQuality"`
		Lowquility  []string `json:"lowQuality"`
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
	Item_id primitive.ObjectID `json:"itemId"`
}
type serch_itme_struct struct {
	Status  string `json:"status"`
	Message string `json:"message "`
	Data    ItemType
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

type respone_struct1 struct {
	Status      int             `json:"status"`
	Message     string          `json:"message"`
	TotalRecord int             `json:"totalRecord"`
	Data        []AutoGenerated `json:"data"`
}
type respone_st struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    AutoGenerated `json:"data"`
}
type AutoGenerated struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ParentId     primitive.ObjectID `json:"parentId"`
	HasDimension bool               `json:"hasDimension"`
	Name         string             `json:"name"`
	Feature      []struct {
		Name string `json:"name"`
	} `json:"feature"`
	Images []struct {
		Image string `json:"image"`
		Color string `json:"color"`
	} `json:"images"`
	Category struct {
		ID primitive.ObjectID `bson:"_id,omitempty"`

		Name        string `json:"name"`
		Gender_flag bool   `json:"gender_flag"`
		Icon        string `json:"icon"`
	} `json:"category"`
	Sub_category struct {
		ID primitive.ObjectID `bson:"_id,omitempty"`

		Name string `json:"name"`

		Icon string `json:"icon"`
	} `json:"sub_category"`
	Gender        string `json:"gender"`
	Price         int    `json:"price"`
	Qty           int    `json:"qty"`
	Remaining_qty int    `json:"remaining_qty"`
	Status        string `json:"status"`
	Color         []struct {
		ID primitive.ObjectID `bson:"_id,omitempty"`

		CSSHex string `json:"cssHex"`
		Name   string `json:"name"`
	} `json:"color"`
	Size []struct {
		ID primitive.ObjectID `bson:"_id,omitempty"`

		Name string `json:"name"`
	} `json:"size"`
	OwnerName string `json:"ownerName"`
	Title     string `json:"title"`
	Country   string `json:"country"`
	Currency  string `json:"currency"`
	Plan      struct {
		// Price      int    `json:"price"`
		Name  string `json:"name"`
		Order int    `json:"order"`
		// AdDuration int    `json:"adDuration"`
	}

	Dimension Dimension `json:"dimension"`
}
type Dimension struct {
	Length *struct {
		Unit  string  `json:"unit"`
		Value float64 `json:"value"`
	} `json:"length"`
	Width *struct {
		Unit  string `json:"unit"`
		Value int    `json:"value"`
	} `json:"width"`
}

func Get_all_items(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	coll := docking.PakTradeDb.Collection("cloths")
	coll1 := docking.PakTradeDb.Collection("items-parent")

	pageN := req.URL.Query().Get("pageNumber")
	pageNu, err := strconv.Atoi(pageN)
	if err != nil || pageNu <= 0 {
		pageNu = 1
	}

	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver
	ctx := context.TODO()
	pageNumber := pageNu
	pageSize := 10
	mongoquery := []bson.D{

		bson.D{
			{"$lookup",
				bson.D{
					{"from", "items-parent"},
					{"localField", "parentId"},
					{"foreignField", "_id"},
					{"as", "parent"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "color"},
					{"localField", "availableColor"},
					{"foreignField", "_id"},
					{"as", "color"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "size"},
					{"localField", "size.availableSize"},
					{"foreignField", "_id"},
					{"as", "size"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "size_chart"},
					{"localField", "size.sizeChart"},
					{"foreignField", "_id"},
					{"as", "sizeChart"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "sub_category"},
					{"localField", "parent.subCategory"},
					{"foreignField", "_id"},
					{"as", "sub_cat"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "categories"},
					{"localField", "parent.category"},
					{"foreignField", "_id"},
					{"as", "cat"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "plans"},
					{"localField", "parent.planId"},
					{"foreignField", "_id"},
					{"as", "plans"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "Mammalas_login"},
					{"localField", "parent.ownerId"},
					{"foreignField", "_id"},
					{"as", "owner"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "tier"},
					{"localField", "plans.tierId"},
					{"foreignField", "_id"},
					{"as", "tier"},
				},
			},
		},
		bson.D{
			{"$set",
				bson.D{
					{"parent", bson.D{{"$first", "$parent"}}},
					{"category", bson.D{{"$first", "$cat"}}},
					{"subCategory", bson.D{{"$first", "$sub_cat"}}},
					{"plans", bson.D{{"$first", "$plans"}}},
					{"tier", bson.D{{"$first", "$tier"}}},
					{"owner", bson.D{{"$first", "$owner"}}},
				},
			},
		},

		bson.D{
			{"$project",
				bson.D{
					{"name", "$name"},
					{"hasDimension", "$hasDimension"},
					{"feature", "$feature"},
					{"images", "$parent.images"},
					{"_id", "$_id"},
					{"parentId", "$parentId"},
					{"category", "$category"},
					{"sub_category", "$subCategory"},
					{"gender", "$gender"},
					{"price", "$parent.price"},
					{"qty", "$parent.qty"},
					{"remaining_qty", "$parent.remainingQty"},
					{"status", "$parent.status"},
					{"country", "$parent.country"},
					{"currency", "$parent.currency"},

					{"color", "$color"},
					{"size", "$size"},
					{"title", "$parent.title"},
					{"ownerName", "$owner.displayName"},
					{"totalRecord", "$string"},
					{"plan",
						bson.D{
							{"price", "$plans.price"},
							{"name", "$tier.name"},
							{"order", "$tier.order"},
							{"adDuration", "$plans.ad_duration"},
						},
					},
					{"dimension",
						bson.D{
							{"width",
								bson.D{
									{"unit", "$dimension.width.unit"},
									{"value", "$dimension.width.value"},
								},
							},
							{"length",
								bson.D{
									{"unit", "$dimension.length.unit"},
									{"value", "$dimension.length.value"},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$sort",
				bson.D{
					{"plan.order", 1},
				},
			},
		},

		bson.D{
			{"$skip", (pageNumber - 1) * pageSize},
		},
		bson.D{
			{"$limit", pageSize},
		},
	}
	// Open an pagination cursor

	// opts := options.Find().SetSkip(int64((page - 1) * pageSize)).SetLimit(int64(pageSize))
	aggOptions := options.Aggregate()
	aggOptions.SetAllowDiskUse(true)

	cursor, err := coll.Aggregate(ctx, mongoquery, aggOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	///////////////////// Totla Recrod
	pipeline := []bson.D{
		bson.D{
			{"$count", "totalRecords"},
		},
	}

	cursor1, err2 := coll1.Aggregate(ctx, pipeline)
	if err2 != nil {
		log.Fatal(err2)
	}
	defer cursor1.Close(context.TODO())

	var Totalcontrecord struct {
		TotalRecords int `json:"totalRecords"`
	}
	if cursor1.Next(context.TODO()) {
		if err := cursor1.Decode(&Totalcontrecord); err != nil {
			log.Fatal(err)
		}
	}

	///////////////// End total record

	var results respone_struct1
	var resp1 []AutoGenerated

	for cursor.Next(context.TODO()) {
		var xy AutoGenerated
		cursor.Decode(&xy)
		resp1 = append(resp1, xy)

	}

	// ///////////
	// Populate resp1 and resp2 with data

	//////////////

	if resp1 != nil {
		results.Status = http.StatusOK
		results.Message = "success"
		results.TotalRecord = Totalcontrecord.TotalRecords
		results.Data = resp1

	} else {
		results.Message = "decline"

	}

	// results.Data =
	output, err := json.MarshalIndent(resp1, "", "    ")
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
		} `json:"lowQuality"`
		HighQuility []string `json:"highQuality"`
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
		HighQuility []string `json:"highQuality"`
		LowQuility  []struct {
			Image string `json:"image"`
			Color string `json:"color"`
		} `json:"lowQuality"`
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
	PlanId        primitive.ObjectID `json:"planId"`
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
				"planId":            parentstruct_in.PlanId,
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
type postadd struct {
	Images []struct {
		Image string `json:"image"`
		Color string `json:"color"`
	} `json:"images"`
	Category primitive.ObjectID `json:"category"`
	Country  string             `json:"country"`
	// CreationTimestamp time.Time `json:"creationTimestamp"`
	Currency string `json:"currency"`
	// NumberRatings     int                `json:"numberRatings"`
	OwnerID primitive.ObjectID `json:"ownerId"`
	Price   int                `json:"price"`
	Qty     int                `json:"qty"`
	// Rating            float64            `json:"rating"`
	RemainingQty int `json:"remainingQty"`
	// Status            string             `json:"status"`
	SubCategory primitive.ObjectID `json:"subCategory"`
	Title       string             `json:"title"`
	PlanID      primitive.ObjectID `json:"planId"`
	// ParentID          primitive.ObjectID `json:"parentId"`
	AvailableColor []primitive.ObjectID `json:"availableColor"`
	Feature        []struct {
		Name string `json:"name"`
	} `json:"feature"`
	Gender       primitive.ObjectID `json:"gender"`
	HasDimension bool               `json:"hasDimension"`
	Name         string             `json:"name"`
	Size         struct {
		AvailableSize []primitive.ObjectID `json:"availableSize"`
	} `json:"size"`
	Fabric           string `json:"fabric"`
	Careinstructions string `josn:"careInstructions"`
	Dimension        struct {
		Length struct {
			Value float32 `json:"value"`
			Unit  string  `json:"unit"`
		} `josn:"lenght"`
		Widht struct {
			Value float32 `json:"value"`
			Unit  string  `json:"unit"`
		} `json:"width"`
	} `json:"dimension"`
}
type add_img_itme_result struct {
	Status  int    `json:"status"`
	Message string `json:"message "`
	Data    PID    `json:"data"`
}
type PID struct {
	ParentID interface{} `json:"parentId"`
}
type adsRemaining1 struct {
	AdsRemaining int `json:"adsRemaining"`
}

func Add_item_wrt_category(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var strcutinit postadd
	err := json.NewDecoder(req.Body).Decode(&strcutinit)
	if err != nil {
		panic(err)
	}
	insertdat := bson.M{
		"images":            strcutinit.Images,
		"category":          strcutinit.Category,
		"country":           strcutinit.Country,
		"creationTimestamp": time.Now(),
		"currency":          strcutinit.Currency,
		"numberRatings":     5,
		"ownerId":           strcutinit.OwnerID,
		"price":             strcutinit.Price,
		"qty":               strcutinit.Qty,
		"rating":            0,
		"remainingQty":      strcutinit.Qty,
		"status":            "pending",
		"subCategory":       strcutinit.SubCategory,
		"title":             strcutinit.Title,
		"planId":            strcutinit.PlanID,
	}
	coll := docking.PakTradeDb.Collection("items-parent")
	coll1 := docking.PakTradeDb.Collection("Mammalas_login")

	// // // find user reming ads

	var result_rem adsRemaining1
	filter := bson.M{"_id": strcutinit.OwnerID}

	err1 := coll1.FindOne(context.TODO(), filter).Decode(&result_rem)
	if err1 != nil {
		fmt.Println("errror retrieving user userid : ")
	}
	if result_rem.AdsRemaining > 0 {
		responceid, err3 := coll.InsertOne(context.TODO(), insertdat)
		if err3 != nil {
			fmt.Print(err3)
		}
		/////////////////// update Remenig  Ads value
		minus := result_rem.AdsRemaining - 1
		_, err := coll1.UpdateOne(
			context.TODO(),
			bson.M{"_id": strcutinit.OwnerID},
			bson.D{
				{Key: "$set", Value: bson.M{
					"adsRemaining": minus}}},
		)
		///////////////
		if strcutinit.HasDimension == false {
			inster_cloths := bson.M{

				"parentId":         responceid.InsertedID,
				"availableColor":   strcutinit.AvailableColor,
				"feature":          strcutinit.Feature,
				"gender":           strcutinit.Gender,
				"hasDimension":     strcutinit.HasDimension,
				"name":             strcutinit.Name,
				"fabric":           strcutinit.Fabric,
				"careInstructions": strcutinit.Careinstructions,
				"size": bson.M{
					"availableSize": strcutinit.Size.AvailableSize,
					"sizeChart":     "",
				},
			}

			if strcutinit.Category.Hex() == "63a9a76fd38789473ba919e6" {
				coll1 := docking.PakTradeDb.Collection("cloths")
				_, err4 := coll1.InsertOne(context.TODO(), inster_cloths)
				if err4 != nil {
					fmt.Print(err4)
				}

			} else if strcutinit.Category.Hex() == "63bdb52116cccb9bb8b48388" {
				coll1 := docking.PakTradeDb.Collection("item-mart")
				_, err4 := coll1.InsertOne(context.TODO(), responceid.InsertedID)
				if err4 != nil {
					fmt.Print(err4)
				}

			}
		} else {
			inster_cloths := bson.M{

				"parentId":         responceid.InsertedID,
				"availableColor":   strcutinit.AvailableColor,
				"feature":          strcutinit.Feature,
				"gender":           strcutinit.Gender,
				"hasDimension":     strcutinit.HasDimension,
				"name":             strcutinit.Name,
				"fabric":           strcutinit.Fabric,
				"careInstructions": strcutinit.Careinstructions,
				"dimension": bson.M{
					"length": bson.M{
						"unit":  strcutinit.Dimension.Length.Unit,
						"value": strcutinit.Dimension.Length.Value,
					},
					"width": bson.M{
						"unit":  strcutinit.Dimension.Widht.Unit,
						"value": strcutinit.Dimension.Widht.Value,
					}},
				"size": bson.M{
					"availableSize": strcutinit.Size.AvailableSize,
					"sizeChart":     "",
				},
			}

			if strcutinit.Category.Hex() == "63a9a76fd38789473ba919e6" {
				coll1 := docking.PakTradeDb.Collection("cloths")
				_, err4 := coll1.InsertOne(context.TODO(), inster_cloths)
				if err4 != nil {
					fmt.Print(err4)
				}

			} else if strcutinit.Category.Hex() == "63bdb52116cccb9bb8b48388" {
				coll1 := docking.PakTradeDb.Collection("item-mart")
				_, err4 := coll1.InsertOne(context.TODO(), responceid.InsertedID)
				if err4 != nil {
					fmt.Print(err4)
				}

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
		////////////// End of image upload_insert
	} else {
		var results add_img_itme_result
		results.Status = http.StatusOK
		results.Message = "limit exceeded"
		output, err := json.MarshalIndent(results, "", "    ")
		if err != nil {
			panic(err)

		}
		fmt.Fprintf(w, "%s\n", output)
	}
}

// //////////// item serch by id
func Serch_item_by_id(w http.ResponseWriter, r *http.Request) {
	coll := docking.PakTradeDb.Collection("cloths")

	var id_get delete_id
	err := json.NewDecoder(r.Body).Decode(&id_get)
	handleError(err)
	// objectId, err := primitive.ObjectIDFromHex(id_get.Item_id)
	// handleError(err)
	mongoquery := []bson.D{
		bson.D{{Key: "$match", Value: bson.D{{Key: "parentId", Value: id_get.Item_id}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "items-parent"},
					{"localField", "parentId"},
					{"foreignField", "_id"},
					{"as", "parent"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "color"},
					{"localField", "availableColor"},
					{"foreignField", "_id"},
					{"as", "color"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "size"},
					{"localField", "size.availableSize"},
					{"foreignField", "_id"},
					{"as", "size"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "size_chart"},
					{"localField", "size.sizeChart"},
					{"foreignField", "_id"},
					{"as", "sizeChart"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "sub_category"},
					{"localField", "parent.subCategory"},
					{"foreignField", "_id"},
					{"as", "sub_cat"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "categories"},
					{"localField", "parent.category"},
					{"foreignField", "_id"},
					{"as", "cat"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "plans"},
					{"localField", "parent.planId"},
					{"foreignField", "_id"},
					{"as", "plans"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "Mammalas_login"},
					{"localField", "parent.ownerId"},
					{"foreignField", "_id"},
					{"as", "owner"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "tier"},
					{"localField", "plans.tierId"},
					{"foreignField", "_id"},
					{"as", "tier"},
				},
			},
		},
		bson.D{
			{"$set",
				bson.D{
					{"parent", bson.D{{"$first", "$parent"}}},
					{"category", bson.D{{"$first", "$cat"}}},
					{"subCategory", bson.D{{"$first", "$sub_cat"}}},
					{"plans", bson.D{{"$first", "$plans"}}},
					{"tier", bson.D{{"$first", "$tier"}}},
					{"owner", bson.D{{"$first", "$owner"}}},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"name", "$name"},
					{"feature", "$feature"},
					{"images", "$parent.images"},
					{"_id", "$_id"},
					{"parentId", "$parent._id"},
					{"category", "$category"},
					{"sub_category", "$subCategory"},
					{"gender", "$gender"},
					{"price", "$parent.price"},
					{"qty", "$parent.qty"},
					{"remaining_qty", "$parent.remainingQty"},
					{"status", "$parent.status"},
					{"country", "$parent.country"},
					{"currency", "$parent.currency"},

					{"color", "$color"},
					{"size", "$size"},
					{"title", "$parent.title"},
					{"ownerName", "$owner.displayName"},
					{"plan",
						bson.D{
							{"price", "$plans.price"},
							{"name", "$tier.name"},
							{"order", "$tier.order"},
							{"adDuration", "$plans.ad_duration"},
						},
					},
					{"dimension",
						bson.D{
							{"width",
								bson.D{
									{"unit", "$dimension.width.unit"},
									{"value", "$dimension.width.value"},
								},
							},
							{"length",
								bson.D{
									{"unit", "$dimension.length.unit"},
									{"value", "$dimension.length.value"},
								},
							},
						},
					},
				},
			},
		},
	}

	cursor, err := coll.Aggregate(context.TODO(), mongoquery)

	if err != nil {
		log.Fatal(err)
	}

	var results respone_st
	var resp1 AutoGenerated
	if cursor.Next(context.TODO()) {
		err := cursor.Decode(&resp1)
		if err != nil {
			log.Fatal(err)
		}
	}

	if resp1.Currency != "" {
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

/////// end item serch by id
