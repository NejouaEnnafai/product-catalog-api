package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/najwa/product-catalog-api/internal/auth"
	"github.com/najwa/product-catalog-api/internal/db"
	"github.com/najwa/product-catalog-api/internal/handlers"
	"github.com/najwa/product-catalog-api/internal/middleware"
	"github.com/najwa/product-catalog-api/internal/models"
)

func TestFavoritesHandler(t *testing.T) {
	// Set up test database
	dbPath := filepath.Join(os.TempDir(), "test_favorites.db")
	
	// Clean up the database file after the test
	defer os.Remove(dbPath)
	
	// Initialize the database
	err := db.Initialize(dbPath)
	if err != nil {
		t.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()
	
	// Seed the database with test data
	userID := seedTestUser()
	seedTestProductsForFavorites()
	
	// Generate a JWT token for the test user
	token, err := auth.GenerateToken(userID)
	if err != nil {
		t.Fatalf("Error generating token: %v", err)
	}
	
	// Test adding a favorite
	t.Run("Add favorite", func(t *testing.T) {
		// Create a request body
		reqBody := models.FavoriteRequest{
			ProductID: 1,
			Notes:     "Test note",
		}
		
		body, err := json.Marshal(reqBody)
		if err != nil {
			t.Fatalf("Error marshaling request body: %v", err)
		}
		
		// Create a request
		req, err := http.NewRequest("POST", "/favorites", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		
		// Add the Authorization header
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		
		// Create a response recorder
		rr := httptest.NewRecorder()
		
		// Create a handler with auth middleware
		handler := middleware.AuthMiddleware(http.HandlerFunc(handlers.AddFavoriteHandler))
		
		// Serve the request
		handler.ServeHTTP(rr, req)
		
		// Check the status code
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}
	})
	
	// Test getting favorites
	t.Run("Get favorites", func(t *testing.T) {
		// Create a request
		req, err := http.NewRequest("GET", "/favorites", nil)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		
		// Add the Authorization header
		req.Header.Set("Authorization", "Bearer "+token)
		
		// Create a response recorder
		rr := httptest.NewRecorder()
		
		// Create a handler with auth middleware
		handler := middleware.AuthMiddleware(http.HandlerFunc(handlers.GetFavoritesHandler))
		
		// Serve the request
		handler.ServeHTTP(rr, req)
		
		// Check the status code
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		
		// Check the response body
		var favorites []models.Product
		err = json.Unmarshal(rr.Body.Bytes(), &favorites)
		if err != nil {
			t.Fatalf("Error unmarshaling response: %v", err)
		}
		
		// We should have 1 favorite (added in the previous test)
		if len(favorites) != 1 {
			t.Errorf("Expected 1 favorite, got %d", len(favorites))
		}
	})
	
	// Test unauthorized access
	t.Run("Unauthorized access", func(t *testing.T) {
		// Create a request
		req, err := http.NewRequest("GET", "/favorites", nil)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		
		// No Authorization header
		
		// Create a response recorder
		rr := httptest.NewRecorder()
		
		// Create a handler with auth middleware
		handler := middleware.AuthMiddleware(http.HandlerFunc(handlers.GetFavoritesHandler))
		
		// Serve the request
		handler.ServeHTTP(rr, req)
		
		// Check the status code
		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})
}

// seedTestUser seeds the database with a test user and returns the user ID
func seedTestUser() int {
	// Clear the users table
	db.DB.Exec("DELETE FROM users")
	
	// Create a test user
	user, _ := db.CreateUser("testuser", "password")
	return user.ID
}

// seedTestProductsForFavorites seeds the database with test products for favorites
func seedTestProductsForFavorites() {
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
			title:    "Test Product 1",
			price:    99.99,
			category: "test",
			image:    "https://example.com/test1.jpg",
		},
		{
			title:    "Test Product 2",
			price:    199.99,
			category: "test",
			image:    "https://example.com/test2.jpg",
		},
	}
	
	for _, p := range products {
		db.DB.Exec(
			"INSERT INTO products (title, price, category, image) VALUES (?, ?, ?, ?)",
			p.title, p.price, p.category, p.image,
		)
	}
}
