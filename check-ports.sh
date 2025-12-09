#!/bin/bash

# Check which ports are used by start.sh services

echo "╔══════════════════════════════════════════════════════════════╗"
echo "║         Realtime Platform - Port Usage Check                ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""

# Check if services are running
echo "🔍 Checking running processes..."
TURBOAUTH_RUNNING=$(ps aux | grep -E "[t]urboauth" | grep -v grep)
TURBOROUTE_RUNNING=$(ps aux | grep -E "[t]urboroute" | grep -v grep)

if [ -z "$TURBOAUTH_RUNNING" ] && [ -z "$TURBOROUTE_RUNNING" ]; then
    echo "❌ No services running"
    echo ""
    echo "Start services with: ./start.sh"
    exit 0
fi

echo ""
echo "📊 Running Processes:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
ps aux | grep -E "(turboauth|turboroute)" | grep -v grep | awk '{printf "%-12s PID: %-6s CPU: %-5s MEM: %-5s\n", $11, $2, $3"%", $4"%"}'

echo ""
echo "🔌 Ports in Use:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Function to check a specific port
check_port() {
    local port=$1
    local service=$2
    local result=$(lsof -i :$port -P -n 2>/dev/null | grep LISTEN)
    
    if [ -n "$result" ]; then
        local process=$(echo "$result" | awk '{print $1}')
        local pid=$(echo "$result" | awk '{print $2}')
        echo "✅ Port $port ($service): $process (PID: $pid)"
    else
        echo "⚪ Port $port ($service): Free"
    fi
}

# Check all expected ports
check_port 8080 "TurboAuth HTTP"
check_port 9090 "TurboAuth gRPC"
check_port 2112 "TurboAuth Metrics"
check_port 8081 "TurboRoute HTTP"
check_port 9091 "TurboRoute gRPC"
check_port 2113 "TurboRoute Metrics"
check_port 6379 "Redis"
check_port 9093 "Prometheus"
check_port 3000 "Grafana"

echo ""
echo "📋 Detailed Port Information:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Show detailed info for turboauth/turboroute
lsof -i -P -n | grep -E "(turboauth|turboroute)" | grep LISTEN | while read line; do
    process=$(echo "$line" | awk '{print $1}')
    pid=$(echo "$line" | awk '{print $2}')
    port=$(echo "$line" | awk '{print $9}' | cut -d: -f2)
    echo "Process: $process | PID: $pid | Port: $port"
done

echo ""
echo "🌐 Service URLs:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if [ -n "$TURBOAUTH_RUNNING" ]; then
    echo "TurboAuth:  http://localhost:8080 (HTTP)"
    echo "            localhost:9090 (gRPC)"
    echo "            http://localhost:2112/metrics (Metrics)"
fi

if [ -n "$TURBOROUTE_RUNNING" ]; then
    echo "TurboRoute: http://localhost:8081 (HTTP)"
    echo "            localhost:9091 (gRPC)"
    echo "            http://localhost:2113/metrics (Metrics)"
fi

echo ""
echo "🧪 Test Services:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "curl http://localhost:8080/health  # TurboAuth"
echo "curl http://localhost:8081/health  # TurboRoute"

echo ""
echo "🛑 Stop Services:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if [ -n "$TURBOAUTH_RUNNING" ]; then
    TURBOAUTH_PID=$(echo "$TURBOAUTH_RUNNING" | awk '{print $2}')
    echo "kill $TURBOAUTH_PID  # Stop TurboAuth"
fi
if [ -n "$TURBOROUTE_RUNNING" ]; then
    TURBOROUTE_PID=$(echo "$TURBOROUTE_RUNNING" | awk '{print $2}')
    echo "kill $TURBOROUTE_PID  # Stop TurboRoute"
fi
echo "# Or use: make kill"

echo ""
