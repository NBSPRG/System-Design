#!/bin/bash

# Circuit Breaker Demo Startup Script
# This script starts all services and runs a demo

set -e

echo "🚀 Starting Circuit Breaker Demo..."
echo "===================================="

# Function to check if a service is running
check_service() {
    local url=$1
    local name=$2
    local max_attempts=30
    local attempt=1
    
    echo "⏳ Waiting for $name to start..."
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s -f "$url" >/dev/null 2>&1; then
            echo "✅ $name is ready!"
            return 0
        fi
        
        echo "   Attempt $attempt/$max_attempts..."
        sleep 2
        attempt=$((attempt + 1))
    done
    
    echo "❌ $name failed to start after $max_attempts attempts"
    return 1
}

# Function to start a service in background
start_service() {
    local binary=$1
    local name=$2
    local port=$3
    
    echo "🔄 Starting $name..."
    if [ -f "./bin/$binary" ]; then
        ./bin/$binary &
        echo $! > ".$binary.pid"
        echo "   $name started with PID $(cat .$binary.pid)"
    else
        echo "❌ Binary ./bin/$binary not found. Please run 'make build' first."
        exit 1
    fi
}

# Function to stop all services
stop_services() {
    echo "🛑 Stopping all services..."
    
    for pidfile in .*.pid; do
        if [ -f "$pidfile" ]; then
            pid=$(cat "$pidfile")
            service_name=$(basename "$pidfile" .pid)
            
            if kill -0 "$pid" 2>/dev/null; then
                echo "   Stopping $service_name (PID: $pid)..."
                kill "$pid"
            fi
            
            rm -f "$pidfile"
        fi
    done
    
    echo "✅ All services stopped"
}

# Trap to stop services on script exit
trap stop_services EXIT INT TERM

# Check if binaries exist
if [ ! -d "./bin" ]; then
    echo "📦 Building binaries..."
    make build
fi

# Start all services
echo "🚀 Starting all services..."
start_service "market-data-service" "Market Data Service" 8082
start_service "portfolio-service" "Portfolio Service" 8081
start_service "trading-gateway" "Trading Gateway" 8080

echo ""
echo "⏳ Waiting for all services to be ready..."

# Wait for services to be ready
check_service "http://localhost:8082/api/v1/health" "Market Data Service"
check_service "http://localhost:8081/api/v1/health" "Portfolio Service"  
check_service "http://localhost:8080/api/v1/health" "Trading Gateway"

echo ""
echo "🎉 All services are ready!"
echo ""
echo "📊 Available Endpoints:"
echo "   Trading Gateway:    http://localhost:8080"
echo "   Portfolio Service:  http://localhost:8081" 
echo "   Market Data:        http://localhost:8082"
echo "   Metrics:            http://localhost:8080/metrics"
echo ""
echo "🧪 Circuit Breaker Status: http://localhost:8080/api/v1/circuit-breaker/status"
echo ""
echo "💡 Example API calls:"
echo ""
echo "   # Execute a trade"
echo "   curl -X POST http://localhost:8080/api/v1/trades \\"
echo "     -H 'Content-Type: application/json' \\"
echo "     -d '{\"userId\":\"user123\",\"symbol\":\"AAPL\",\"quantity\":10,\"orderType\":\"BUY\",\"price\":150.00}'"
echo ""
echo "   # Get portfolio"
echo "   curl http://localhost:8080/api/v1/portfolio/user123"
echo ""
echo "   # Get market data"
echo "   curl http://localhost:8080/api/v1/market-data/AAPL"
echo ""
echo "   # Simulate market data failures"
echo "   curl -X POST http://localhost:8082/api/v1/simulate/failure \\"
echo "     -H 'Content-Type: application/json' \\"
echo "     -d '{\"failure_rate\":0.7,\"is_healthy\":true}'"
echo ""

# Ask user if they want to run the demo
echo "🎬 Would you like to run the circuit breaker demo? (y/n)"
read -r response

if [[ "$response" =~ ^[Yy]$ ]]; then
    echo ""
    echo "🎬 Running Circuit Breaker Demo..."
    echo "=================================="
    go run scripts/demo.go
else
    echo ""
    echo "✅ Services are running. Use Postman or curl to test the APIs."
    echo "📖 See docs/POSTMAN_GUIDE.md for detailed API examples."
    echo ""
    echo "Press Ctrl+C to stop all services..."
    
    # Keep script running
    while true; do
        sleep 1
    done
fi
