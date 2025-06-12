package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"github.com/najwa/product-catalog-api/internal/db"
	"github.com/najwa/product-catalog-api/internal/handlers"
	"github.com/najwa/product-catalog-api/internal/middleware"
)

func main() {
	// Parse command-line flags
	port := flag.String("port", "8080", "Port to listen on")
	dbPath := flag.String("db", "./product_catalog.db", "Path to SQLite database file")
	flag.Parse()

	// Initialize the database
	absDBPath, err := filepath.Abs(*dbPath)
	if err != nil {
		log.Fatalf("Error resolving database path: %v", err)
	}

	err = db.Initialize(absDBPath)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	// Set up routes
	setupRoutes()

	// Start the server
	log.Printf("Server starting on port %s...", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func setupRoutes() {
	// Public routes
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/products", handlers.ProductsHandler)

	// Protected routes
	// Create a subrouter for protected routes
	favoritesHandler := http.HandlerFunc(handlers.AddFavoriteHandler)
	getFavoritesHandler := http.HandlerFunc(handlers.GetFavoritesHandler)

	// Apply auth middleware to protected routes
	http.Handle("/favorites", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			favoritesHandler.ServeHTTP(w, r)
		case http.MethodGet:
			getFavoritesHandler.ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))
}
