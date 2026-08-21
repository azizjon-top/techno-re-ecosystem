.PHONY: help build run test clean docker-up docker-down migrate lint fmt

# Variables
BINARY_NAME=server
GO=go
DOCKER_COMPOSE=docker-compose
PROJECT_NAME=techno-re-ecosystem

# Help command
help:
	@echo "Techno RE Ecosystem - Available Commands"
	@echo "========================================"
	@echo "make build          - Build the application"
	@echo "make run            - Run the application locally"
	@echo "make test           - Run all tests"
	@echo "make test-coverage  - Run tests with coverage report"
	@echo "make clean          - Clean build artifacts"
	@echo "make docker-build   - Build Docker image"
	@echo "make docker-up      - Start services with docker-compose"
	@echo "make docker-down    - Stop services with docker-compose"
	@echo "make docker-logs    - View docker-compose logs"
	@echo "make migrate        - Run database migrations"
	@echo "make lint           - Run linter (golangci-lint)"
	@echo "make fmt            - Format code with gofmt"
	@echo "make deps           - Download dependencies"
	@echo "make dev            - Start development environment"

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build -o bin/$(BINARY_NAME) ./cmd/server

# Run the application locally
run: build
	@echo "Running $(BINARY_NAME)..."
	./bin/$(BINARY_NAME)

# Run all tests
test:
	@echo "Running tests..."
	$(GO) test -v -race -timeout 10m ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -v -race -coverprofile=coverage.out -timeout 10m ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	$(GO) clean

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(PROJECT_NAME):latest .

# Start services with docker-compose
docker-up:
	@echo "Starting services..."
	$(DOCKER_COMPOSE) up -d

# Stop services with docker-compose
docker-down:
	@echo "Stopping services..."
	$(DOCKER_COMPOSE) down

# View docker-compose logs
docker-logs:
	$(DOCKER_COMPOSE) logs -f

# Run database migrations
migrate:
	@echo "Running migrations..."
	# Add migration commands here when migrations are set up

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...
	goimports -w .

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Setup development environment
dev: docker-up
	@echo "Development environment started"
	@echo "Visit http://localhost:8080/health to check server health"
	@echo "Run 'make docker-logs' to view logs"

# Setup environment file
env-setup:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo ".env file created from .env.example"; \
	else \
		echo ".env file already exists"; \
	fi

# Install development tools
install-tools:
	@echo "Installing development tools..."
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install golang.org/x/tools/cmd/goimports@latest

# Run the server in Docker
docker-run: docker-build docker-up
	@echo "Server is running. Access http://localhost:8080"

# Full clean and rebuild
full-rebuild: clean deps build test
	@echo "Full rebuild completed"

# Development hot reload (requires air)
watch:
	@echo "Starting development server with hot reload..."
	air

# Get database shell access
db-shell:
	$(DOCKER_COMPOSE) exec postgres psql -U techno -d techno_re

# Get redis shell access
redis-shell:
	$(DOCKER_COMPOSE) exec redis redis-cli

# View server logs
logs:
	$(DOCKER_COMPOSE) logs -f server

# Reset all containers
reset: docker-down
	@echo "Removing volumes..."
	docker volume rm $$(docker volume ls -q | grep techno)
	@echo "All containers and volumes removed"

# Status check
status:
	@echo "Checking service status..."
	@echo "Server: " && curl -s http://localhost:8080/health || echo "Server not running"
	@echo "\nDatabase: " && $(DOCKER_COMPOSE) exec -T postgres pg_isready -U techno || echo "Database not running"
	@echo "\nCache: " && $(DOCKER_COMPOSE) exec -T redis redis-cli ping || echo "Cache not running"
