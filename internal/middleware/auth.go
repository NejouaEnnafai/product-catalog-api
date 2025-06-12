package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/najwa/product-catalog-api/internal/auth"
	"github.com/najwa/product-catalog-api/internal/models"
)

// UserIDKey is the key used to store the user ID in the request context
type userIDKey string

const UserIDKey userIDKey = "userID"

// AuthMiddleware is a middleware that validates JWT tokens
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the request
		tokenString, err := auth.ExtractTokenFromRequest(r)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Validate the token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Add the user ID to the request context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID retrieves the user ID from the request context
func GetUserID(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	return userID, ok
}

// respondWithError responds with an error message
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, models.ErrorResponse{Error: message})
}

// respondWithJSON responds with JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
