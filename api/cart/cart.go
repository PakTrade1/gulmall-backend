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
	Orders []struct {
		Mammal_id       primitive.ObjectID `json:"user_id"`
		Item_id         primitive.ObjectID `json:"item_id"`
		Quantity        int                `json:"quantity"`
		Payement_method primitive.ObjectID `json:"payement_method"`
		Color_id        primitive.ObjectID `json:"color_id"`
		Size_id         primitive.ObjectID `json:"size_id"`
		SellerInfo      primitive.ObjectID `json:"seller_info"`
		Price           float32            `json:"price"`
		Discount        string             `json:"discount"`
		Total_price     float32            `json:"total_price"`
		Currency        string             `json:"currency"`
		Rem             int                `json:"items_remaining_quantity"`
	} `json:"orders"`
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

	var cart_init CartMammals
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

	for i := 0; i < len(cart_init.Orders); i++ {

		// // insert a user

		mongo_query := bson.M{
			"user_id":         cart_init.Orders[i].Mammal_id,
			"item_id":         cart_init.Orders[i].Item_id,
			"delivery_status": "pending",
			"orderDate":       time.Now(),
			"color_id":        cart_init.Orders[i].Color_id,
			"size_id":         cart_init.Orders[i].Size_id,
			"quantity":        cart_init.Orders[i].Quantity,
			"price":           cart_init.Orders[i].Price,
			"discount":        cart_init.Orders[i].Discount,
			"payement_method": cart_init.Orders[i].Payement_method,
			"total_price":     cart_init.Orders[i].Total_price,
			"seller_info":     cart_init.Orders[i].SellerInfo,
			"currency":        cart_init.Orders[i].Currency,
		}

		inset_data, err3 := coll.InsertOne(context.TODO(), mongo_query)
		if err3 != nil {
			fmt.Fprintf(w, "%s\n", err3)
		}
		inset.InsertedID = inset_data.InsertedID
		Qty_minus := cart_init.Orders[i].Rem
		_, err1 := coll1.UpdateOne(
			context.TODO(),
			bson.M{"_id": cart_init.Orders[i].Item_id},
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
	ID          primitive.ObjectID `bson:"_id"`
	OrderDate   string             `json:"order_date"`
	Seller_info string             `json:"seller_info"`
	Qty         int                `json:"qty"`
	Size        []struct {
		ID   primitive.ObjectID `bson:"_id"`
		Name string             `json:"name"`
	} `json:"size"`
	Color []struct {
		ID     string `bson:"_id"`
		CssHex string `json:"cssHex"`
		Name   string `json:"name"`
	} `json:"color"`
	Total_price int    `json:"total_price"`
	Discount    string `json:"discount"`
	Item_name   string `json:"item_name"`
	Images      []struct {
		Image string `json:"image"`
		Color string `json:"color"`
	} `json:"images"`
	Item_price int `json:"item_price"`
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
	UserID primitive.ObjectID `json:"user_id"`
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
					{"mammal_id", cart_init.UserID},
					{"delivery_status", "pending"},
				},
			},
		},
		bson.D{{"$unwind", bson.D{{"path", "$color_id"}}}},
		bson.D{{"$unwind", bson.D{{"path", "$size_id"}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "color"},
					{"localField", "color_id"},
					{"foreignField", "_id"},
					{"as", "color"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "size"},
					{"localField", "size_id"},
					{"foreignField", "_id"},
					{"as", "size"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "items-parent"},
					{"localField", "item_id"},
					{"foreignField", "_id"},
					{"as", "items"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "Mammalas_login"},
					{"localField", "seller_info"},
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
					{"order_date", "$orderDate"},
					{"seller_info", "$seller.displayName"},
					{"qty", "$quantity"},
					{"size", "$size"},
					{"color", "$color"},
					{"total_price", "$total_price"},
					{"discount", "$discount"},
					{"item_name", "$items.title"},
					{"images", "$items.images"},
					{"item_price", "$items.price"},
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
