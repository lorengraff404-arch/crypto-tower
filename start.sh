#!/bin/bash

# Crypto Tower Defense - Startup Script
# This script ensures all services start correctly with proper error handling

set -e  # Exit on any error

PROJECT_ROOT="/Users/lorengraff/Development/1.-EN-DESARROLLO/Crypto_Towell_Defense"
BACKEND_LOG="/tmp/crypto-tower-backend.log"
FRONTEND_LOG="/tmp/crypto-tower-frontend.log"

echo "🚀 Starting Crypto Tower Defense..."
echo "======================================"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to check if port is in use
check_port() {
    lsof -ti:$1 > /dev/null 2>&1
}

# Function to kill process on port
kill_port() {
    if check_port $1; then
        echo -e "${YELLOW}⚠️  Port $1 is in use, killing process...${NC}"
        lsof -ti:$1 | xargs kill -9 2>/dev/null || true
        sleep 2
    fi
}

# 1. PostgreSQL Check
echo -n "1️⃣  Checking PostgreSQL... "
if pgrep -x postgres > /dev/null; then
    echo -e "${GREEN}✅ Running${NC}"
else
    echo -e "${RED}❌ Not running${NC}"
    echo "   Starting PostgreSQL..."
    brew services start postgresql@14
    sleep 3
fi

# Verify database exists
if psql -d tower_defense_dev -c '\q' 2>/dev/null; then
    echo -e "   Database: ${GREEN}✅ Connected${NC}"
else
    echo -e "   ${RED}❌ Database not accessible${NC}"
    echo "   Run: createdb tower_defense_dev"
    exit 1
fi

# 2. Check Anvil (optional for local blockchain)
echo -n "2️⃣  Checking Anvil... "
if check_port 8545; then
    echo -e "${GREEN}✅ Running${NC}"
else
    echo -e "${YELLOW}⚠️  Not running (optional)${NC}"
    echo "   To start: anvil --port 8545"
fi

# 3. Backend
echo -n "3️⃣  Starting Backend... "
kill_port 8080

cd "$PROJECT_ROOT/backend"

# Check if binary exists
if [ ! -f "./api" ]; then
    echo -e "${YELLOW}Binary not found, building...${NC}"
    go build -o api ./cmd/api || {
        echo -e "${RED}❌ Build failed${NC}"
        exit 1
    }
fi

# Start backend
nohup ./api > "$BACKEND_LOG" 2>&1 &
BACKEND_PID=$!
echo $BACKEND_PID > /tmp/backend.pid

# Wait and verify
sleep 4
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Running (PID: $BACKEND_PID)${NC}"
    echo "   Log: $BACKEND_LOG"
else
    echo -e "${RED}❌ Failed to start${NC}"
    echo "   Check log: tail -f $BACKEND_LOG"
    exit 1
fi

# 4. Frontend
echo -n "4️⃣  Starting Frontend... "
kill_port 3000

cd "$PROJECT_ROOT/game-client"

nohup python3 -m http.server 3000 > "$FRONTEND_LOG" 2>&1 &
FRONTEND_PID=$!
echo $FRONTEND_PID > /tmp/frontend.pid

sleep 2
if curl -s http://localhost:3000/ > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Running (PID: $FRONTEND_PID)${NC}"
else
    echo -e "${RED}❌ Failed to start${NC}"
    exit 1
fi

echo ""
echo "======================================"
echo -e "${GREEN}✅ All services started!${NC}"
echo ""
echo "📊 URLs:"
echo "   Game Client:  http://localhost:3000/"
echo "   Backend API:  http://localhost:8080/health"
echo "   Admin Panel:  http://localhost:8080/admin/"
echo ""
echo "📝 Logs:"
echo "   Backend:  tail -f $BACKEND_LOG"
echo "   Frontend: tail -f $FRONTEND_LOG"
echo ""
echo "🛑 To stop all services:"
echo "   kill \$(cat /tmp/backend.pid) \$(cat /tmp/frontend.pid)"
echo ""
