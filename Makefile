.PHONY: help dev build deploy test clean

help: ## Show this help
	@echo 'Realtime Platform - Available Commands'
	@echo ''
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
dev: ## Run services locally (fastest for development)
	@echo "Starting infrastructure (Redis, Prometheus, Grafana)..."
	cd infrastructure && docker-compose up -d redis prometheus grafana
	@echo ""
	@echo "✅ Infrastructure running!"
	@echo "   Redis: localhost:6379"
	@echo "   Prometheus: http://localhost:9092"
	@echo "   Grafana: http://localhost:3000"
	@echo ""
	@echo "Now run services:"
	@echo "   Terminal 1: cd services/turboauth/backend && go run ./cmd/api"
	@echo "   Terminal 2: cd services/turboroute/backend && go run ./cmd/api"

dev-turboauth: ## Run TurboAuth locally
	cd services/turboauth/backend && go run ./cmd/api

dev-turboroute: ## Run TurboRoute locally
	cd services/turboroute/backend && go run ./cmd/api

# Build
build: ## Build all services
	@echo "Building TurboAuth..."
	cd services/turboauth/backend && go build -o turboauth ./cmd/api
	@echo "Building TurboRoute..."
	cd services/turboroute/backend && go build -o turboroute ./cmd/api
	@echo "✅ All services built"

build-turboauth: ## Build TurboAuth
	cd services/turboauth/backend && go build -o turboauth ./cmd/api

build-turboroute: ## Build TurboRoute
	cd services/turboroute/backend && go build -o turboroute ./cmd/api

# Docker
docker-build: ## Build Docker images
	cd infrastructure && docker-compose build

deploy: ## Deploy all services with Docker
	cd infrastructure && docker-compose up -d
	@echo "✅ All services deployed!"
	@echo "   TurboAuth: http://localhost:8080"
	@echo "   TurboRoute: http://localhost:8081"
	@echo "   Prometheus: http://localhost:9092"
	@echo "   Grafana: http://localhost:3000"

deploy-turboauth: ## Deploy only TurboAuth
	cd infrastructure && docker-compose up -d turboauth redis

deploy-turboroute: ## Deploy only TurboRoute
	cd infrastructure && docker-compose up -d turboroute redis

stop: ## Stop all services
	cd infrastructure && docker-compose down

restart: ## Restart all services
	cd infrastructure && docker-compose restart

# Testing
test: ## Run all tests
	@echo "Testing TurboAuth..."
	cd services/turboauth/backend && go test -v ./...
	@echo "Testing TurboRoute..."
	cd services/turboroute/backend && go test -v ./...

test-turboauth: ## Test TurboAuth
	cd services/turboauth/backend && go test -v ./...

test-turboroute: ## Test TurboRoute
	cd services/turboroute/backend && go test -v ./...

# Utilities
logs: ## View all logs
	cd infrastructure && docker-compose logs -f

logs-turboauth: ## View TurboAuth logs
	cd infrastructure && docker-compose logs -f turboauth

logs-turboroute: ## View TurboRoute logs
	cd infrastructure && docker-compose logs -f turboroute

health: ## Check health of all services
	@echo "TurboAuth:"
	@curl -s http://localhost:8080/health | jq . || echo "Not running"
	@echo "\nTurboRoute:"
	@curl -s http://localhost:8081/health | jq . || echo "Not running"

ps: ## Show running processes
	@echo "Running processes:"
	@ps aux | grep -E "(turboauth|turboroute)" | grep -v grep || echo "No processes running"
	@echo "\nPorts in use:"
	@lsof -i :8080 2>/dev/null | grep LISTEN || echo "Port 8080: free"
	@lsof -i :8081 2>/dev/null | grep LISTEN || echo "Port 8081: free"
	@lsof -i :9090 2>/dev/null | grep LISTEN || echo "Port 9090: free"
	@lsof -i :9091 2>/dev/null | grep LISTEN || echo "Port 9091: free"

kill: ## Kill all running local processes
	@echo "Killing local processes..."
	@pkill -f "turboauth|turboroute" || echo "No processes to kill"
	@echo "✅ Done"

check-ports: ## Check which ports are in use
	@./check-ports.sh

# Cleanup
clean: ## Clean build artifacts
	rm -f services/turboauth/backend/turboauth
	rm -f services/turboroute/backend/turboroute
	cd infrastructure && docker-compose down -v

# Setup
init: ## Initialize project (first time setup)
	@echo "Installing Go dependencies..."
	cd services/turboauth/backend && go mod download
	cd services/turboroute/backend && go mod download
	@echo "✅ Initialization complete"

tidy: ## Tidy Go modules
	cd services/turboauth/backend && go mod tidy
	cd services/turboroute/backend && go mod tidy

# Quick commands
run: dev ## Alias for dev

up: deploy ## Alias for deploy

down: stop ## Alias for stop

start-local: ## Start services with start.sh
	@./start.sh

stop-local: ## Stop services started by start.sh
	@./stop.sh
