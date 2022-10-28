package items

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	docking "pak-trade-go/Docking"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type itemsstruct struct {
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
	// Available_color []struct {
	// 	Available_color []string `json:"available_color,omitempty"`
	// } `json:"available_color,omitempty"`
	Available_color []string `json:"available_color,omitempty"`
	Available_size  []string `json:"available_size,omitempty"`

	// Available_size []struct {
	// 	Available_size []string `json:"available_size,omitempty"`
	// } `json:"available_size,omitempty"`

	Images []struct {
		Low_quility  []string `json:"low_quility,omitempty"`
		High_quility []string `json:"high_quility,omitempty"`
	} `json:"images,omitempty"`
}

// type Images []struct {
// 	LowQuility  []string `json:"low_quility,omitempty"`
// 	HighQuility []string `json:"high_quility,omitempty"`
// }

// type AvailableColor struct {
// 	Availablecolor []string `json:"availablecolor,omitempty"`
// }

// // type Namefeature struct {
// // 	En             string `json:"en,omitempty"`
// // 	Ar             string `json:"ar,omitempty"`
// // 	AvailableColor `json:"availableColo,omitempty"`
// // }

// // type Feature struct {
// // 	Name           `json:"name,omitempty"`
// // 	AvailableColor `json:"availableColo,omitempty"`
// // }

// type Name struct {
// 	En      string `json:"en,omitempty"`
// 	Ar      string `json:"ar,omitempty"`
// 	Feature struct {
// 		Name struct {
// 			En string `json:"en,omitempty"`
// 			Ar string `json:"ar,omitempty"`
// 		} `json:"name,omitempty"`
// 	} `json:"feature,omitempty"`
// }

// type itemsstruct struct {
// 	Id             primitive.ObjectID `bson:"_id,omitempty"`
// 	Name           `json:"name,omitempty"`
// 	AvailableColor `json:"availablecolor,omitempty"`
// 	Images         `json:"images,omitempty"`
// }

// type itemsstruct1 struct {
// 	ID struct {
// 		Id string `json:"$oid,omitempty"`
// 	} `json:"_id,omitempty"`
// 	Name struct {
// 		En string `json:"en,omitempty"`
// 		Ar string `json:"ar,omitempty"`
// 	} `json:"name,omitempty"`
// 	Feature []struct {
// 		Name struct {
// 			En string `json:"en,omitempty"`
// 			Ar string `json:"ar,omitempty"`
// 		} `json:"name,omitempty"`
// 	} `json:"feature,omitempty"`
// 	AvailableColor []string `json:"available_color,omitempty"`
// 	Images         []struct {
// 		LowQuility  []string `json:"low_quility,omitempty"`
// 		HighQuility []string `json:"high_quility,omitempty"`
// 	} `json:"images,omitempty"`
// }

func Items(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	coll := docking.ItemDb.Collection("cloths")
	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}
	var results []itemsstruct
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for cursor.Next(context.TODO()) {
		var abc itemsstruct
		cursor.Decode(&abc)
		results = append(results, abc)

	}

	// for _, result := range results {
	// 	cursor.Decode(&results)
	// 	results = append(results, result)
	// }
	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)

	}
	fmt.Fprintf(w, "%s\n", output)

}
