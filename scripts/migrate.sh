#!/bin/bash

# Script untuk menjalankan migrations ke semua database
# Pastikan PostgreSQL containers sudah running

set -e

echo "🚀 Running Database Migrations..."
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Database credentials
POSTGRES_USER="postgres"
POSTGRES_PASSWORD="postgres"

# Function to run migration
run_migration() {
    local db_name=$1
    local db_port=$2
    local migration_file=$3
    
    echo -e "${BLUE}📦 Migrating $db_name on port $db_port...${NC}"
    
    # Check if database is accessible
    if ! PGPASSWORD=$POSTGRES_PASSWORD psql -h localhost -p $db_port -U $POSTGRES_USER -d $db_name -c '\q' 2>/dev/null; then
        echo -e "${RED}❌ Cannot connect to $db_name on port $db_port${NC}"
        echo -e "${RED}   Make sure PostgreSQL container is running!${NC}"
        return 1
    fi
    
    # Run migration
    PGPASSWORD=$POSTGRES_PASSWORD psql -h localhost -p $db_port -U $POSTGRES_USER -d $db_name -f $migration_file
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Successfully migrated $db_name${NC}"
    else
        echo -e "${RED}❌ Failed to migrate $db_name${NC}"
        return 1
    fi
    echo ""
}

# Check if PostgreSQL client is installed
if ! command -v psql &> /dev/null; then
    echo -e "${RED}❌ psql (PostgreSQL client) is not installed${NC}"
    echo "   Install it with: sudo apt-get install postgresql-client"
    exit 1
fi

echo "Starting migrations for all services..."
echo "========================================"
echo ""

# Run migrations for each service
run_migration "db_order" "5432" "pkg/migrations/001_order_service.sql"
run_migration "db_inventory" "5434" "pkg/migrations/002_inventory_service.sql"
run_migration "db_payment" "5433" "pkg/migrations/003_payment_service.sql"
run_migration "db_notification" "5435" "pkg/migrations/005_notification_service.sql"

echo "========================================"
echo -e "${GREEN}🎉 All migrations completed successfully!${NC}"
