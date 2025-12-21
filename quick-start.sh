#!/bin/bash

# Quick Start Script - Crypto Tower Defense
echo "🚀 Starting Crypto Tower Defense Development Environment"
echo "========================================================"

# Kill existing processes
echo "🧹 Cleaning up existing processes..."
lsof -ti:8080 | xargs kill -9 2>/dev/null
lsof -ti:8545 | xargs kill -9 2>/dev/null
sleep 2

# Start Anvil (blockchain)
echo "⛓️  Starting Anvil local blockchain..."
anvil --port 8545 --accounts 10 --balance 10000 > /tmp/anvil.log 2>&1 &
ANVIL_PID=$!
echo "   Anvil PID: $ANVIL_PID"
sleep 3

# Deploy contracts (if needed)
if [ ! -f "smart-contracts/.deployments.env" ]; then
    echo "📜 Deploying smart contracts..."
    cd smart-contracts
    export PRIVATE_KEY="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
    forge script script/Deploy.s.sol --rpc-url http://localhost:8545 --broadcast --legacy | tail -20
    cd ..
fi

# Start Backend
echo "🔧 Starting Go backend API..."
cd backend
./api > /tmp/backend.log 2>&1 &
BACKEND_PID=$!
echo "   Backend PID: $BACKEND_PID"
cd ..
sleep 4

# Check status
echo ""
echo "✅ Status Check:"
echo "========================================================"
curl -s http://localhost:8080/health | jq '.' || echo "❌ Backend not responding"

echo ""
echo "🎮 Crypto Tower Defense is Ready!"
echo "========================================================"
echo "Backend API: http://localhost:8080"
echo "Admin Dashboard: http://localhost:8080/admin/"
echo "Health Check: http://localhost:8080/health"
echo ""
echo "Blockchain RPC: http://localhost:8545"
echo "Contracts: smart-contracts/.deployments.env"
echo ""
echo "📊 Logs:"
echo "  Backend: tail -f /tmp/backend.log"
echo "  Anvil: tail -f /tmp/anvil.log"
echo ""
echo "🛑 To stop:"
echo "  kill $BACKEND_PID $ANVIL_PID"
echo ""
echo "🌐 Open http://localhost:8080/admin/ in your browser!"
