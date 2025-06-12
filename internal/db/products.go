package db

import (
	"fmt"
	"strings"

	"github.com/najwa/product-catalog-api/internal/models"
)

// GetProducts retrieves products with filtering, sorting, and pagination
func GetProducts(page, limit int, category, sort, search string) ([]models.Product, int, error) {
	// Build the query
	query := "SELECT id, title, price, category, image FROM products"
	countQuery := "SELECT COUNT(*) FROM products"
	
	// Add WHERE clauses
	whereClause := []string{}
	args := []interface{}{}
	
	if category != "" {
		whereClause = append(whereClause, "category = ?")
		args = append(args, category)
	}
	
	if search != "" {
		whereClause = append(whereClause, "title LIKE ?")
		args = append(args, "%"+search+"%")
	}
	
	if len(whereClause) > 0 {
		whereStr := " WHERE " + strings.Join(whereClause, " AND ")
		query += whereStr
		countQuery += whereStr
	}
	
	// Add ORDER BY clause
	if sort != "" {
		switch sort {
		case "price_asc":
			query += " ORDER BY price ASC"
		case "price_desc":
			query += " ORDER BY price DESC"
		default:
			// Default sort by ID
			query += " ORDER BY id ASC"
		}
	} else {
		// Default sort by ID
		query += " ORDER BY id ASC"
	}
	
	// Add LIMIT and OFFSET for pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	
	offset := (page - 1) * limit
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	
	// Execute the count query
	var total int
	err := DB.QueryRow(countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting products: %w", err)
	}
	
	// Execute the main query
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying products: %w", err)
	}
	defer rows.Close()
	
	// Parse the results
	products := []models.Product{}
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.Title,
			&product.Price,
			&product.Category,
			&product.Image,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning product: %w", err)
		}
		products = append(products, product)
	}
	
	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating products: %w", err)
	}
	
	return products, total, nil
}

// GetProductByID retrieves a product by ID
func GetProductByID(id int) (*models.Product, error) {
	var product models.Product
	err := DB.QueryRow("SELECT id, title, price, category, image FROM products WHERE id = ?", id).Scan(
		&product.ID,
		&product.Title,
		&product.Price,
		&product.Category,
		&product.Image,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying product: %w", err)
	}
	return &product, nil
}
