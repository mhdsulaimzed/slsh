.PHONY: build run test clean install dev help

# Variables
BINARY_NAME=slsh
BUILD_DIR=bin
MAIN_FILE=main.go

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Built $(BUILD_DIR)/$(BINARY_NAME)"

# Run the shell directly
run:
	@echo "Starting $(BINARY_NAME)..."
	go run $(MAIN_FILE)

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Install to system (requires sudo)
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin/"
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installed successfully!"

# Development mode with auto-rebuild
dev:
	@echo "Starting development mode..."
	@echo "Note: Install 'air' for auto-reload: go install github.com/cosmtrek/air@latest"
	@which air > /dev/null && air || go run $(MAIN_FILE)

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@which golangci-lint > /dev/null && golangci-lint run || echo "golangci-lint not found, skipping..."

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Check for vulnerabilities
vuln-check:
	@echo "Checking for vulnerabilities..."
	@which govulncheck > /dev/null && govulncheck ./... || echo "govulncheck not found, install with: go install golang.org/x/vuln/cmd/govulncheck@latest"

# Release build with optimizations
release: clean
	@echo "Building release version..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Release build completed: $(BUILD_DIR)/$(BINARY_NAME)"

# Help target
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  run          - Run the shell directly"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  clean        - Clean build artifacts"
	@echo "  install      - Install to system (requires sudo)"
	@echo "  dev          - Development mode with auto-rebuild"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code (requires golangci-lint)"
	@echo "  tidy         - Tidy dependencies"
	@echo "  vuln-check   - Check for vulnerabilities"
	@echo "  release      - Build optimized release version"
	@echo "  help         - Show this help"