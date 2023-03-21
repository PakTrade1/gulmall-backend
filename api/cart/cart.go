package cart

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

type CartMammals struct {
	Mammal_id       primitive.ObjectID `bson:"mammal_id,omitempty"`
	Item_id         primitive.ObjectID `bson:"item_id,omitempty"`
	Quantity        int                `json:"quantity"`
	Price           int                `json:"price"`
	Payement_method primitive.ObjectID `bson:"payement_method,omitempty"`
	color_id        primitive.ObjectID `bson:"_id,omitempty"`
	size_id         primitive.ObjectID `bson:"size_id,omitempty"`
	Discount        string             `json:"discount"`
	Total_price     float32            `json:"total_price"`
}

func Cart_getall(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	coll := docking.PakTradeDb.Collection("cart_mammals")

	cursor, err := coll.Find(context.Background(), coll)
	if err != nil {
		panic(err)
	}

	var results []CartMammals
	for cursor.Next(context.TODO()) {
		var abc CartMammals
		cursor.Decode(&abc)
		results = append(results, abc)

	}

	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)

	}
	fmt.Fprintf(w, "%s\n", output)

}

type resp_insert struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Id      interface{} `json:"id"`
}

func Cart_insertone(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	var cart_init CartMammals
	err := json.NewDecoder(req.Body).Decode(&cart_init)
	if err != nil {
		panic(err)
	}
	// mongo query
	Total_price_cart := cart_init.Price * cart_init.Quantity
	mongo_query := bson.M{
		"mammal_id":       cart_init.Mammal_id,
		"item_id":         cart_init.Item_id,
		"color_id":        cart_init.color_id,
		"size_id":         cart_init.size_id,
		"quantity":        cart_init.Quantity,
		"price":           cart_init.Price,
		"discount":        cart_init.Discount,
		"payement_method": cart_init.Payement_method,
		"total_price":     Total_price_cart,
	}

	coll := docking.PakTradeDb.Collection("cart_mammals")

	// // // insert a user

	inset, err3 := coll.InsertOne(context.TODO(), mongo_query)
	if err3 != nil {
		fmt.Fprintf(w, "%s\n", err3)
	}
	var results resp_insert
	if inset != nil {
		results.Status = http.StatusOK
		results.Message = "success"

	} else {
		results.Message = "decline"

	}

	results.Id = inset.InsertedID
	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)

	}

	fmt.Fprintf(w, "%s\n", output)

}

type resp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
type cart_responce struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Price       int                `json:"price"`
	Quantity    int                `json:"quantity"`
	Total_price int                `json:"total_price"`

	Mamals struct {
		ID   primitive.ObjectID `bson:"_id,omitempty"`
		Name struct {
			Firt_name string `json:"firt_name"`
			Last_name string `json:"last_name"`
		} `json:"name"`
		Address []struct {
			// Home_address struct {
			// 	Address  string `json:"address"`
			// 	Country  string `json:"country"`
			// 	City     string `json:"city"`
			// 	Province string `json:"province"`
			// 	ZipCode  int    `json:"zip_code"`
			// } `json:"home_address"`
			Shipping_address struct {
				Address  string `json:"address"`
				Country  string `json:"country"`
				City     string `json:"city"`
				Province string `json:"province"`
				ZipCode  int    `json:"zip_code"`
			} `json:"shipping_address"`
		} `json:"address"`
		Email string `json:"email"`
		Phone string `json:"phone"`
		// Profile string `json:"profile"`
	} `json:"user_details"`
	Item struct {
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
		// Price int `json:"price"`
	} `json:"item_details"`
	Color struct {
		//CSSHex string `json:"cssHex"`
		Name struct {
			En string `json:"en"`
			Ar string `json:"ar"`
		} `json:"name"`
	} `json:"color"`
	Size struct {
		Size string `json:"size"`
	} `json:"size"`
	Payment struct {
		Payment_status string `json:"payment_status"`
	} `json:"payment"`
	// Total_price float32 `json:"total_price"`
}
type respone_struct struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    []cart_responce `json:"data"`
}

func Get_cart_all_with_id_data(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// var Status_req status_req
	// err := json.NewDecoder(req.Body).Decode(&Status_req)
	// if err != nil {
	// 	panic(err)
	// }

	coll := docking.PakTradeDb.Collection("cart_mammals")

	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver
	ctx := context.TODO()

	mongo_query := bson.A{
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "color"},
					{Key: "localField", Value: "color_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "color"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "cloths"},
					{Key: "localField", Value: "item_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "item"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "size"},
					{Key: "localField", Value: "size_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "size"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "mammals"},
					{Key: "localField", Value: "mammal_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "mamals"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "mamals_paymentInfo"},
					{Key: "localField", Value: "payement_method"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "payment"},
				},
			},
		},
		bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key: "color", Value: bson.D{{Key: "$first", Value: "$color"}}},
					{Key: "item", Value: bson.D{{Key: "$first", Value: "$item"}}},
					{Key: "mamals", Value: bson.D{{Key: "$first", Value: "$mamals"}}},
					{Key: "size", Value: bson.D{{Key: "$first", Value: "$size"}}},
					{Key: "payment", Value: bson.D{{Key: "$first", Value: "$payment"}}},
				},
			},
		},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: "$_id"},
					{Key: "price", Value: "$price"},
					{Key: "total_price", Value: "$total_price"},
					{Key: "quantity", Value: "$quantity"},
					{Key: "mamals", Value: "$mamals"},
					{Key: "item", Value: "$item"},
					{Key: "color", Value: "$color"},
					{Key: "size", Value: "$size"},
					{Key: "payment", Value: "$payment"},
				},
			},
		},
	}

	// Open an aggregation cursor
	cursor, err := coll.Aggregate(ctx, mongo_query)
	if err != nil {
		log.Fatal(err)
	}
	var results respone_struct
	var resp1 []cart_responce
	for cursor.Next(context.TODO()) {

		var xy cart_responce
		cursor.Decode(&xy)
		resp1 = append(resp1, xy)

	}
	if cursor != nil {
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

type cart_id_req struct {
	Cart_id string `json:"cart_id"`
}

func Get_cart_with_id(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var Id_req cart_id_req
	err := json.NewDecoder(req.Body).Decode(&Id_req)
	if err != nil {
		panic(err)
	}

	coll := docking.PakTradeDb.Collection("cart_mammals")

	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver
	ctx := context.TODO()
	objectIDS, _ := primitive.ObjectIDFromHex(Id_req.Cart_id)

	mongo_query := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: objectIDS}}}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "color"},
					{Key: "localField", Value: "color_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "color"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "cloths"},
					{Key: "localField", Value: "item_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "item"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "size"},
					{Key: "localField", Value: "size_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "size"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "mammals"},
					{Key: "localField", Value: "mammal_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "mamals"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "mamals_paymentInfo"},
					{Key: "localField", Value: "payement_method"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "payment"},
				},
			},
		},
		bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key: "color", Value: bson.D{{Key: "$first", Value: "$color"}}},
					{Key: "item", Value: bson.D{{Key: "$first", Value: "$item"}}},
					{Key: "mamals", Value: bson.D{{Key: "$first", Value: "$mamals"}}},
					{Key: "size", Value: bson.D{{Key: "$first", Value: "$size"}}},
					{Key: "payment", Value: bson.D{{Key: "$first", Value: "$payment"}}},
				},
			},
		},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: "$_id"},
					{Key: "price", Value: "$price"},
					{Key: "total_price", Value: "$total_price"},
					{Key: "quantity", Value: "$quantity"},
					{Key: "mamals", Value: "$mamals"},
					{Key: "item", Value: "$item"},
					{Key: "color", Value: "$color"},
					{Key: "size", Value: "$size"},
					{Key: "payment", Value: "$payment"},
				},
			},
		},
	}

	// Open an aggregation cursor
	cursor, err := coll.Aggregate(ctx, mongo_query)
	if err != nil {
		log.Fatal(err)
	}
	var results respone_struct
	var resp1 []cart_responce
	for cursor.Next(context.TODO()) {

		var xy cart_responce
		cursor.Decode(&xy)
		resp1 = append(resp1, xy)
	}
	if cursor != nil {
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

// update cart w.r.t changess like color size and qty
type cart_id_data_validation struct {
	Mammal_id       primitive.ObjectID `json:"mammal_id"`
	Item_id         primitive.ObjectID `json:"item_id"`
	Quantity        int                `json:"quantity"`
	Price           int                `json:"price"`
	Payement_method primitive.ObjectID `bson:"payement_method"`
	color_id        primitive.ObjectID `json:"color_id"`
	size_id         primitive.ObjectID `json:"size_id"`
	Discount        string             `json:"discount"`
	Total_price     float32            `json:"total_price"`
}
type cart_id_req_update struct {
	Cart_id  string             `json:"cart_id"`
	color_id primitive.ObjectID `json:"color_id"`
	size_id  primitive.ObjectID `json:"size_id"`
}
type resp_update struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Id      string `json:"id"`
}

func Update_cart_all_with_id_data(w http.ResponseWriter, req *http.Request) {
	var search1 cart_id_req_update
	err := json.NewDecoder(req.Body).Decode(&search1)
	if err != nil {
		panic(err)
	}
	coll := docking.PakTradeDb.Collection("cart_mammals")
	objectIDS, _ := primitive.ObjectIDFromHex(search1.Cart_id)

	var result cart_id_data_validation
	filter := bson.M{"_id": objectIDS}

	err1 := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err1 != nil {
		fmt.Println("errror retrieving user userid : " + objectIDS.Hex())
	}

	// end findOne
	// update data if add to cart thing not change only qty changed

	if result.color_id == search1.color_id && result.size_id == search1.size_id {
		coll := docking.PakTradeDb.Collection("cart_mammals")
		objectIDS, _ := primitive.ObjectIDFromHex(search1.Cart_id)

		//total price updete and qty
		qty_update := result.Quantity + 1
		total_price_update := result.Price * qty_update
		result1, err := coll.UpdateOne(
			context.TODO(),
			bson.M{"_id": objectIDS},
			bson.D{
				{Key: "$set", Value: bson.M{
					"quantity":    qty_update,
					"total_price": total_price_update,
				}},
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		//end update
		var results resp_update
		if result1 != nil {
			results.Status = http.StatusOK
			results.Message = "success"

		} else {
			results.Message = "decline"

		}

		results.Id = search1.Cart_id
		output, err := json.MarshalIndent(results, "", "    ")
		if err != nil {
			panic(err)

		}

		fmt.Fprintf(w, "%s\n", output)

	} else {
		if search1.color_id != result.color_id && search1.size_id == result.size_id {
			mongo_query := bson.M{
				"mammal_id":       result.Mammal_id,
				"item_id":         result.Item_id,
				"color_id":        search1.color_id,
				"size_id":         result.size_id,
				"quantity":        1,
				"price":           result.Price,
				"discount":        result.Discount,
				"payement_method": result.Payement_method,
				"total_price":     1 * result.Quantity,
			}

			coll := docking.PakTradeDb.Collection("cart_mammals")
			inset, err3 := coll.InsertOne(context.TODO(), mongo_query)
			if err3 != nil {
				fmt.Fprintf(w, "%s\n", err3)
			}
			var results resp_insert
			if inset != nil {
				results.Status = http.StatusOK
				results.Message = "success"

			} else {
				results.Message = "decline"

			}

			results.Id = inset.InsertedID
			output, err := json.MarshalIndent(results, "", "    ")
			if err != nil {
				panic(err)

			}

			fmt.Fprintf(w, "%s\n", output)

		} else if search1.size_id != result.size_id && search1.color_id == result.color_id {
			mongo_query := bson.M{
				"mammal_id":       result.Mammal_id,
				"item_id":         result.Item_id,
				"color_id":        result.color_id,
				"size_id":         search1.size_id,
				"quantity":        1,
				"price":           result.Price,
				"discount":        result.Discount,
				"payement_method": result.Payement_method,
				"total_price":     1 * result.Quantity,
			}

			coll := docking.PakTradeDb.Collection("cart_mammals")

			// // // insert a user

			inset, err3 := coll.InsertOne(context.TODO(), mongo_query)
			if err3 != nil {
				fmt.Fprintf(w, "%s\n", err3)
			}
			var results resp_insert
			if inset != nil {
				results.Status = http.StatusOK
				results.Message = "success"

			} else {
				results.Message = "decline"

			}

			results.Id = inset.InsertedID
			output, err := json.MarshalIndent(results, "", "    ")
			if err != nil {
				panic(err)

			}

			fmt.Fprintf(w, "%s\n", output)

		} else if search1.size_id != result.color_id && search1.color_id != result.color_id {

			mongo_query := bson.M{
				"mammal_id":       result.Mammal_id,
				"item_id":         result.Item_id,
				"color_id":        search1.color_id,
				"size_id":         search1.size_id,
				"quantity":        1,
				"price":           result.Price,
				"discount":        result.Discount,
				"payement_method": result.Payement_method,
				"total_price":     1 * result.Quantity,
			}

			coll := docking.PakTradeDb.Collection("cart_mammals")

			// // // insert a user

			inset, err3 := coll.InsertOne(context.TODO(), mongo_query)
			if err3 != nil {
				fmt.Fprintf(w, "%s\n", err3)
			}
			var results resp_insert
			if inset != nil {
				results.Status = http.StatusOK
				results.Message = "success"

			} else {
				results.Message = "decline"

			}

			results.Id = inset.InsertedID
			output, err := json.MarshalIndent(results, "", "    ")
			if err != nil {
				panic(err)

			}

			fmt.Fprintf(w, "%s\n", output)

		}
	}

}

type delte_status struct {
	Cart_id string `json:"cart_id"`
}
type resp_delete struct {
	Status        int    `json:"status"`
	Message       string `json:"message"`
	Status_delete string `json:"status_delete"`
}

func Cart_delete(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var strcutinit delte_status
	err := json.NewDecoder(req.Body).Decode(&strcutinit)
	if err != nil {
		panic(err)
	}
	coll := docking.PakTradeDb.Collection("cart_mammals")
	objectIDS, _ := primitive.ObjectIDFromHex(strcutinit.Cart_id)
	// fmt.Print(objectIDS)

	//end update
	result, err := coll.DeleteOne(context.TODO(), bson.M{"_id": objectIDS}, nil)
	if err != nil {
		log.Fatal(err)
	}

	var results resp_delete
	if result.DeletedCount == 1 {
		results.Status = http.StatusOK
		results.Message = "success"
		results.Status_delete = "delete successfully"

	} else {
		results.Message = "decline"
		results.Status = http.StatusBadRequest
		results.Status_delete = "no record found"

	}

	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)

	}

	fmt.Fprintf(w, "%s\n", output)

}

type cart_id_data struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Mammal_id       primitive.ObjectID `json:"mammal_id"`
	Item_id         primitive.ObjectID `json:"item_id"`
	Quantity        int                `json:"quantity"`
	Price           int                `json:"price"`
	Payement_method primitive.ObjectID `bson:"payement_method"`
	Color_id        primitive.ObjectID `json:"color_id"`
	Size_id         primitive.ObjectID `json:"size_id"`
	Discount        string             `json:"discount"`
	Total_price     float32            `json:"total_price"`
}
type resp_update1 struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Id      interface{} `json:"id"`
}
type price struct {
	Price float32 `json:"price"`
}
type paymanet_status struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	Payment_status string `json:"payment_status"`
}

func Update_cart(w http.ResponseWriter, req *http.Request) {
	mammal_id := req.Header.Get("mammal_id")
	var search1 cart_id_data
	err := json.NewDecoder(req.Body).Decode(&search1)
	if err != nil {
		panic(err)
	}
	objectIDS, _ := primitive.ObjectIDFromHex(mammal_id)
	var Payment_status paymanet_status
	coll2 := docking.PakTradeDb.Collection("mamals_paymentInfo")
	filter2 := bson.M{"user_id": mammal_id}
	err3 := coll2.FindOne(context.TODO(), filter2).Decode(&Payment_status)
	if err3 != nil {
		log.Fatal(err3)
	}
	//fmt.Println(Payment_status.ID, "abbasi")

	var Price price
	coll1 := docking.PakTradeDb.Collection("cloths")
	filter1 := bson.M{"_id": search1.Item_id}
	err2 := coll1.FindOne(context.TODO(), filter1).Decode(&Price)
	if err2 != nil {
		log.Fatal(err2)
	}
	coll := docking.PakTradeDb.Collection("cart_mammals")
	var result cart_id_data
	filter := bson.M{"mammal_id": objectIDS, "item_id": search1.Item_id, "color_id": search1.Color_id, "size_id": search1.Size_id}
	err1 := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err1 != nil {
		// fmt.Println(err1)
		tp := int(Price.Price) * 1.
		mongo_query := bson.M{
			"mammal_id":       objectIDS,
			"item_id":         search1.Item_id,
			"color_id":        search1.Color_id,
			"size_id":         search1.Size_id,
			"quantity":        1,
			"price":           Price.Price,
			"discount":        "0%",
			"payement_method": Payment_status.ID,
			"total_price":     tp,
		}

		//coll := docking.PakTradeDb.Collection("cart_mammals")

		// // // insert a user

		inset, err3 := coll.InsertOne(context.TODO(), mongo_query)
		if err3 != nil {
			fmt.Fprintf(w, "%s\n", err3)
		}
		var results resp_insert
		if inset != nil {
			results.Status = http.StatusOK
			results.Message = "success"

		} else {
			results.Message = "decline"

		}

		results.Id = inset.InsertedID
		output, err := json.MarshalIndent(results, "", "    ")
		if err != nil {
			panic(err)

		}

		fmt.Fprintf(w, "%s\n", output)

	} else {
		qty_update := result.Quantity + 1
		total_price_update := result.Price * qty_update
		result1, err := coll.UpdateOne(
			context.TODO(),
			bson.M{"_id": result.ID},
			bson.D{
				{Key: "$set", Value: bson.M{
					"quantity":    qty_update,
					"total_price": total_price_update,
					"price":       Price.Price,
				}},
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		//end update
		var results resp_update1
		if result1 != nil {
			results.Status = http.StatusOK
			results.Message = "success"
		} else {
			results.Message = "decline"

		}
		//objectIDS, _ := primitive.ObjectIDFromHex(string(result.Id))

		results.Id = 1
		output, err := json.MarshalIndent(results, "", "    ")
		if err != nil {
			panic(err)

		}

		fmt.Fprintf(w, "%s\n", output)
	}

}
//Test
