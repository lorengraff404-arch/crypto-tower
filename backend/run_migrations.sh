#!/bin/bash
echo "🔄 Running all migrations..."
for file in migrations/*.sql; do
    echo "Executing: $file"
    PGPASSWORD=postgres psql -h localhost -U postgres -d crypto_tower_defense -f "$file" 2>&1 | grep -v "already exists" | grep -v "duplicate"
done
echo "✅ All migrations executed"
