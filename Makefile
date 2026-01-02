.PHONY: help build run test clean proto docker-build docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building GoAuth..."
	@go build -o bin/goauth cmd/api/main.go
	@echo "Build complete: bin/goauth"

run: ## Run the application
	@echo "Starting GoAuth server..."
	@go run cmd/api/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f *.log
	@echo "Clean complete"

proto: ## Regenerate protocol buffers
	@echo "Regenerating protobuf files..."
	@protoc --go_out=. --go-grpc_out=. proto/main.proto
	@echo "Protobuf generation complete"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker-compose build
	@echo "Docker build complete"

docker-up: ## Start services with Docker Compose
	@echo "Starting services..."
	@docker-compose up -d
	@echo "Services started"

docker-down: ## Stop services
	@echo "Stopping services..."
	@docker-compose down
	@echo "Services stopped"

docker-logs: ## View Docker logs
	@docker-compose logs -f

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	@echo "Tools installed"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated"
