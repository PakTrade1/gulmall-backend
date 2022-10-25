package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbdata *mongo.Collection

func color(w http.ResponseWriter, req *http.Request) {

	// find code goes here
	coll := dbdata
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	// iterate code goes here
	for cursor.Next(context.TODO()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			panic(err)
		}
		//fmt.Println(result)

		output, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "%s\n", output)
	}
	if err := cursor.Err(); err != nil {
		panic(err)
	}

}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func dbconcet() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://developer001:IAmMuslim@cluster0.qeqntol.mongodb.net/test"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// database and colletion code goes hear
	db := client.Database("PakTrade")
	coll := db.Collection("color")
	dbdata = coll
	//fmt.Println("database conection seccessfully")
}
func main() {

	dbconcet()

	//fmt.Println("runging server port 9900")
	http.HandleFunc("/color", color)
	http.HandleFunc("/headers", headers)

	http.ListenAndServe(":9900", nil)
}
