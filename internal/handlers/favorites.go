package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/najwa/product-catalog-api/internal/db"
	"github.com/najwa/product-catalog-api/internal/middleware"
	"github.com/najwa/product-catalog-api/internal/models"
)

// AddFavoriteHandler handles adding a product to the user's favorites
func AddFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get the user ID from the context
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse the request body
	var req models.FavoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate the request
	if req.ProductID <= 0 {
		respondWithError(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	// Add the favorite to the database
	err := db.AddFavorite(userID, req.ProductID, req.Notes)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error adding favorite: "+err.Error())
		return
	}

	// Return success
	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "Favorite added successfully"})
}

// GetFavoritesHandler handles retrieving the user's favorite products
func GetFavoritesHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get the user ID from the context
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get the favorites from the database
	favorites, err := db.GetFavorites(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving favorites")
		return
	}

	// Return the favorites
	respondWithJSON(w, http.StatusOK, favorites)
}
