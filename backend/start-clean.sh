#!/bin/bash

echo "🚀 Starting Crypto Tower Defense Backend..."
echo "📍 Working directory: $(pwd)"
echo ""

# Kill any existing process on port 8080
echo "🔍 Checking for existing processes on port 8080..."
lsof -ti:8080 | xargs kill -9 2>/dev/null && echo "✅ Killed existing process" || echo "✅ Port 8080 is free"

echo ""
echo "🏗️  Building backend..."
cd cmd/api
go build -o ../../bin/server

if [ $? -ne 0 ]; then
    echo "❌ Build failed!"
    exit 1
fi

echo "✅ Build successful!"
echo ""
echo "🎯 Starting server on port 8080..."
echo "📊 Revenue stats endpoint: http://localhost:8080/api/v1/revenue/stats"
echo "🏥 Health check: http://localhost:8080/health"
echo "🔍 V1 Ping: http://localhost:8080/api/v1/ping"
echo ""
echo "============================================"
echo "Server logs will appear below:"
echo "============================================"
echo ""

cd ../..
./bin/server
