package signin

import (
	"context"
	"encoding/json"
	"net/http"
	docking "pak-trade-go/Docking"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	PublicId int                `json:"publicId" bson:"publicId"`
	Email    string             `json:"email" bson:"email"`
	ID       primitive.ObjectID `json:"id" bson:"_id"`
}

type EmailCheckResponse struct {
	PublicId int                `json:"publicId,omitempty"`
	Found    bool               `json:"found"`
	Message  string             `json:"message"`
	Status   int                `json:"status"`
	ID       primitive.ObjectID `json:"id" bson:"_id"`
}

type respone_struct1 struct {
	Status   int                `json:"status"`
	Message  string             `json:"message"`
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

	user, err := findUserByEmail(email)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, false, err.Error())
		return
	} else {
		response := EmailCheckResponse{
			PublicId: user.PublicId,
			Found:    true,
			Message:  "Email found",
			ID:       user.ID,
			Status:   200,
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

	_phoneInt, err := strconv.Atoi(phone)

	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, false, "Invalid phone number")
		return
	}

	if phone == "" {
		respondWithJSON(w, http.StatusBadRequest, false, "Phone parameter is missing")
		return
	}

	user, err := findUserByPhone(_phoneInt)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, false, err.Error())
		return
	} else {
		response := EmailCheckResponse{
			PublicId: user.PublicId,
			Found:    true,
			ID:       user.ID,
			Message:  "Phone found",
			Status:   200,
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

func findUserByPhone(phone int) (*User, error) {

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
		"exists":  exists,
		"message": message,
		"status":  statusCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
