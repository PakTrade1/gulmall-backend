package signin

import (
	"encoding/json"
	"net/http"
	"pak-trade-go/api/mammals"
)

type respone_struct1 struct {
	Status   int    `json:"status"`
	Message  string `json:"message"`
	PublicID int    `json:"publicId"`
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

	exist, err := mammals.CheckEmailExists(email)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, false, "Internal server error")
		return
	}

	if exist {
		// look
		respondWithJSON(w, http.StatusOK, true, "Email is Exist")
	} else {
		respondWithJSON(w, http.StatusOK, false, "Email is not Exist")
	}
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
