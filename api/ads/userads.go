package ads

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

type ads_init struct {
	Postid primitive.ObjectID `json:"postId"`
}
type Ads_resp struct {
	// ID         primitive.ObjectID `json:"id"`
	// EndDate    time.Time `json:"endDate"`
	// StartDate  time.Time `json:"startDate"`
	Post_status struct {
		Name  string `json:"name"`
		Order int    `json:"order"`
	} `json:"post_status"`
	Price  int `json:"price"`
	Images struct {
		HighQuility []string `json:"highQuility"`
		LowQuility  []struct {
			Image string `json:"image"`
			Color string `json:"color"`
		} `json:"lowQuility"`
	} `json:"images"`
	Status       string  `json:"status"`
	Rating       float64 `json:"rating"`
	ViewsCount   int     `json:"viewsCount"`
	CallCount    int     `json:"callCount"`
	ChatCount    int     `json:"chatCount"`
	FavCount     int     `json:"favCount"`
	CategoryName string  `json:"categoryName"`
	DaysRemening int     `json:"daysRemening"`
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

type respone_struct struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    []Ads_resp `json:"data"`
}

func Get_ads_user_by_post_id(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// urlParts := strings.Split(req.URL.Path, "/")
	// id := urlParts[len(urlParts)-1]

	id := req.URL.Query().Get("id")

	var strcutinit ads_init
	err := json.NewDecoder(req.Body).Decode(&strcutinit)
	if err != nil {
		panic(err)
	}
	handleError(err)

	coll := docking.PakTradeDb.Collection("post-status")
	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver
	ctx := context.TODO()
	cursor, err := coll.Aggregate(
		ctx, bson.A{
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "items-parent"},
						{"localField", "postId"},
						{"foreignField", "ownerId"},
						{"as", "test"},
					},
				},
			},
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "Mammalas_login"},
						{"localField", "result.ownerId"},
						{"foreignField", "ownerId"},
						{"as", "u"},
					},
				},
			},
			bson.D{{"$unwind", bson.D{{"path", "$u"}}}},
			bson.D{{"$match", bson.D{{"u.publicId", id}}}},
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "items-parent"},
						{"localField", "u._id"},
						{"foreignField", "ownerId"},
						{"as", "allpost"},
					},
				},
			},
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "plans"},
						{"localField", "u.planId"},
						{"foreignField", "_id"},
						{"as", "plans"},
					},
				},
			},
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "tier"},
						{"localField", "plans.tierId"},
						{"foreignField", "_id"},
						{"as", "tier"},
					},
				},
			},
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "post-statistics"},
						{"localField", "postId"},
						{"foreignField", "postId"},
						{"as", "statistics"},
					},
				},
			},
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "categories"},
						{"localField", "allpost.category"},
						{"foreignField", "_id"},
						{"as", "catName"},
					},
				},
			},
			bson.D{{"$unwind", bson.D{{"path", "$allpost"}}}},
			bson.D{
				{"$set",
					bson.D{
						{"cat", bson.D{{"$first", "$catName.name.en"}}},
						{"statistics", bson.D{{"$first", "$statistics"}}},
						{"post_status", bson.D{{"$first", "$tier"}}},
					},
				},
			},
			bson.D{
				{"$project",
					bson.D{
						{"price", "$allpost.price"},
						{"images", "$allpost.images"},
						{"status", "$allpost.status"},
						{"rating", "$allpost.rating"},
						{"viewsCount", "$statistics.views"},
						{"callCount", "$statistics.calls"},
						{"chatCount", "$statistics.chats"},
						{"favCount", "$statistics.favorites"},
						{"categoryName", "$cat"},
						{"post_status",
							bson.D{
								{"name", "$post_status.name"},
								{"order", "$post_status.order"},
							},
						},
						{"daysRemening",
							bson.D{
								{"$floor",
									bson.D{
										{"$divide",
											bson.A{
												bson.D{
													{"$subtract",
														bson.A{
															bson.D{{"$toDate", "$endDate"}},
															bson.D{{"$toDate", "$allpost.creationTimestamp"}},
														},
													},
												},
												86400000,
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
	)
	if err != nil {
		log.Fatal(err)
	}

	var resp1 []Ads_resp
	var results respone_struct
	//  fmt.Println(cursor)

	for cursor.Next(context.TODO()) {

		var xy Ads_resp
		cursor.Decode(&xy)
		resp1 = append(resp1, xy)

	}
	if resp1 != nil {
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
