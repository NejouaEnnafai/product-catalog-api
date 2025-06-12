package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

// DB is the database connection
var DB *sql.DB

// Initialize initializes the database connection and creates tables if they don't exist
func Initialize(dbPath string) error {
	var err error

	// Check if the database file exists
	_, err = os.Stat(dbPath)
	if os.IsNotExist(err) {
		// Create the database file
		file, err := os.Create(dbPath)
		if err != nil {
			return fmt.Errorf("error creating database file: %w", err)
		}
		file.Close()
	}

	// Open the database
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	// Create tables if they don't exist
	if err = createTables(); err != nil {
		return fmt.Errorf("error creating tables: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// createTables creates the necessary tables if they don't exist
func createTables() error {
	// Create users table
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}

	// Create products table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			price REAL NOT NULL,
			category TEXT NOT NULL,
			image TEXT NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating products table: %w", err)
	}

	// Create favorites table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS favorites (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			notes TEXT,
			FOREIGN KEY (user_id) REFERENCES users (id),
			FOREIGN KEY (product_id) REFERENCES products (id),
			UNIQUE (user_id, product_id)
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating favorites table: %w", err)
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
