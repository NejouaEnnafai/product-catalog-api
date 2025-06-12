# Product Catalog API

A RESTful API backend in Go that powers a mobile product catalog. The API supports product listing, filtering, user authentication using JWT, and favoriting functionality per user.

## Tech Stack

- Go (standard library + JWT only)
- SQLite (with database/sql)
- JSON API responses
- No web frameworks

## Features

- User authentication with JWT
- Product listing with filtering, sorting, and pagination
- User-specific product favorites
- Search functionality

## API Endpoints

1. **POST /login**
   - Body: `{ "username": "john", "password": "1234" }`
   - Returns a JWT token

2. **GET /products**
   - Public route
   - Supports query params:
     - page, limit
     - category
     - sort=price_asc | price_desc
     - search (search in product title)

3. **POST /favorites**
   - Protected route (Authorization: Bearer <token>)
   - Body: `{ "product_id": 123 }`

4. **GET /favorites**
   - Protected route
   - Returns user's favorite products

## Getting Started

1. Clone the repository
2. Run `go mod tidy` to install dependencies
3. Run `go run cmd/main.go` to start the server

## Testing

Run `go test ./...` to run all tests
