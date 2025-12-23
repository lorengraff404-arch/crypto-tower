#!/bin/bash

# Stop all Crypto Tower Defense services

echo "🛑 Stopping Crypto Tower Defense..."

# Kill backend
if [ -f /tmp/backend.pid ]; then
    BACKEND_PID=$(cat /tmp/backend.pid)
    if kill $BACKEND_PID 2>/dev/null; then
        echo "✅ Backend stopped (PID: $BACKEND_PID)"
    fi
    rm /tmp/backend.pid
fi

# Kill frontend
if [ -f /tmp/frontend.pid ]; then
    FRONTEND_PID=$(cat /tmp/frontend.pid)
    if kill $FRONTEND_PID 2>/dev/null; then
        echo "✅ Frontend stopped (PID: $FRONTEND_PID)"
    fi
    rm /tmp/frontend.pid
fi

# Kill any stragglers on ports
lsof -ti:8080 | xargs kill -9 2>/dev/null
lsof -ti:3000 | xargs kill -9 2>/dev/null

echo "✅ All services stopped"
