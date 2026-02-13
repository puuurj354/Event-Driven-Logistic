#!/bin/bash

# Script untuk membuat semua database dan menjalankan migrasi
# Usage: ./setup_databases.sh

echo "🗄️  Setting up databases for Event-Driven-Logistic System"
echo ""

# Database credentials
DB_USER="postgres"
DB_PASS="postgres"
DB_HOST="localhost"
DB_PORT="5432"

# Function to create database if it doesn't exist
create_database() {
    local db_name=$1
    echo "📦 Creating database: $db_name"
    
    # Check if database exists
    PGPASSWORD=$DB_PASS psql -h $DB_HOST -p $DB_PORT -U $DB_USER -tc "SELECT 1 FROM pg_database WHERE datname = '$db_name'" | grep -q 1
    
    if [ $? -eq 0 ]; then
        echo "   ⚠️  Database $db_name already exists, dropping and recreating..."
        PGPASSWORD=$DB_PASS psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "DROP DATABASE IF EXISTS $db_name;"
    fi
    
    # Create database
    PGPASSWORD=$DB_PASS psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "CREATE DATABASE $db_name;"
    
    # Enable UUID extension
    echo "   🔧 Enabling UUID extension..."
    PGPASSWORD=$DB_PASS psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $db_name -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"
    
    echo "   ✅ Database $db_name ready"
    echo ""
}

# Create all databases
create_database "db_order"
create_database "db_payment"
create_database "db_inventory"
create_database "db_delivery"

echo "✨ All databases created successfully!"
echo ""
echo "📋 Next steps:"
echo "   1. Run: make run-order       # Start Order Service (port 8081)"
echo "   2. Run: make run-payment     # Start Payment Service (port 8082)"
echo "   3. Run: make run-inventory   # Start Inventory Service (port 8084)"
echo "   4. Run: cd cmd/delivery-service && go run main.go  # Start Delivery Service (port 8085)"
echo ""
echo "Each service will automatically create its tables via AutoMigrate"
