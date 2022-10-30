package items

import (
	"context"
	"fmt"
	"net/http"
	docking "pak-trade-go/Docking"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Name struct {
	Ar string `json:"ar,omitempty"`
	En string `json:"en,omitempty"`
}

type ColorType struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	CSSHex string             `json:"cssHex,omitempty"`
	Name   `json:"name,omitempty"`
}

type ItemType struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	Name struct {
		En string `json:"en,omitempty"`
		Ar string `json:"ar,omitempty"`
	} `json:"name,omitempty"`

	Feature []struct {
		Name struct {
			En string `json:"en,omitempty"`
			Ar string `json:"ar,omitempty"`
		} `json:"name,omitempty"`
	} `json:"feature,omitempty"`

	Available_color []string `json:"available_color,omitempty"`
	Available_size  []string `json:"available_size,omitempty"`

	Images []struct {
		Low_quility  []string `json:"low_quility,omitempty"`
		High_quility []string `json:"high_quility,omitempty"`
	} `json:"images,omitempty"`
}

func Items(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	coll := docking.ItemDb.Collection("cloths")
	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}
	var results []ItemType
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for cursor.Next(context.TODO()) {
		var abc ItemType
		cursor.Decode(&abc)
		results = append(results, abc)

	}
	//id, _ := primitive.ObjectIDFromHex("6352f8123e006819c56246c6")

	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "color"}, {Key: "localField", Value: "color"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "color"}}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$color"}, {Key: "total", Value: bson.D{{Key: "$preserveNullAndEmptyArrays", Value: false}}}}}}

	showLoadedStructCursor, err := docking.PakTradeDb.Aggregate(context.TODO(), mongo.Pipeline{lookupStage, unwindStage})
	if err != nil {
		panic(err)
	}
	var showsLoadedStruct []ColorType
	if err = showLoadedStructCursor.All(context.TODO(), &showsLoadedStruct); err != nil {
		panic(err)
	}
	fmt.Println(showsLoadedStruct)

	// output, err := json.MarshalIndent(results, "", "    ")
	// if err != nil {
	// 	panic(err)

	// }
	// fmt.Fprintf(w, "%s\n", output)

}
