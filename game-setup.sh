#!/bin/bash
set -e

# Crypto Tower Defense - System Setup & Hardening Script
# Ensures definitive database consistency and zero-error restarts.

echo -e "\033[0;34m🛡️  Crypto Tower Defense - System Setup & Hardening\033[0m"
echo "================================================="

# Check if Postgres is reachable
pg_isready -h localhost -p 5432 > /dev/null 2>&1 || {
    echo -e "\033[0;31m❌ Postgres is not running on port 5432.\033[0m"
    echo "Please start Postgres or check connection."
    exit 1
}

echo -e "\033[1;33m❓ Choose initialization mode:\033[0m"
echo "   1) [RESET] Drop DB, Recreate, Import Golden Schema, Seed (DATA LOSS)"
echo "   2) [VERIFY] Run Seeders only (Safe for existing data)"
read -p "   Enter choice [1/2]: " choice

if [ "$choice" == "1" ]; then
    echo ""
    echo -e "\033[0;31m⚠️  WARNING: This will DESTROY 'tower_defense_dev' database.\033[0m"
    read -p "   Are you sure? (y/n): " confirm
    if [ "$confirm" == "y" ]; then
        echo ""
        echo -e "\033[0;34m🗑️  Dropping database...\033[0m"
        dropdb tower_defense_dev --if-exists
        
        echo -e "\033[0;34m✨ Creating database...\033[0m"
        createdb tower_defense_dev
        
        echo -e "\033[0;34m📜 Importing Golden Schema...\033[0m"
        psql -d tower_defense_dev < backend/migrations/001_initial_schema.sql
        
        echo -e "\033[0;34m🌱 Seeding data...\033[0m"
        cd backend && go run cmd/seed/main.go
        
        echo ""
        echo -e "\033[0;32m✅ Reset Complete! System is now consistent.\033[0m"
    else
        echo "🚫 Cancelled."
        exit 0
    fi
elif [ "$choice" == "2" ]; then
    echo ""
    echo -e "\033[0;34m🌱 Running Seeders (Idempotent Check)...\033[0m"
    cd backend && go run cmd/seed/main.go
    echo ""
    echo -e "\033[0;32m✅ Verification Complete!\033[0m"
else
    echo -e "\033[0;31m❌ Invalid choice.\033[0m"
    exit 1
fi
