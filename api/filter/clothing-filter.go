package clothingFilter

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type filter_list struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name" json:"name"`
}

type Filter struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	FilterName string             `bson:"filter_name" json:"filter_name"`
	FilterList []filter_list      `bson:"filter_data" json: "filter_data"`
}

func getFilters(client *mongo.Client) ([]Filter, error) {

	ctx := context.TODO()
	coll := client.Database("PakTrade").Collection("clothing_filter")
	cursor, err := coll.Aggregate(ctx, bson.A{
		bson.D{{"$unwind", bson.D{{"path", "$list"}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "brands"},
					{"localField", "list"},
					{"foreignField", "_id"},
					{"as", "brand_data"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "color"},
					{"localField", "list"},
					{"foreignField", "_id"},
					{"as", "color_data"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "size"},
					{"localField", "list"},
					{"foreignField", "_id"},
					{"as", "size_data"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "fabric"},
					{"localField", "list"},
					{"foreignField", "_id"},
					{"as", "fabric_data"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "gender"},
					{"localField", "list"},
					{"foreignField", "_id"},
					{"as", "gender_data"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "season"},
					{"localField", "list"},
					{"foreignField", "_id"},
					{"as", "season_data"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "cloth_type"},
					{"localField", "list"},
					{"foreignField", "_id"},
					{"as", "cloth_data"},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"filter_data",
						bson.D{
							{"$concatArrays",
								bson.A{
									bson.D{
										{"$ifNull",
											bson.A{
												"$brand_data",
												bson.A{},
											},
										},
									},
									bson.D{
										{"$ifNull",
											bson.A{
												"$color_data",
												bson.A{},
											},
										},
									},
									bson.D{
										{"$ifNull",
											bson.A{
												"$size_data",
												bson.A{},
											},
										},
									},
									bson.D{
										{"$ifNull",
											bson.A{
												"$fabric_data",
												bson.A{},
											},
										},
									},
									bson.D{
										{"$ifNull",
											bson.A{
												"$gender_data",
												bson.A{},
											},
										},
									},
									bson.D{
										{"$ifNull",
											bson.A{
												"$season_data",
												bson.A{},
											},
										},
									},
									bson.D{
										{"$ifNull",
											bson.A{
												"$cloth_data",
												bson.A{},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$_id"},
					{"filter_name", bson.D{{"$first", "$filter_name"}}},
					{"filter_data", bson.D{{"$push", "$filter_data"}}},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"filter_name", 1},
					{"list", 1},
					{"filter_data",
						bson.D{
							{"$reduce",
								bson.D{
									{"input", "$filter_data"},
									{"initialValue", bson.A{}},
									{"in",
										bson.D{
											{"$concatArrays",
												bson.A{
													"$$value",
													"$$this",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.TODO())
	var filters []Filter

	for cursor.Next(ctx) {
		var filter Filter
		if err := cursor.Decode(&filter); err != nil {
			return nil, err
		}
		filters = append(filters, filter)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return filters, nil

}

func FiltersHandler(w http.ResponseWriter, r *http.Request) {
	clientOptions := options.Client().ApplyURI("mongodb+srv://developer001:IAmMuslim@cluster0.qeqntol.mongodb.net/test")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.TODO())

	filters, err := getFilters(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filters)

}
