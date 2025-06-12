package handlers

import (
	"net/http"
	"strconv"

	"github.com/najwa/product-catalog-api/internal/db"
	"github.com/najwa/product-catalog-api/internal/models"
)

// ProductsHandler handles product listing with filtering, sorting, and pagination
func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse query parameters
	query := r.URL.Query()
	
	// Parse pagination parameters
	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))
	
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	
	// Parse filtering parameters
	category := query.Get("category")
	sort := query.Get("sort")
	search := query.Get("search")
	
	// Get products from the database
	products, total, err := db.GetProducts(page, limit, category, sort, search)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving products")
		return
	}
	
	// Return the products
	respondWithJSON(w, http.StatusOK, models.PaginatedResponse{
		Total:   total,
		Page:    page,
		Limit:   limit,
		Results: products,
	})
}
