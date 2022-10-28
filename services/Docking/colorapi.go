package docking

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type Color1 struct {
// 	ID     primitive.ObjectID `json:"_id,omitempty"`
// 	CSSHex string             `json:"cssHex,omitempty"`
// 	Name   struct {
// 		Ar string `json:"ar,omitempty"`
// 		En string `json:"en,omitempty"`
// 	} `json:"name,omitempty"`
// }

type Name struct {
	Ar string `json:"ar,omitempty"`
	En string `json:"en,omitempty"`
}

type Color1 struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	CSSHex string             `json:"cssHex,omitempty"`
	Name   `json:"name,omitempty"`
}

// type color1 struct {
// 	ID     primitive.ObjectID `bson:"_id,omitempty"`
// 	CSSHEX string             ` bson:"csshex" ,omitempty`
// 	NAME   Name               `bson:"inline"`
// }
// type Name struct {
// 	EN string `bson:"en" ,omitempty`
// 	AR string `bson:"ar" ,omitempty`
// }

func Color(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	coll := Database.Collection("color")
	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}

	var results []Color1
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	for _, result := range results {
		cursor.Decode(&results)

		results = append(results, result)

	}
	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "%s\n", output)

}
