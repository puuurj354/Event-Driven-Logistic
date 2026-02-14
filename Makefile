# Makefile untuk Event-Driven Logistic System

.PHONY: help docker-up docker-down docker-logs migrate test clean run-order run-payment run-inventory run-notification run-delivery run-all gateway-start gateway-stop gateway-reload

# Default target
help:
	@echo "📦 Event-Driven Logistic System - Available Commands:"
	@echo ""
	@echo "  🏠 LOCAL Development (No Docker):"
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

run-delivery:
	@echo "🚀 Running Delivery Service on port 8085..."
	cd cmd/delivery-service && go run main.go

run-all:
	@echo "🚀 Starting ALL services concurrently..."
	@$(MAKE) run-order & \
	 $(MAKE) run-payment & \
	 $(MAKE) run-inventory & \
	 $(MAKE) run-delivery & \
	 $(MAKE) run-notification & \
	 wait

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

# ── API Gateway (Nginx Reverse Proxy) ──
# Menjalankan Nginx sebagai reverse proxy pada port 8080

# Start Nginx gateway (local, non-root)
gateway-start:
	@echo "🚪 Starting API Gateway (Nginx) on port 8080..."
	@nginx -c $(PWD)/deployments/nginx.conf
	@echo "✅ API Gateway started: http://localhost:8080"

# Stop Nginx gateway
gateway-stop:
	@echo "🛑 Stopping API Gateway..."
	@nginx -c $(PWD)/deployments/nginx.conf -s quit 2>/dev/null || true
	@echo "✅ API Gateway stopped"

# Reload Nginx config (tanpa downtime)
gateway-reload:
	@echo "🔄 Reloading API Gateway config..."
	@nginx -c $(PWD)/deployments/nginx.conf -s reload
	@echo "✅ API Gateway config reloaded"

# Seed inventory with test products
seed-products:
	@echo "🌱 Seeding inventory products..."
	@curl -s -X POST http://localhost:8080/api/products \
		-H "Content-Type: application/json" \
		-d '{"name":"Laptop ASUS ROG","stock":100}' | jq .
	@curl -s -X POST http://localhost:8080/api/products \
		-H "Content-Type: application/json" \
		-d '{"name":"iPhone 15 Pro","stock":50}' | jq .
	@curl -s -X POST http://localhost:8080/api/products \
		-H "Content-Type: application/json" \
		-d '{"name":"Samsung Galaxy S24","stock":75}' | jq .
	@echo "✅ Products seeded!"

