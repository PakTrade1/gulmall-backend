package mammals

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	docking "pak-trade-go/Docking"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Mammals_serch struct {
	ID int `json:"id"`
}
type respone_struct struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    []Mammals_user1 `json:"data"`
}
type respone_struct1 struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    Mammals_user1 `json:"data"`
}
type respone_update struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Status string `json:"status_update"`
	} `json:"data"`
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
type Mammals_user1 struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ProviderID   interface{}        `json:"providerId"`
	PublicID     int                `json:"publicId"`
	RefreshToken interface{}        `json:"refreshToken"`
	DisplayName  interface{}        `json:"displayName"`
	LastSignedIn interface{}        `json:"lastSignedIn"`
	ProviderInfo []struct {
		PhotoURL    interface{} `json:"photoURL"`
		ProviderID  interface{} `json:"providerId"`
		UID         interface{} `json:"uid"`
		DisplyName  interface{} `json:"displyName"`
		Email       interface{} `json:"email"`
		PhoneNumber interface{} `json:"phoneNumber"`
	} `json:"providerInfo"`
	Credit          int         `json:"credit"`
	PhoneNumber     interface{} `json:"phoneNumber"`
	PhotoURL        interface{} `json:"photoUrl"`
	IsAnonymous     bool        `json:"isAnonymous"`
	IsEmailVerified bool        `json:"isEmailVerified"`
	CreationDate    interface{} `json:"creationDate"`
	ServerDate      time.Time   `json:"serverDate"`
	Email           interface{} `json:"email"`
	PlanID          string      `json:"planId"`
}

func Mammals_getall(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	coll := docking.PakTradeDb.Collection("Mammalas_login")

	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver
	// get all user in  cursor
	cursor, err := coll.Find(context.Background(), coll)
	if err != nil {
		panic(err)
	}

	var results []Mammals_user1
	var data respone_struct

	for cursor.Next(context.TODO()) {
		var abc Mammals_user1
		cursor.Decode(&abc)
		results = append(results, abc)

	}
	if results != nil {
		data.Status = http.StatusOK
		data.Message = "success"
		data.Data = results

	} else {
		data.Message = "decline"

	}

	output, err := json.MarshalIndent(data, "", "    ")
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

type datanotfound struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func Mammals_select_one(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var search1 Mammals_serch
	err := json.NewDecoder(req.Body).Decode(&search1)
	if err != nil {
		panic(err)
	}
	coll := docking.PakTradeDb.Collection("Mammalas_login")
	// objectIDS, _ := primitive.ObjectIDFromHex(search1.ID)

	var result Mammals_user1
	filter := bson.M{"publicId": search1.ID}

	err1 := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err1 != nil {
		fmt.Println("errror retrieving user userid : ")
	}

	// end findOne
	var data respone_struct1
	var datanot datanotfound
	if result.CreationDate != nil {
		data.Status = http.StatusOK
		data.Message = "success"
		data.Data = result
		output, err2 := json.MarshalIndent(data, "", "    ")
		if err2 != nil {
			panic(err2)
		}
		fmt.Fprintf(w, "%s\n", output)

	} else {
		datanot.Message = "decline"
		datanot.Status = "not found"
		output, err2 := json.MarshalIndent(datanot, "", "    ")
		if err2 != nil {
			panic(err2)
		}
		fmt.Fprintf(w, "%s\n", output)

	}

}

type mamal_login_update struct {
	ID           primitive.ObjectID `json:"userId"`
	DisplayName  string             `json:"displayName"`
	PhoneNumber  string             `json:"phoneNumber"`
	PhotoURL     string             `json:"photoUrl"`
	ProviderInfo []struct {
		PhoneNumber string `json:"phoneNumber"`
		PhotoURL    string `json:"photoURL"`
		ProviderID  string `json:"providerId"`
		UID         string `json:"uid"`
		DisplyName  string `json:"displyName"`
		Email       string `json:"email"`
	} `json:"providerInfo"`
	Email string `json:"email"`
}

func Mammals_update_one(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var payloadMap map[string]interface{}
	err1 := json.Unmarshal([]byte(body), &payloadMap)
	if err1 != nil {
		log.Fatal(err1)
	}
	userID := payloadMap["userId"].(string)
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"_id": objectID} // Assuming userId is the unique identifier for the document
	update := bson.M{"$set": payloadMap}

	coll := docking.PakTradeDb.Collection("Mammalas_login")
	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	var result_updte respone_update
	if result.ModifiedCount != 0 {
		result_updte.Status = http.StatusOK
		result_updte.Message = "success"
		result_updte.Data.Status = "Record updated"
	} else {
		result_updte.Message = "decline"
		result_updte.Data.Status = "No Change"
	}

	// fmt.Print(objectIDS)

	//end update

	output, err2 := json.MarshalIndent(result_updte, "", "    ")
	if err2 != nil {
		panic(err2)
	}

	fmt.Fprintf(w, "%s\n", output)

}

////////////// Mammalas Registration...././//////././

type mammals_reg_insert struct {
	// ID              primitive.ObjectID `bson:"_id,omitempty"`

	CreationDate    string `json:"creationDate"`
	DisplayName     string `json:"displayName"`
	Email           string `json:"email"`
	IsAnonymour     bool   `json:"isAnonymous"`
	IsEmailVerified bool   `json:"isEmailVerified"`
	LastSignedIn    string `json:"lastSignedIn"`
	PhoneNumber     string `json:"phoneNumber"`
	PhotoURL        string `json:"photoUrl"`
	ProviderID      string `json:"providerId"`
	ProviderInfo    []struct {
		DisplyName  string `json:"displyName"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phoneNumber"`
		PhotoURL    string `json:"photoUrl"`
		ProviderID  string `json:"providerId"`
		UID         string `json:"uid"`
	} `json:"providerInfo"`
	PublicID     int    `json:"publicId"`
	RefreshToken string `json:"refreshToken"`
}
type resp_insert struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Id      interface{} `json:"id"`
}

func getStringValue(value string) interface{} {
	if value == "" {
		return primitive.Null{}
	}
	return value
}

func Mammals_user_registration(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var mammals_reg mammals_reg_insert
	err := json.NewDecoder(req.Body).Decode(&mammals_reg)
	if err != nil {
		panic(err)
	}
	// mongo
	inputString := "64735fe18f737b74c13bd6d3"

	// Convert string to ObjectID
	Planid, err := primitive.ObjectIDFromHex(inputString)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	providerID := getStringValue(mammals_reg.ProviderID)
	providerID1 := getStringValue(mammals_reg.ProviderInfo[0].ProviderID)
	displayName := getStringValue(mammals_reg.DisplayName)
	email := getStringValue(mammals_reg.Email)
	phoneNumber := getStringValue(mammals_reg.PhoneNumber)
	photoUrl := getStringValue(mammals_reg.PhotoURL)
	displyName1 := getStringValue(mammals_reg.ProviderInfo[0].DisplyName)
	email1 := getStringValue(mammals_reg.ProviderInfo[0].Email)
	phoneNumber1 := getStringValue(mammals_reg.ProviderInfo[0].PhoneNumber)
	photoUrl1 := getStringValue(mammals_reg.ProviderInfo[0].PhotoURL)
	uid1 := getStringValue(mammals_reg.ProviderInfo[0].UID)
	refreshToken1 := getStringValue(mammals_reg.RefreshToken)
	mongo_query := bson.M{
		"creationDate":    mammals_reg.CreationDate,
		"serverDate":      time.Now(),
		"displayName":     displayName,
		"email":           email,
		"isAnonymous":     mammals_reg.IsAnonymour,
		"isEmailVerified": mammals_reg.IsEmailVerified,
		"lastSignedIn":    mammals_reg.LastSignedIn,
		"phoneNumber":     phoneNumber,
		"photoUrl":        photoUrl,
		"providerId":      providerID,
		"providerInfo": bson.A{
			bson.M{
				"displyName":  displyName1,
				"email":       email1,
				"phoneNumber": phoneNumber1,
				"photoURL":    photoUrl1,
				"providerId":  providerID1,
				"uid":         uid1,
			},
		},
		"publicId":     mammals_reg.PublicID,
		"refreshToken": refreshToken1,
		"credit":       5,
		"planId":       Planid,
		"adsRemaining": 5,
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

func CheckEmailHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In mail")
	if r.Method != http.MethodGet {
		respondWithJSON(w, http.StatusMethodNotAllowed, false, "Invalid request method")
		return
	}

	email := r.URL.Query().Get("email")
	if email == "" {
		respondWithJSON(w, http.StatusBadRequest, false, "Email parameter is missing")
		return
	}

	exists, err := CheckEmailExists(email)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, false, "Internal server error")
		return
	}
	if exists {
		respondWithJSON(w, http.StatusOK, true, "Email exists")
	} else {
		respondWithJSON(w, http.StatusOK, false, "Email does not exist")
	}
}

func CheckEmailExists(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := docking.PakTradeDb.Collection("Mammalas_login")

	filter := bson.M{"email": email}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func checkPhoneExists(phone int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := docking.PakTradeDb.Collection("Mammalas_login")

	filter := bson.M{"primaryPhone": phone}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func respondWithJSON(w http.ResponseWriter, statusCode int, exists bool, message string) {
	response := map[string]interface{}{
		"exists":  exists,
		"message": message,
		"status":  statusCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func checkIsEmailVerified(email string) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := docking.PakTradeDb.Collection("Mammalas_login")

	filter := bson.M{"email": email}
	var result struct {
		IsEmailVerified bool `bson:"isEmailVerified"`
	}
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil // Email does not exist
		}
		return false, err
	}
	return result.IsEmailVerified, nil
}

func CheckEmailVerifiedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithJSON(w, http.StatusMethodNotAllowed, false, "Invalid request method")
		return
	}

	email := r.URL.Query().Get("email")
	if email == "" {
		respondWithJSON(w, http.StatusBadRequest, false, "Email parameter is missing")
		return
	}

	verified, err := checkIsEmailVerified(email)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, false, "Internal server error")
		return
	}

	if verified {
		respondWithJSON(w, http.StatusOK, true, "Email is verified")
	} else {
		respondWithJSON(w, http.StatusOK, false, "Email is not verified")
	}
}

func CheckPhoneHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		respondWithJSON(w, http.StatusMethodNotAllowed, false, "Invalid request method")
		return
	}

	phone := r.URL.Query().Get("phone")

	if phone == "" {
		respondWithJSON(w, http.StatusBadRequest, false, "Phone parameter is missing")
		return
	}
	_phoneInt, err := strconv.Atoi(phone)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, false, "Invalid phone number")
		return
	}
	phoneInt, err := checkPhoneExists(_phoneInt)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, false, "Internal server error")
		return
	}
	if phoneInt {
		respondWithJSON(w, http.StatusOK, true, "Phone exists")
	} else {
		respondWithJSON(w, http.StatusOK, false, "Phone does not exist")
	}
}
