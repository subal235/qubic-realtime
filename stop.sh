#!/bin/bash

# Stop all Realtime services

echo "🛑 Stopping Realtime Platform services..."
echo ""

# Check what's running
TURBOAUTH_RUNNING=$(ps aux | grep -E "[t]urboauth" | grep -v grep)
TURBOROUTE_RUNNING=$(ps aux | grep -E "[t]urboroute" | grep -v grep)

if [ -z "$TURBOAUTH_RUNNING" ] && [ -z "$TURBOROUTE_RUNNING" ]; then
    echo "ℹ️  No services running"
    exit 0
fi

# Show what will be stopped
echo "📋 Services to stop:"
if [ -n "$TURBOAUTH_RUNNING" ]; then
    TURBOAUTH_PID=$(echo "$TURBOAUTH_RUNNING" | awk '{print $2}')
    echo "  • TurboAuth (PID: $TURBOAUTH_PID)"
fi
if [ -n "$TURBOROUTE_RUNNING" ]; then
    TURBOROUTE_PID=$(echo "$TURBOROUTE_RUNNING" | awk '{print $2}')
    echo "  • TurboRoute (PID: $TURBOROUTE_PID)"
fi

echo ""
echo "⏳ Stopping services..."

# Kill processes
pkill -f "turboauth|turboroute"

# Wait a moment
sleep 1

# Verify they're stopped
STILL_RUNNING=$(ps aux | grep -E "(turboauth|turboroute)" | grep -v grep)

if [ -z "$STILL_RUNNING" ]; then
    echo "✅ All services stopped successfully"
    echo ""
    echo "🔓 Ports are now free:"
    echo "  • 8080 (TurboAuth HTTP)"
    echo "  • 9090 (TurboAuth gRPC)"
    echo "  • 8081 (TurboRoute HTTP)"
    echo "  • 9091 (TurboRoute gRPC)"
else
    echo "⚠️  Some processes still running. Forcing..."
    pkill -9 -f "turboauth|turboroute"
    sleep 1
    echo "✅ Force stopped"
fi

echo ""
echo "🚀 You can now run:"
echo "  • ./start.sh          (background)"
echo "  • make dev-turboauth  (foreground)"
echo "  • make deploy         (Docker)"
echo ""
