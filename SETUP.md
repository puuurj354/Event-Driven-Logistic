# Quick Setup Guide

## Prerequisites
- Go 1.25.4+
- Docker & Docker Compose
- PostgreSQL Client (psql)
- Make (optional, untuk shortcut commands)

## 🚀 Quick Start (Development)

### 1. Start Infrastructure (RabbitMQ + 4 PostgreSQL DBs)

```bash
# Menggunakan Make
make docker-up

# Atau tanpa Make
cd deployments
docker-compose up -d
```

Ini akan menjalankan:
- **RabbitMQ** di port 5672 (Management UI: http://localhost:15672)
- **PostgreSQL Order DB** di port 5432
- **PostgreSQL Payment DB** di port 5433
- **PostgreSQL Inventory DB** di port 5434
- **PostgreSQL Notification DB** di port 5435

### 2. Run Database Migrations

```bash
# Menggunakan Make
make migrate

# Atau manual
chmod +x scripts/migrate.sh
./scripts/migrate.sh
```

### 3. Verify Database Connections

```bash
make db-status
```

### 4. Run Services Locally

Buka 4 terminal terpisah dan jalankan:

**Terminal 1 - Order Service:**
```bash
make run-order
# atau
cd cmd/order-service && go run main.go
```

**Terminal 2 - Payment Service:**
```bash
make run-payment
```

**Terminal 3 - Inventory Service:**
```bash
make run-inventory
```

**Terminal 4 - Notification Service:**
```bash
make run-notification
```

## 🗄️ Database Architecture

Setiap service memiliki database PostgreSQL sendiri:

| Service       | Database Name      | Port  | Tables                          |
|---------------|-------------------|-------|---------------------------------|
| Order         | db_order          | 5432  | orders                          |
| Payment       | db_payment        | 5433  | payments                        |
| Inventory     | db_inventory      | 5434  | products, stock_reservations    |
| Notification  | db_notification   | 5435  | notifications, websocket_sessions|

### Connection Strings

**Development (Local):**
```
Order:        postgres://postgres:postgres@localhost:5432/db_order?sslmode=disable
Payment:      postgres://postgres:postgres@localhost:5433/db_payment?sslmode=disable
Inventory:    postgres://postgres:postgres@localhost:5434/db_inventory?sslmode=disable
Notification: postgres://postgres:postgres@localhost:5435/db_notification?sslmode=disable
```

**Docker (Container to Container):**
```
Order:        postgres://postgres:postgres@postgres-order:5432/db_order?sslmode=disable
Payment:      postgres://postgres:postgres@postgres-payment:5432/db_payment?sslmode=disable
Inventory:    postgres://postgres:postgres@postgres-inventory:5432/db_inventory?sslmode=disable
Notification: postgres://postgres:postgres@postgres-notification:5432/db_notification?sslmode=disable
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage
```

## 🛑 Stop & Clean

```bash
# Stop containers
make docker-down

# Stop and remove volumes (⚠️ akan menghapus data!)
make docker-clean
```

## 📝 Available Make Commands

Run `make help` untuk melihat semua commands yang tersedia.

## 🔧 Troubleshooting

### Database connection refused
```bash
# Check if containers are running
docker ps

# Check database logs
docker logs postgres-order
docker logs postgres-payment
docker logs postgres-inventory
docker logs postgres-notification
```

### Migration failed
```bash
# Connect manually to check
PGPASSWORD=postgres psql -h localhost -p 5432 -U postgres -d db_order

# Re-run specific migration
PGPASSWORD=postgres psql -h localhost -p 5432 -U postgres -d db_order -f pkg/migrations/001_order_service.sql
```

## 📚 Next Steps

1. Implement business logic di `internal/*/service/`
2. Implement repository layer di `internal/*/repository/`
3. Implement HTTP handlers di `internal/*/delivery/`
4. Implement event publishers/subscribers di `internal/*/event/`
