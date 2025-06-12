package db

import (
	"database/sql"
	"fmt"

	"github.com/najwa/product-catalog-api/internal/models"
)

// AddFavorite adds a product to a user's favorites
func AddFavorite(userID, productID int, notes string) error {
	// Check if the product exists
	_, err := GetProductByID(productID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	// Check if the favorite already exists
	var id int
	err = DB.QueryRow("SELECT id FROM favorites WHERE user_id = ? AND product_id = ?", userID, productID).Scan(&id)
	if err == nil {
		// Favorite already exists, update the notes
		_, err = DB.Exec("UPDATE favorites SET notes = ? WHERE id = ?", notes, id)
		if err != nil {
			return fmt.Errorf("error updating favorite: %w", err)
		}
		return nil
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("error checking favorite: %w", err)
	}

	// Add the favorite
	_, err = DB.Exec("INSERT INTO favorites (user_id, product_id, notes) VALUES (?, ?, ?)", userID, productID, notes)
	if err != nil {
		return fmt.Errorf("error adding favorite: %w", err)
	}

	return nil
}

// GetFavorites retrieves a user's favorite products
func GetFavorites(userID int) ([]models.Product, error) {
	// Query for favorite products
	rows, err := DB.Query(`
		SELECT p.id, p.title, p.price, p.category, p.image, f.notes
		FROM favorites f
		JOIN products p ON f.product_id = p.id
		WHERE f.user_id = ?
		ORDER BY f.id DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying favorites: %w", err)
	}
	defer rows.Close()

	// Parse the results
	favorites := []models.Product{}
	for rows.Next() {
		var product models.Product
		var notes sql.NullString
		err := rows.Scan(
			&product.ID,
			&product.Title,
			&product.Price,
			&product.Category,
			&product.Image,
			&notes,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning favorite: %w", err)
		}
		favorites = append(favorites, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating favorites: %w", err)
	}

	return favorites, nil
}

// RemoveFavorite removes a product from a user's favorites
func RemoveFavorite(userID, productID int) error {
	result, err := DB.Exec("DELETE FROM favorites WHERE user_id = ? AND product_id = ?", userID, productID)
	if err != nil {
		return fmt.Errorf("error removing favorite: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("favorite not found")
	}

	return nil
}
