package signin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	docking "pak-trade-go/Docking"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	PublicId int64              `json:"publicId" bson:"publicId"`
	Email    string             `json:"email" bson:"email"`
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	IP       string             `json:"ip" bson:"ip"`
}

type EmailCheckResponse struct {
	Found  bool `json:"found"`
	Status int  `json:"status"`
}

type respone_struct1 struct {
	Status   int                `json:"status"`
	PublicID int                `json:"publicId"`
	ID       primitive.ObjectID `json:"id" bson:"_id"`
}

func SignInEmailHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		respondWithJSON(w, http.StatusMethodNotAllowed, false, "Invalid request method")
		return
	}

	email := r.URL.Query().Get("email")
	if email == "" {
		respondWithJSON(w, http.StatusBadRequest, false, "Email parameter is missing")
		return
	}

	_, err := findUserByEmail(email)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, false, err.Error())
		return
	} else {
		response := EmailCheckResponse{

			Found:  true,
			Status: 200,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	}

}

func SignInPhoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithJSON(w, http.StatusMethodNotAllowed, false, "Invalid request method")
		return
	}

	phone := r.URL.Query().Get("phone")

	if phone == "" {
		respondWithJSON(w, http.StatusBadRequest, false, "Phone parameter is missing")
		return
	}

	_, err := FindUserByPhone(phone)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, false, err.Error())
		return
	} else {
		response := EmailCheckResponse{
			Found:  true,
			Status: 200,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	}
}

func findUserByEmail(email string) (*User, error) {

	collection := docking.PakTradeDb.Collection("Mammalas_login")
	var user User
	filter := bson.M{"email": email}
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FindUserByID(idStr string) (*User, error) {
	collection := docking.PakTradeDb.Collection("Mammalas_login")
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	var user User
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &user, nil
}

func FindUserByPhone(phone string) (*User, error) {

	collection := docking.PakTradeDb.Collection("Mammalas_login")
	var user User
	filter := bson.M{"primaryPhone": phone}
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func respondWithJSON(w http.ResponseWriter, statusCode int, exists bool, message string) {
	response := map[string]interface{}{
		"exists": exists,
		"status": statusCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
