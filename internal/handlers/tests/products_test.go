package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/najwa/product-catalog-api/internal/db"
	"github.com/najwa/product-catalog-api/internal/handlers"
	"github.com/najwa/product-catalog-api/internal/models"
)

func TestProductsHandler(t *testing.T) {
	// Set up test database
	dbPath := filepath.Join(os.TempDir(), "test_products.db")
	
	// Clean up the database file after the test
	defer os.Remove(dbPath)
	
	// Initialize the database
	err := db.Initialize(dbPath)
	if err != nil {
		t.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()
	
	// Seed the database with test data
	seedTestProducts()
	
	// Test cases
	testCases := []struct {
		name           string
		url            string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "Get all products",
			url:            "/products",
			expectedStatus: http.StatusOK,
			expectedCount:  3, // We'll seed 3 test products
		},
		{
			name:           "Filter by category",
			url:            "/products?category=electronics",
			expectedStatus: http.StatusOK,
			expectedCount:  2, // We'll seed 2 electronics products
		},
		{
			name:           "Sort by price ascending",
			url:            "/products?sort=price_asc",
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
		{
			name:           "Search by title",
			url:            "/products?search=phone",
			expectedStatus: http.StatusOK,
			expectedCount:  1, // We'll seed 1 product with "phone" in the title
		},
		{
			name:           "Pagination",
			url:            "/products?page=1&limit=2",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request
			req, err := http.NewRequest("GET", tc.url, nil)
			if err != nil {
				t.Fatalf("Error creating request: %v", err)
			}
			
			// Create a response recorder
			rr := httptest.NewRecorder()
			
			// Create a handler
			handler := http.HandlerFunc(handlers.ProductsHandler)
			
			// Serve the request
			handler.ServeHTTP(rr, req)
			
			// Check the status code
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}
			
			// Check the response body
			var response struct {
				Total   int             `json:"total"`
				Page    int             `json:"page"`
				Limit   int             `json:"limit"`
				Results []models.Product `json:"results"`
			}
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Error unmarshaling response: %v", err)
			}
			
			// Check the number of products
			if len(response.Results) != tc.expectedCount {
				t.Errorf("Expected %d products, got %d", tc.expectedCount, len(response.Results))
			}
		})
	}
}

// seedTestProducts seeds the database with test products
func seedTestProducts() {
	// Clear the products table
	db.DB.Exec("DELETE FROM products")
	
	// Insert test products
	products := []struct {
		title    string
		price    float64
		category string
		image    string
	}{
		{
			title:    "Smartphone",
			price:    499.99,
			category: "electronics",
			image:    "https://example.com/smartphone.jpg",
		},
		{
			title:    "Laptop",
			price:    999.99,
			category: "electronics",
			image:    "https://example.com/laptop.jpg",
		},
		{
			title:    "T-Shirt",
			price:    19.99,
			category: "clothing",
			image:    "https://example.com/tshirt.jpg",
		},
	}
	
	for _, p := range products {
		db.DB.Exec(
			"INSERT INTO products (title, price, category, image) VALUES (?, ?, ?, ?)",
			p.title, p.price, p.category, p.image,
		)
	}
}
