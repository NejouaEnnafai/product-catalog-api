package db

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/najwa/product-catalog-api/internal/models"
)

// GetUserByUsername retrieves a user by username
func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := DB.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).Scan(
		&user.ID, &user.Username, &user.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error querying user: %w", err)
	}
	return &user, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := DB.QueryRow("SELECT id, username, password FROM users WHERE id = ?", id).Scan(
		&user.ID, &user.Username, &user.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error querying user: %w", err)
	}
	return &user, nil
}

// CreateUser creates a new user
func CreateUser(username, password string) (*models.User, error) {
	// Hash the password
	hashedPassword := hashPassword(password)

	result, err := DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %w", err)
	}

	return &models.User{
		ID:       int(id),
		Username: username,
		Password: hashedPassword,
	}, nil
}

// hashPassword hashes a password using SHA-256
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// ValidatePassword validates a password against a hash
func ValidatePassword(password, hash string) bool {
	return hashPassword(password) == hash
}
