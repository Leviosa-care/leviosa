# Notification Service Makefile

.PHONY: help build run test clean deps

# Default target
help:
	@echo "Available targets:"
	@echo "  build    - Build the notification service binary"
	@echo "  run      - Run the notification service"
	@echo "  test     - Run tests"
	@echo "  clean    - Clean build artifacts"
	@echo "  deps     - Download and tidy dependencies"

# Build the service
build:
	go build -o bin/notification main.go

# Run the service
run:
	go run main.go

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Download and tidy dependencies
deps:
	go mod download
	go mod tidy

# Run with development environment
dev: deps
	go run main.go