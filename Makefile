.PHONY: run build test clean docker-build help

# Variables
BINARY_NAME=peso
DOCKER_IMAGE=peso:latest

# Default target
help: ## Show this help
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

run: ## Run the application in development mode
	@echo "🚀 Starting Peso application..."
	go run cmd/main.go

build: ## Build the application binary
	@echo "🔨 Building application..."
	go build -o bin/$(BINARY_NAME) cmd/main.go
	@echo "✅ Built bin/$(BINARY_NAME)"

test: ## Run all tests
	@echo "🧪 Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	go test -cover ./...

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	rm -f bin/$(BINARY_NAME)
	rm -f peso.db
	rm -f *.db

docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

lint: ## Run golangci-lint (if available)
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "🔍 Running linter..."; \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not found, running go fmt instead"; \
		go fmt ./...; \
	fi

fmt: ## Format code
	@echo "✨ Formatting code..."
	go fmt ./...

tidy: ## Tidy go modules
	@echo "📦 Tidying modules..."
	go mod tidy

dev-setup: tidy fmt ## Set up development environment
	@echo "🛠️  Development environment ready!"

# Quick commands for common operations
up: build run ## Build and run
test-domain: ## Run domain tests only
	@echo "🧪 Testing domain layer..."
	go test -v ./internal/domain/...