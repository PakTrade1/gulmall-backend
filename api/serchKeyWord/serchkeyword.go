package serchkeyword

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	docking "pak-trade-go/Docking"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type serchkey struct {
	User_public_id string `json:"user_id"`
	Keyword        string `json:"keyword"`
}
type resp_insert struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Id      interface{} `json:"id"`
}

func Serchkeywordinsert(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	var keyserch serchkey
	err := json.NewDecoder(req.Body).Decode(&keyserch)
	if err != nil {
		panic(err)
	}

	// mongo

	// Convert string to ObjectID
	inputString := keyserch.User_public_id

	// Convert string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(inputString)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	mongo_query := bson.M{
		"keyword": keyserch.Keyword,
		"userId":  objectID,
		"time":    time.Now(),
	}

	coll := docking.PakTradeDb.Collection("searched_keyword")

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
