#!/bin/bash

# Crypto Tower Defense - Backend Development Server

echo "🚀 Starting Crypto Tower Defense Backend..."

# Kill any existing processes on port 8080
echo "Cleaning up port 8080..."
lsof -ti:8080 | xargs kill -9 2>/dev/null || true
sleep 1

# Set database credentials
export DB_USER=lorengraff
export DB_PASSWORD=""

# Navigate to backend directory
cd "$(dirname "$0")/backend"

# Start server
echo "Starting Go server..."
go run cmd/api/main.go
