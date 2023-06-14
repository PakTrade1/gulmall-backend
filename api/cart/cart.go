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
	Mammal_id       primitive.ObjectID `json:"user_id"`
	Item_id         primitive.ObjectID `json:"item_id"`
	Quantity        int                `json:"quantity"`
	Payement_method primitive.ObjectID `json:"payement_method"`
	Color_id        primitive.ObjectID `json:"color_id"`
	Size_id         primitive.ObjectID `json:"size_id"`
	SellerInfo      primitive.ObjectID `json:"seller_info"`
}

type Resp_insert struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Id      interface{} `json:"id"`
}
type Get_qty struct {
	Qty      int     `json:"qty"`
	Price    float32 `json:"price"`
	Discount string  `json:"discount"`
}

func Cart_insertone_fashion(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var cart_init []CartMammals
	err := json.NewDecoder(req.Body).Decode(&cart_init)
	if err != nil {
		panic(err)
	}

	coll := docking.PakTradeDb.Collection("cart_mammals")
	coll1 := docking.PakTradeDb.Collection("items-parent")

	inset := struct {
		InsertedID interface{}
	}{
		InsertedID: "",
	}

	for i := 0; i < len(cart_init); i++ {

		filter := bson.M{"_id": cart_init[i].Item_id}
		var qty Get_qty

		err12 := coll1.FindOne(context.TODO(), filter).Decode(&qty)
		if err12 != nil {
			log.Fatal(err12)
		}
		// // insert a user

		Price := qty.Price
		Discount := qty.Discount
		Total_price_cart := int(Price) * cart_init[i].Quantity
		mongo_query := bson.M{
			"mammal_id":       cart_init[i].Mammal_id,
			"item_id":         cart_init[i].Item_id,
			"delivery_status": "pending",
			"orderDate":       time.Now(),
			"color_id":        cart_init[i].Color_id,
			"size_id":         cart_init[i].Size_id,
			"quantity":        cart_init[i].Quantity,
			"price":           Price,
			"discount":        Discount,
			"payement_method": cart_init[i].Payement_method,
			"total_price":     Total_price_cart,
			"seller_info":     cart_init[i].SellerInfo,
		}

		inset_data, err3 := coll.InsertOne(context.TODO(), mongo_query)
		if err3 != nil {
			fmt.Fprintf(w, "%s\n", err3)
		}
		inset.InsertedID = inset_data.InsertedID
		Qty_minus := qty.Qty - cart_init[i].Quantity
		_, err1 := coll1.UpdateOne(
			context.TODO(),
			bson.M{"_id": cart_init[i].Item_id},
			bson.D{
				{Key: "$set", Value: bson.M{
					"qty": Qty_minus,
				}},
			},
		)
		if err1 != nil {
			log.Fatal(err1)
		}
	}
	var results Resp_insert
	if inset.InsertedID != nil {
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
