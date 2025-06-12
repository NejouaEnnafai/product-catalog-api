package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/najwa/product-catalog-api/internal/db"
)

// Sample products data
var products = []struct {
	Title    string
	Price    float64
	Category string
	Image    string
}{
	{
		Title:    "Smartphone X",
		Price:    999.99,
		Category: "electronics",
		Image:    "https://example.com/smartphone.jpg",
	},
	{
		Title:    "Laptop Pro",
		Price:    1499.99,
		Category: "electronics",
		Image:    "https://example.com/laptop.jpg",
	},
	{
		Title:    "Wireless Headphones",
		Price:    199.99,
		Category: "electronics",
		Image:    "https://example.com/headphones.jpg",
	},
	{
		Title:    "Smart Watch",
		Price:    299.99,
		Category: "electronics",
		Image:    "https://example.com/smartwatch.jpg",
	},
	{
		Title:    "Cotton T-Shirt",
		Price:    19.99,
		Category: "clothing",
		Image:    "https://example.com/tshirt.jpg",
	},
	{
		Title:    "Jeans",
		Price:    49.99,
		Category: "clothing",
		Image:    "https://example.com/jeans.jpg",
	},
	{
		Title:    "Running Shoes",
		Price:    89.99,
		Category: "footwear",
		Image:    "https://example.com/shoes.jpg",
	},
	{
		Title:    "Backpack",
		Price:    39.99,
		Category: "accessories",
		Image:    "https://example.com/backpack.jpg",
	},
	{
		Title:    "Water Bottle",
		Price:    14.99,
		Category: "accessories",
		Image:    "https://example.com/bottle.jpg",
	},
	{
		Title:    "Fitness Tracker",
		Price:    79.99,
		Category: "electronics",
		Image:    "https://example.com/tracker.jpg",
	},
}

// Sample user data
var users = []struct {
	Username string
	Password string
}{
	{
		Username: "john",
		Password: "1234",
	},
}

func main() {
	// Initialize the database
	dbPath := filepath.Join(".", "product_catalog.db")
	err := db.Initialize(dbPath)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	// Seed users
	for _, user := range users {
		_, err := db.CreateUser(user.Username, user.Password)
		if err != nil {
			log.Printf("Error seeding user %s: %v", user.Username, err)
		} else {
			log.Printf("Seeded user: %s", user.Username)
		}
	}

	// Seed products
	for _, product := range products {
		_, err := db.DB.Exec(
			"INSERT INTO products (title, price, category, image) VALUES (?, ?, ?, ?)",
			product.Title, product.Price, product.Category, product.Image,
		)
		if err != nil {
			log.Printf("Error seeding product %s: %v", product.Title, err)
		} else {
			log.Printf("Seeded product: %s", product.Title)
		}
	}

	log.Println("Database seeded successfully!")

	// Run the main application
	cmd := exec.Command("go", "run", filepath.Join("cmd", "main.go"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println("Starting the application...")
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error running the application: %v", err)
	}
}
