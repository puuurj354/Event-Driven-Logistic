# Makefile untuk Event-Driven Logistic System

.PHONY: help docker-up docker-down docker-logs migrate test clean run-order run-payment run-inventory run-notification

# Default target
help:
	@echo "📦 Event-Driven Logistic System - Available Commands:"
	@echo ""
	@echo "  🏠 LOCAL Development (No Docker):"
	@echo "    make setup-local        - Setup local PostgreSQL databases (4 DBs)"
	@echo "    make migrate-local      - Run migrations to local PostgreSQL"
	@echo "    make local-db-status    - Check local database connections"
	@echo "    make run-order          - Run Order Service locally"
	@echo "    make run-payment        - Run Payment Service locally"
	@echo "    make run-inventory      - Run Inventory Service locally"
	@echo "    make run-notification   - Run Notification Service locally"
	@echo ""
	@echo "  🐳 DOCKER Development:"
	@echo "    make docker-up          - Start all Docker containers (RabbitMQ + 4 PostgreSQL DBs)"
	@echo "    make docker-down        - Stop all Docker containers"
	@echo "    make docker-logs        - Show Docker logs"
	@echo "    make docker-clean       - Stop containers and remove volumes"
	@echo "    make migrate-docker     - Run migrations to Docker containers"
	@echo ""
	@echo "  🧪 Testing:"
	@echo "    make test               - Run all tests"
	@echo "    make test-coverage      - Run tests with coverage"
	@echo ""

# ========================================
# LOCAL Development Commands
# ========================================

setup-local:
	@echo "🏠 Setting up local PostgreSQL databases..."
	@chmod +x scripts/setup-local-db.sh
	@./scripts/setup-local-db.sh

migrate-local:
	@echo "🔄 Running migrations to local databases..."
	@chmod +x scripts/migrate-local.sh
	@./scripts/migrate-local.sh

local-db-status:
	@echo "🔍 Checking LOCAL database connections..."
	@echo ""
	@echo "Order DB (localhost:5432):"
	@PGPASSWORD=postgres psql -h localhost -p 5432 -U postgres -d db_order -c "SELECT 'Connected ✅' as status;" 2>/dev/null || echo "❌ Not connected"
	@echo ""
	@echo "Payment DB (localhost:5432):"
	@PGPASSWORD=postgres psql -h localhost -p 5432 -U postgres -d db_payment -c "SELECT 'Connected ✅' as status;" 2>/dev/null || echo "❌ Not connected"
	@echo ""
	@echo "Inventory DB (localhost:5432):"
	@PGPASSWORD=postgres psql -h localhost -p 5432 -U postgres -d db_inventory -c "SELECT 'Connected ✅' as status, COUNT(*) as products FROM products;" 2>/dev/null || echo "❌ Not connected"
	@echo ""
	@echo "Notification DB (localhost:5432):"
	@PGPASSWORD=postgres psql -h localhost -p 5432 -U postgres -d db_notification -c "SELECT 'Connected ✅' as status;" 2>/dev/null || echo "❌ Not connected"

# ========================================
# Docker commands
# ========================================

docker-up:
	@echo "🐳 Starting Docker containers..."
	cd deployments && docker-compose up -d
	@echo "⏳ Waiting for databases to be ready..."
	@sleep 10
	@echo "✅ All containers are up!"
	@echo ""
	@echo "🔗 Access points:"
	@echo "   RabbitMQ Management: http://localhost:15672 (guest/guest)"
	@echo "   PostgreSQL Order:    localhost:5432"
	@echo "   PostgreSQL Payment:  localhost:5433"
	@echo "   PostgreSQL Inventory: localhost:5434"
	@echo "   PostgreSQL Notification: localhost:5435"

docker-down:
	@echo "🛑 Stopping Docker containers..."
	cd deployments && docker-compose down

docker-logs:
	cd deployments && docker-compose logs -f

docker-clean:
	@echo "🧹 Cleaning up Docker containers and volumes..."
	cd deployments && docker-compose down -v
	@echo "✅ Cleanup complete!"

migrate-docker:
	@echo "🔄 Running database migrations to Docker containers..."
	@chmod +x scripts/migrate.sh
	@./scripts/migrate.sh

# Service run commands
run-order:
	@echo "🚀 Running Order Service on port 8081..."
	cd cmd/order-service && go run main.go

run-payment:
	@echo "🚀 Running Payment Service on port 8082..."
	cd cmd/payment-service && go run main.go

run-inventory:
	@echo "🚀 Running Inventory Service on port 8084..."
	cd cmd/inventory-service && go run main.go

run-notification:
	@echo "🚀 Running Notification Service on port 8083..."
	cd cmd/notification-service && go run main.go

# Testing
test:
	@echo "🧪 Running tests..."
	go test -v ./...

test-coverage:
	@echo "🧪 Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Build commands
build-all:
	@echo "🔨 Building all services..."
	@go build -o bin/order-service ./cmd/order-service
	@go build -o bin/payment-service ./cmd/payment-service
	@go build -o bin/inventory-service ./cmd/inventory-service
	@go build -o bin/notification-service ./cmd/notification-service
	@echo "✅ All services built successfully!"

# Clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "✅ Clean complete!"
