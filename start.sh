#!/bin/bash

# Crypto Tower Defense - SIMPLIFIED Startup Script
# Just starts services, assumes DB is already configured

set -e

PROJECT_ROOT="/Users/lorengraff/Development/1.-EN-DESARROLLO/Crypto_Towell_Defense"
BACKEND_LOG="/tmp/crypto-tower-backend.log"
FRONTEND_LOG="/tmp/crypto-tower-frontend.log"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo "🚀 Crypto Tower Defense - Startup"
echo "========================================"
echo ""

# Function to check port
check_port() {
    lsof -ti:$1 > /dev/null 2>&1
}

# Function to kill port
kill_port() {
    if check_port $1; then
        echo -e "${YELLOW}⚠️  Port $1 in use, killing...${NC}"
        lsof -ti:$1 | xargs kill -9 2>/dev/null || true
        sleep 2
    fi
}

# 1. Backend
echo -e "${BLUE}1️⃣  Backend${NC}"
kill_port 8080

cd "$PROJECT_ROOT/backend"

if [ ! -f "./api" ]; then
    echo -e "   ${YELLOW}Building...${NC}"
    go build -o api ./cmd/api || {
        echo -e "   ${RED}❌ Build failed${NC}"
        exit 1
    }
fi

nohup ./api > "$BACKEND_LOG" 2>&1 &
BACKEND_PID=$!
echo $BACKEND_PID > /tmp/backend.pid

sleep 4
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "   ${GREEN}✅ Running (PID: $BACKEND_PID)${NC}"
    echo -e "   ${GREEN}📝 Log: $BACKEND_LOG${NC}"
else
    echo -e "   ${RED}❌ Failed - check: tail -f $BACKEND_LOG${NC}"
    exit 1
fi

# 2. Frontend
echo -e "${BLUE}2️⃣  Frontend${NC}"
kill_port 3000

cd "$PROJECT_ROOT/game-client"
nohup python3 -m http.server 3000 > "$FRONTEND_LOG" 2>&1 &
FRONTEND_PID=$!
echo $FRONTEND_PID > /tmp/frontend.pid

sleep 2
if curl -s http://localhost:3000/ > /dev/null 2>&1; then
    echo -e "   ${GREEN}✅ Running (PID: $FRONTEND_PID)${NC}"
else
    echo -e "   ${RED}❌ Failed${NC}"
    exit 1
fi

echo ""
echo "========================================"
echo -e "${GREEN}✅ Services started!${NC}"
echo ""
echo -e "${BLUE}📊 URLs:${NC}"
echo "   🎮 Game:     http://localhost:3000/"
echo "   🔧 Backend:  http://localhost:8080/health"
echo "   👥 Admin:    http://localhost:8080/admin/"
echo ""
echo -e "${BLUE}📝 Logs:${NC}"
echo "   Backend:  tail -f $BACKEND_LOG"
echo "   Frontend: tail -f $FRONTEND_LOG"
echo ""
echo -e "${BLUE}🛑 Stop:${NC}"
echo "   ./stop.sh"
echo ""
echo -e "${YELLOW}� Database:${NC}"
echo "   Already configured in backend/.env"
echo "   Migrations: psql -d tower_defense_dev < backend/migrations/001_initial_schema.sql"
echo "   Seed data:  cd backend && go run cmd/seed/main.go"
echo ""
