#!/bin/bash
# Crypto Tower Defense - Backend Management Script
# USE THIS to start/stop/restart backend reliably

set -e

PROJECT_DIR="/Users/lorengraff/Development/1.-EN-DESARROLLO/Crypto_Towell_Defense"
BACKEND_DIR="$PROJECT_DIR/backend"
BACKEND_BIN="$BACKEND_DIR/api"
BACKEND_LOG="/tmp/crypto_backend.log"
BACKEND_PID="/tmp/crypto_backend.pid"

case "${1:-start}" in
  build)
    echo "🔨 Building backend..."
    cd "$BACKEND_DIR"
    go build -o api ./cmd/api || exit 1
    echo "✅ Built: $(ls -lh api | awk '{print $5}')"
    ;;
    
  start)
    if [ -f "$BACKEND_PID" ] && kill -0 $(cat "$BACKEND_PID") 2>/dev/null; then
      echo "⚠️  Backend already running (PID: $(cat $BACKEND_PID))"
      exit 0
    fi
    
    echo "🚀 Starting backend..."
    cd "$BACKEND_DIR"
    nohup ./api > "$BACKEND_LOG" 2>&1 &
    echo $! > "$BACKEND_PID"
    sleep 3
    
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
      echo "✅ Backend running (PID: $(cat $BACKEND_PID))"
      curl -s http://localhost:8080/health | jq '.'
    else
      echo "❌ Backend failed to start. Log:"
      tail -20 "$BACKEND_LOG"
      exit 1
    fi
    ;;
    
  stop)
    if [ -f "$BACKEND_PID" ]; then
      PID=$(cat "$BACKEND_PID")
      echo "🛑 Stopping backend (PID: $PID)..."
      kill $PID 2>/dev/null || true
      rm -f "$BACKEND_PID"
      echo "✅ Stopped"
    else
      echo "No PID file found"
    fi
    ;;
    
  restart)
    $0 stop
    sleep 2
    $0 build
    $0 start
    ;;
    
  status)
    if [ -f "$BACKEND_PID" ] && kill -0 $(cat "$BACKEND_PID") 2>/dev/null; then
      echo "✅ Running (PID: $(cat $BACKEND_PID))"
      curl -s http://localhost:8080/health | jq '.'
    else
      echo "❌ Not running"
      exit 1
    fi
    ;;
    
  logs)
    tail -f "$BACKEND_LOG"
    ;;
    
  *)
    echo "Usage: $0 {build|start|stop|restart|status|logs}"
    exit 1
    ;;
esac
