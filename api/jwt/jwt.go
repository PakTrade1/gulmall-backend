package jwt

import (
	"encoding/json"
	"net/http"
	"os"
	"pak-trade-go/api/signin"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtoken = os.Getenv("JWT_TOKEN")

func GenerateJWT(user *signin.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID.Hex(),
		"exp": time.Now().Add(72 * time.Hour).Unix(), // expires in 3 days
	})
	tokenString, _ := token.SignedString([]byte(jwtoken))
	return tokenString
}

func ValidateTokenHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization") // should be "Bearer <token>"

	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("your_secret_key"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["id"].(string)
		user, err := signin.FindUserByID(userID)

		if err != nil || user == nil {
			http.Error(w, "Invalid user", http.StatusUnauthorized)
			return
		}

		json.NewEncoder(w).Encode(map[string]any{
			"valid": true,
			"user":  user,
		})
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}
