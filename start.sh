#!/bin/bash

# Realtime Services Startup Script

echo "🚀 Starting Realtime Platform..."
echo ""

# Check if infrastructure is running
echo "📊 Checking infrastructure..."
if ! docker ps | grep -q realtime-redis; then
    echo "⚠️  Infrastructure not running. Starting..."
    cd infrastructure && docker-compose up -d redis prometheus grafana
    cd ..
    sleep 3
else
    echo "✅ Infrastructure already running"
fi

echo ""
echo "🔐 Starting TurboAuth on port 8080..."
cd services/turboauth/backend
./turboauth &
TURBOAUTH_PID=$!
cd ../../..

sleep 2

echo "💸 Starting TurboRoute on port 8081..."
cd services/turboroute/backend
./turboroute &
TURBOROUTE_PID=$!
cd ../../..

sleep 2

echo ""
echo "✅ All services started!"
echo ""
echo "📍 Service URLs:"
echo "   TurboAuth:  http://localhost:8080"
echo "   TurboRoute: http://localhost:8081"
echo "   Prometheus: http://localhost:9093"
echo "   Grafana:    http://localhost:3000"
echo ""
echo "🧪 Test with:"
echo "   curl http://localhost:8080/health"
echo "   curl http://localhost:8081/health"
echo ""
echo "🛑 Stop with: Ctrl+C or kill $TURBOAUTH_PID $TURBOROUTE_PID"
echo ""

# Wait for interrupt
trap "echo ''; echo '🛑 Stopping services...'; kill $TURBOAUTH_PID $TURBOROUTE_PID 2>/dev/null; echo '✅ Services stopped'; exit 0" INT TERM

# Keep script running
wait
