package items

import (
	"context"
	"encoding/json"
	"net/http"
	docking "pak-trade-go/Docking"
	user "pak-trade-go/models"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Response struct {
	User  user.User   `json:"user"`
	Items []user.Item `json:"items"`
}

func GetUserAndItemsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	publicId := r.URL.Query().Get("publicId")
	if publicId == "" {
		http.Error(w, "publicId parameter is missing", http.StatusBadRequest)
		return
	}

	publicIdInt, err := strconv.Atoi(publicId)
	if err != nil {
		http.Error(w, "Invalid publicId parameter", http.StatusBadRequest)
		return
	}

	user, err := findUserByPublicId(publicIdInt)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	items, err := findItemsByOwnerId(user.ID)
	if err != nil {
		http.Error(w, "Error retrieving items", http.StatusInternalServerError)
		return
	}

	response := Response{
		User:  *user,
		Items: items,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func findUserByPublicId(publicId int) (*user.User, error) {
	userCollection := docking.PakTradeDb.Collection("Mammalas_login")

	var user user.User
	filter := bson.M{"publicId": publicId}
	err := userCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func findItemsByOwnerId(ownerId primitive.ObjectID) ([]user.Item, error) {
	itemCollection := docking.PakTradeDb.Collection("items-parent")

	var items []user.Item
	filter := bson.M{"ownerId": ownerId}
	cursor, err := itemCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var item user.Item
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
