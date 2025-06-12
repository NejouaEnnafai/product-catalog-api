package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/najwa/product-catalog-api/internal/auth"
	"github.com/najwa/product-catalog-api/internal/db"
	"github.com/najwa/product-catalog-api/internal/models"
)

// LoginHandler handles user login and returns a JWT token
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse the request body
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate the request
	if req.Username == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	// Get the user from the database
	user, err := db.GetUserByUsername(req.Username)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Validate the password
	if !db.ValidatePassword(req.Password, user.Password) {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate a JWT token
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	// Return the token
	respondWithJSON(w, http.StatusOK, models.LoginResponse{Token: token})
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
