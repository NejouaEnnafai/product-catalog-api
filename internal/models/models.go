package models

// User represents a user in the system
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // Password is not included in JSON responses
}

// Product represents a product in the catalog
type Product struct {
	ID       int     `json:"id"`
	Title    string  `json:"title"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
	Image    string  `json:"image"`
}

// Favorite represents a user's favorite product
type Favorite struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	ProductID int    `json:"product_id"`
	Notes     string `json:"notes,omitempty"` // Optional notes (bonus feature)
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the login response body
type LoginResponse struct {
	Token string `json:"token"`
}

// FavoriteRequest represents the request to add a favorite
type FavoriteRequest struct {
	ProductID int    `json:"product_id"`
	Notes     string `json:"notes,omitempty"` // Optional notes (bonus feature)
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Total   int         `json:"total"`
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
	Results interface{} `json:"results"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}
