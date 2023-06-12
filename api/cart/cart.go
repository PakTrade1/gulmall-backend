package cart

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	docking "pak-trade-go/Docking"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartMammals struct {
	Mammal_id       primitive.ObjectID   `json:"userId"`
	Item_id         primitive.ObjectID   `json:"itemId"`
	Quantity        int                  `json:"quantity"`
	Price           int                  `json:"price"`
	Payement_method primitive.ObjectID   `json:"payementMethod"`
	Color_id        []primitive.ObjectID `json:"colorId"`
	Size_id         []primitive.ObjectID `json:"sizeId"`
	Discount        string               `json:"discount"`
	SellerInfo      primitive.ObjectID   `json:"sellerInfo"`
}

type Resp_insert struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Id      interface{} `json:"id"`
}
type Get_qty struct {
	Qty int `json:"qty"`
}

func Cart_insertone(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var cart_init CartMammals
	err := json.NewDecoder(req.Body).Decode(&cart_init)
	if err != nil {
		panic(err)
	}
	Total_price_cart := cart_init.Price * cart_init.Quantity
	mongo_query := bson.M{
		"mammalId":       cart_init.Mammal_id,
		"itemId":         cart_init.Item_id,
		"deliveryStatus": "pending",
		"orderDate":      time.Now(),
		"colorId":        cart_init.Color_id,
		"sizeId":         cart_init.Size_id,
		"quantity":       cart_init.Quantity,
		"price":          cart_init.Price,
		"discount":       cart_init.Discount,
		"payementMethod": cart_init.Payement_method,
		"totalPrice":     Total_price_cart,
		"sellerInfo":     cart_init.SellerInfo,
	}

	coll := docking.PakTradeDb.Collection("cart_mammals")
	coll1 := docking.PakTradeDb.Collection("items-parent")
	var qty Get_qty
	filter := bson.M{"_id": cart_init.Item_id}

	err12 := coll1.FindOne(context.TODO(), filter).Decode(&qty)
	if err12 != nil {
		log.Fatal(err12)
	}
	// // insert a user

	inset, err3 := coll.InsertOne(context.TODO(), mongo_query)
	if err3 != nil {
		fmt.Fprintf(w, "%s\n", err3)
	}
	Qty_minus := qty.Qty - cart_init.Quantity
	_, err1 := coll1.UpdateOne(
		context.TODO(),
		bson.M{"_id": cart_init.Item_id},
		bson.D{
			{Key: "$set", Value: bson.M{
				"qty": Qty_minus,
			}},
		},
	)
	if err1 != nil {
		log.Fatal(err1)
	}

	var results Resp_insert
	if inset != nil {
		results.Status = http.StatusOK
		results.Message = "success"
		results.Id = inset.InsertedID

	} else {
		results.Message = "decline"

	}

	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)

	}

	fmt.Fprintf(w, "%s\n", output)

}

type Item struct {
	ID         primitive.ObjectID `bson:"_id"`
	OrderDate  time.Time          `json:"orderDate"`
	SellerInfo string             `json:"sellerInfo"`
	Qty        int                `json:"qty"`
	Size       []struct {
		ID   primitive.ObjectID `bson:"_id"`
		Name string             `json:"name"`
	} `json:"size"`
	Color []struct {
		ID     string `bson:"_id"`
		CssHex string `json:"cssHex"`
		Name   string `json:"name"`
	} `json:"color"`
	TotalPrice int    `json:"totalPrice"`
	Discount   string `json:"discount"`
	ItemName   string `json:"itemName"`
	Images     []struct {
		Image string `json:"image"`
		Color string `json:"color"`
	} `json:"images"`
	ItemPrice int `json:"itemPrice"`
}

type Size struct {
	ID   primitive.ObjectID `json:"_id"`
	Name string             `json:"name"`
}

type Color struct {
	ID     primitive.ObjectID `json:"_id"`
	CssHex string             `json:"cssHex"`
	Name   string             `json:"name"`
}
type UserID struct {
	UserID primitive.ObjectID `json:"userId"`
}

func Cart_getall(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var cart_init UserID

	err := json.NewDecoder(req.Body).Decode(&cart_init)
	if err != nil {
		panic(err)
	}

	coll := docking.PakTradeDb.Collection("cart_mammals")

	mongoQ := bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"userId", cart_init.UserID},
					{"deliveryStatus", "pending"},
				},
			},
		},
		bson.D{{"$unwind", bson.D{{"path", "$colorId"}}}},
		bson.D{{"$unwind", bson.D{{"path", "$sizeId"}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "color"},
					{"localField", "colorId"},
					{"foreignField", "_id"},
					{"as", "color"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "size"},
					{"localField", "sizeId"},
					{"foreignField", "_id"},
					{"as", "size"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "items-parent"},
					{"localField", "itemId"},
					{"foreignField", "_id"},
					{"as", "items"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "Mammalas_login"},
					{"localField", "sellerInfo"},
					{"foreignField", "_id"},
					{"as", "seller"},
				},
			},
		},
		bson.D{
			{"$set",
				bson.D{
					{"seller", bson.D{{"$first", "$seller"}}},
					{"items", bson.D{{"$first", "$items"}}},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"orderDate", "$orderDate"},
					{"sellerInfo", "$seller.displayName"},
					{"qty", "$quantity"},
					{"size", "$size"},
					{"color", "$color"},
					{"totalPrice", "$totalPrice"},
					{"discount", "$discount"},
					{"itemName", "$items.title"},
					{"images", "$items.images"},
					{"itemPrice", "$items.price"},
				},
			},
		},
	}

	cursor, err := coll.Aggregate(context.Background(), mongoQ)
	if err != nil {
		panic(err)
	}

	var results []Item
	for cursor.Next(context.TODO()) {
		var abc Item
		cursor.Decode(&abc)
		results = append(results, abc)

	}

	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)

	}
	fmt.Fprintf(w, "%s\n", output)

}
