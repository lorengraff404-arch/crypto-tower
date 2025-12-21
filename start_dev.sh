#!/bin/bash
set -e

echo "🚀 Starting Master Dev Environment..."

# 1. Clean & Build Frontend
echo "📦 Building Frontend..."
cd game-client
npm run clean
npm run build
cd ..

# 2. Reset & Seed Database
echo "🌱 Seeding Database..."
cd backend
# Verify Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed!"
    exit 1
fi
go run cmd/seeder/main.go
cd ..

# 3. Start Backend (in background)
echo "🔥 Starting Backend API..."
cd backend
go run cmd/api/main.go &
BACKEND_PID=$!
cd ..

# 4. Serve Frontend (Python simple server)
echo "🌐 Serving Frontend at http://localhost:8000"
cd game-client
# Check for python3
if command -v python3 &> /dev/null; then
    python3 -m http.server 8000 &
    FRONTEND_PID=$!
else
    echo "⚠️ Python3 not found, skipping frontend server. Open game-client/index.html manually."
fi

echo "✅ Environment Ready!"
echo "   - API: http://localhost:8080"
echo "   - Frontend: http://localhost:8000"
echo "   - Admin: http://localhost:8080/admin"
echo ""
echo "Press [CTRL+C] to stop."

# Wait for user interrupt
trap "kill $BACKEND_PID $FRONTEND_PID; exit" SIGINT
wait
