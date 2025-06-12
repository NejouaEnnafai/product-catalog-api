.PHONY: build run test seed clean

# Build the application
build:
	go build -o bin/product-catalog-api ./cmd

# Run the application
run:
	go run ./cmd/main.go

# Run the application with database seeding
run-with-seed:
	go run ./scripts/run.go

# Run tests
test:
	go test -v ./...

# Seed the database
seed:
	go run ./scripts/seed.go

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f product_catalog.db
