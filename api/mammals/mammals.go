package mammals

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

type Mammals_serch struct {
	ID string `json:"_id"`
}
type Mammals_user struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	//ID string `json:"_id"`

	Name struct {
		Firt_name string `json:"firt_name" bson:"firt_name"`
		Last_name string `json:"last_name" bson:"last_name"`
	} `json:"name" bson:"name"`
	Address []struct {
		Home_address struct {
			Address  string `json:"address" bson:"address"`
			Country  string `json:"country"  bson:"country"`
			City     string `json:"city"  bson:"city"`
			Province string `json:"province" bson:"province"`
			Zip_code int    `json:"zip_code" bson:"zip_code"`
		} `json:"home_address,omitempty" bson:"home_address,omitempty"`
		Shipping_address struct {
			Address  string `json:"address" bson:"address"`
			Country  string `json:"country"  bson:"country"`
			City     string `json:"city"  bson:"city"`
			Province string `json:"province" bson:"province"`
			Zip_code int    `json:"zip_code" bson:"zip_code"`
		} `json:"shipping_address,omitempty" bson:"shipping_address,omitempty"`
	} `json:"address"`
	Email   string `json:"email" bson:"email"`
	Phone   string `json:"phone" bson:"phone"`
	Profile string `json:"profile" bson:"profile"`
}

type Mammals_user_update struct {
	//ID primitive.ObjectID `bson:"_id,omitempty"`
	ID string `json:"_id"`

	Name struct {
		Firt_name string `json:"firt_name" bson:"firt_name"`
		Last_name string `json:"last_name" bson:"last_name"`
	} `json:"name" bson:"name"`
	Address []struct {
		Home_address struct {
			Address  string `json:"address" bson:"address"`
			Country  string `json:"country"  bson:"country"`
			City     string `json:"city"  bson:"city"`
			Province string `json:"province" bson:"province"`
			Zip_code int    `json:"zip_code" bson:"zip_code"`
		} `json:"home_address,omitempty" bson:"home_address,omitempty"`
		Shipping_address struct {
			Address  string `json:"address" bson:"address"`
			Country  string `json:"country"  bson:"country"`
			City     string `json:"city"  bson:"city"`
			Province string `json:"province" bson:"province"`
			Zip_code int    `json:"zip_code" bson:"zip_code"`
		} `json:"shipping_address,omitempty" bson:"shipping_address,omitempty"`
	} `json:"address"`
	Email   string `json:"email" bson:"email"`
	Phone   string `json:"phone" bson:"phone"`
	Profile string `json:"profile" bson:"profile"`
}

func Mammals_getall(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	coll := docking.PakTradeDb.Collection("mammals")

	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver
	// get all user in  cursor
	cursor, err := coll.Find(context.Background(), coll)
	if err != nil {
		panic(err)
	}

	var results []Mammals_user
	for cursor.Next(context.TODO()) {
		var abc Mammals_user
		cursor.Decode(&abc)
		results = append(results, abc)

	}

	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)

	}
	fmt.Fprintf(w, "%s\n", output)

}
func Mammals_insertone(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	var strcutinit Mammals_user
	err := json.NewDecoder(req.Body).Decode(&strcutinit)
	if err != nil {
		panic(err)
	}

	insertdat := bson.M{
		"name": bson.M{
			"firt_name": strcutinit.Name.Firt_name,
			"last_name": strcutinit.Name.Last_name,
		},
		"address": bson.A{
			bson.M{
				"home_address": bson.M{
					"address":  strcutinit.Address[0].Home_address.Address,
					"country":  strcutinit.Address[0].Home_address.Country,
					"city":     strcutinit.Address[0].Home_address.City,
					"province": strcutinit.Address[0].Home_address.Province,
					"zip_code": strcutinit.Address[0].Home_address.Zip_code,
				},
				"shipping_address": bson.M{
					"address":  "",
					"country":  "",
					"city":     "",
					"province": "",
					"zip_code": 0,
				},
			},
		},
		"email":   strcutinit.Email,
		"phone":   strcutinit.Phone,
		"profile": strcutinit.Profile,
	}

	//fmt.Print(body)
	coll := docking.PakTradeDb.Collection("mammals")

	// // // insert a user

	inset, err3 := coll.InsertOne(context.TODO(), insertdat)
	if err3 != nil {
		fmt.Fprintf(w, "%s\n", err3)
	}

	fmt.Fprintf(w, "%s\n", inset)

}

func Mammals_select_one(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var search1 Mammals_serch
	err := json.NewDecoder(req.Body).Decode(&search1)
	if err != nil {
		panic(err)
	}
	coll := docking.PakTradeDb.Collection("mammals")
	objectIDS, _ := primitive.ObjectIDFromHex(search1.ID)

	var result Mammals_user
	filter := bson.M{"_id": objectIDS}

	err1 := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err1 != nil {
		fmt.Println("errror retrieving user userid : " + objectIDS.Hex())
	}

	// end findOne

	output, err2 := json.MarshalIndent(result, "", "    ")
	if err2 != nil {
		panic(err2)
	}

	fmt.Fprintf(w, "%s\n", output)

}
func Mammals_update_one(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var search1 Mammals_user_update
	err := json.NewDecoder(req.Body).Decode(&search1)
	if err != nil {
		panic(err)
	}
	coll := docking.PakTradeDb.Collection("mammals")
	objectIDS, _ := primitive.ObjectIDFromHex(search1.ID)
	// fmt.Print(objectIDS)

	result1, err := coll.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectIDS},
		bson.D{
			{Key: "$set", Value: bson.M{
				"name": bson.M{
					"firt_name": search1.Name.Firt_name,
					"last_name": search1.Name.Last_name,
				},
				"address": bson.A{
					bson.M{
						"home_address": bson.M{
							"address":  search1.Address[0].Home_address.Address,
							"country":  search1.Address[0].Home_address.Country,
							"city":     search1.Address[0].Home_address.City,
							"province": search1.Address[0].Home_address.Province,
							"zip_code": search1.Address[0].Home_address.Zip_code,
						},
						"shipping_address": bson.M{
							"address":  search1.Address[0].Shipping_address.Address,
							"country":  search1.Address[0].Shipping_address.Country,
							"city":     search1.Address[0].Shipping_address.City,
							"province": search1.Address[0].Shipping_address.Province,
							"zip_code": search1.Address[0].Shipping_address.Zip_code,
						},
					},
				},
				"email":   search1.Email,
				"phone":   search1.Phone,
				"profile": search1.Profile,
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

////////////// Mammalas Registration...././//////././

type mammals_reg_insert struct {
	// ID              primitive.ObjectID `bson:"_id,omitempty"`
	CreationDate    time.Time `json:"creationDate"`
	DisplayName     string    `json:"displayName"`
	Email           string    `json:"email"`
	IsAnonymour     bool      `json:"isAnonymour"`
	IsEmailVerified bool      `json:"isEmailVerified"`
	LastSignedIn    time.Time `json:"lastSignedIn"`
	PhoneNumber     string    `json:"phoneNumber"`
	PhotoURL        string    `json:"photoUrl"`
	ProviderID      string    `json:"providerId"`
	ProviderInfo    []struct {
		DisplyName  string `json:"displyName"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phoneNumber"`
		PhotoURL    string `json:"photoUrl"`
		ProviderID  string `json:"providerId"`
		UID         string `json:"uid"`
	} `json:"providerInfo"`
	PublicID     string `json:"publicId"`
	RefreshToken string `json:"refreshToken"`
}
type resp_insert struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Id      interface{} `json:"id"`
}

func Mammals_user_registration(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	var mammals_reg mammals_reg_insert
	err := json.NewDecoder(req.Body).Decode(&mammals_reg)
	if err != nil {
		panic(err)
	}
	// mongo query
	mongo_query := bson.M{
		"creationDate":    mammals_reg.CreationDate,
		"displayName":     mammals_reg.DisplayName,
		"email":           mammals_reg.Email,
		"isAnonymour":     mammals_reg.IsAnonymour,
		"isEmailVerified": mammals_reg.IsEmailVerified,
		"lastSignedIn":    mammals_reg.LastSignedIn,
		"phoneNumber":     mammals_reg.PhoneNumber,
		"photoUrl":        mammals_reg.PhotoURL,
		"providerId":      mammals_reg.ProviderID,
		"providerInfo": bson.A{
			bson.M{
				"displyName":  mammals_reg.ProviderInfo[0].DisplyName,
				"email":       mammals_reg.ProviderInfo[0].Email,
				"phoneNumber": mammals_reg.ProviderInfo[0].PhoneNumber,
				"photoUrl":    mammals_reg.ProviderInfo[0].PhotoURL,
				"providerId":  mammals_reg.ProviderInfo[0].ProviderID,
				"uid":         mammals_reg.ProviderInfo[0].UID,
			},
		},
		"publicId":     mammals_reg.PublicID,
		"refreshToken": mammals_reg.RefreshToken,
		"credit":       5,
		"planId":       "64735fe18f737b74c13bd6d3",
	}

	coll := docking.PakTradeDb.Collection("Mammalas_login")

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
