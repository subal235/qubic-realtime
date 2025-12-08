.PHONY: help dev build deploy test proto clean

help: ## Show this help message
	@echo 'Realtime Monorepo - Available Commands'
	@echo ''
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
dev: ## Start all services in development mode
	cd infrastructure && docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

dev-turboauth: ## Start only MicroAuth in development mode
	cd services/turboauth/backend && go run ./cmd/api

# Build
build: ## Build all services
	@echo "Building MicroAuth..."
	cd services/turboauth/backend && go build -o turboauth ./cmd/api
	@echo "✅ All services built"

build-turboauth: ## Build MicroAuth service
	cd services/turboauth/backend && go build -o turboauth ./cmd/api

build-docker: ## Build all Docker images
	cd infrastructure && docker-compose build

# Deploy
deploy: ## Deploy all services in production mode
	cd infrastructure && docker-compose up -d

deploy-turboauth: ## Deploy only MicroAuth
	cd infrastructure && docker-compose up -d turboauth

stop: ## Stop all services
	cd infrastructure && docker-compose down

restart: ## Restart all services
	cd infrastructure && docker-compose restart

# Testing
test: ## Run all tests
	@echo "Testing MicroAuth..."
	cd services/turboauth/backend && go test -v ./...

test-coverage: ## Run tests with coverage
	cd services/turboauth/backend && go test -v -race -coverprofile=coverage.out ./...
	cd services/turboauth/backend && go tool cover -html=coverage.out

test-integration: ## Run integration tests
	cd services/turboauth/backend && go test -v -tags=integration ./...

# Proto
proto: ## Generate gRPC code from proto files
	@echo "Generating proto files..."
	protoc --go_out=services/turboauth/backend --go_opt=paths=source_relative \
		--go-grpc_out=services/turboauth/backend --go-grpc_opt=paths=source_relative \
		services/turboauth/api/proto/*.proto
	@echo "✅ Proto files generated"

# Utilities
logs: ## View logs from all services
	cd infrastructure && docker-compose logs -f

logs-turboauth: ## View MicroAuth logs
	cd infrastructure && docker-compose logs -f turboauth

health: ## Check health of all services
	@echo "Checking MicroAuth..."
	@curl -s http://localhost:8080/health | jq . || echo "MicroAuth not running"

metrics: ## View metrics
	@curl -s http://localhost:2112/metrics | grep turboauth | head -20

# Cleanup
clean: ## Clean build artifacts
	rm -f services/turboauth/backend/turboauth
	rm -rf services/turboauth/backend/api/proto/*.pb.go
	cd infrastructure && docker-compose down -v

clean-all: clean ## Clean everything including Docker volumes
	docker system prune -af --volumes

# Lint & Format
lint: ## Run linters
	cd services/turboauth/backend && golangci-lint run

fmt: ## Format code
	cd services/turboauth/backend && go fmt ./...

tidy: ## Tidy go modules
	cd services/turboauth/backend && go mod tidy

# Initialize
init: ## Initialize project (install dependencies)
	@echo "Installing Go dependencies..."
	cd services/turboauth/backend && go mod download
	@echo "Installing protoc plugins..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "✅ Initialization complete"

# Benchmarks
benchmark: ## Run benchmarks
	cd services/turboauth/backend && go test -bench=. -benchmem ./...

# Docker
docker-rebuild: ## Rebuild Docker images from scratch
	cd infrastructure && docker-compose build --no-cache

docker-clean: ## Clean Docker resources
	cd infrastructure && docker-compose down -v
	docker system prune -f
