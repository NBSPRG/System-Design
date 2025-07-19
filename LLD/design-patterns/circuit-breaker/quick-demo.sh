#!/bin/bash

echo "🎯 Circuit Breaker Demo - Financial Trading Platform"
echo "=================================================="
echo

# 1. Show normal operation
echo "�� 1. Normal Operation - Executing successful trades"
echo "---------------------------------------------------"
for i in {1..3}; do
    echo "Trade $i:"
    curl -s -X POST http://localhost:8080/api/v1/trades \
        -H "Content-Type: application/json" \
        -d "{\"userId\":\"demo_user\",\"symbol\":\"AAPL\",\"quantity\":$i,\"orderType\":\"BUY\",\"price\":150.00}" \
        | grep -o '"status":"[^"]*"' | sed 's/"status":/Status: /'
    sleep 1
done
echo

# 2. Check circuit breaker status
echo "🔍 2. Circuit Breaker Status"
echo "----------------------------"
curl -s http://localhost:8080/api/v1/circuit-breaker/status | \
    grep -o '"market_data_service":{[^}]*}' | \
    sed 's/.*"state":"\([^"]*\)".*/Market Data Service Circuit Breaker: \1/'
echo

# 3. Introduce failures
echo "⚠️  3. Introducing Market Data Service Failures (70% failure rate)"
echo "-----------------------------------------------------------------"
curl -s -X POST http://localhost:8082/api/v1/simulate/failure \
    -H "Content-Type: application/json" \
    -d '{"failure_rate": 0.7, "is_healthy": true}' > /dev/null
echo "Market data service configured to fail 70% of requests"
echo

# 4. Execute trades during failures
echo "🔄 4. Executing trades during service failures"
echo "----------------------------------------------"
for i in {1..5}; do
    echo "Trade $i:"
    curl -s -X POST http://localhost:8080/api/v1/trades \
        -H "Content-Type: application/json" \
        -d "{\"userId\":\"demo_user\",\"symbol\":\"TSLA\",\"quantity\":1,\"orderType\":\"BUY\",\"price\":800.00}" \
        | grep -o '"status":"[^"]*"' | sed 's/"status":/Status: /' || echo "Status: FAILED"
    sleep 1
done
echo

# 5. Make service completely unavailable
echo "💥 5. Making Market Data Service Completely Unavailable"
echo "------------------------------------------------------"
curl -s -X POST http://localhost:8082/api/v1/simulate/failure \
    -H "Content-Type: application/json" \
    -d '{"failure_rate": 1.0, "is_healthy": false}' > /dev/null
echo "Market data service is now completely down"
echo

# 6. Test rapid failures
echo "⚡ 6. Testing Fast Failures (Circuit Breaker Protection)"
echo "-------------------------------------------------------"
echo "Measuring response time for trades when service is down:"
for i in {1..3}; do
    echo -n "Trade $i: "
    time_output=$(time (curl -s -X POST http://localhost:8080/api/v1/trades \
        -H "Content-Type: application/json" \
        -d "{\"userId\":\"demo_user\",\"symbol\":\"MSFT\",\"quantity\":1,\"orderType\":\"BUY\",\"price\":300.00}" \
        | grep -o '"status":"[^"]*"' | sed 's/"status":/Status: /') 2>&1)
    
    real_time=$(echo "$time_output" | grep real | awk '{print $2}')
    status=$(echo "$time_output" | grep "Status:" || echo "Status: FAILED")
    echo "$status (Time: $real_time)"
    sleep 1
done
echo

# 7. Restore service
echo "🛠️  7. Restoring Market Data Service"
echo "-----------------------------------"
curl -s -X POST http://localhost:8082/api/v1/simulate/failure \
    -H "Content-Type: application/json" \
    -d '{"failure_rate": 0.0, "is_healthy": true}' > /dev/null
echo "Market data service restored to normal operation"
echo

# 8. Final status check
echo "✅ 8. Final Status Check"
echo "------------------------"
echo "Services are now healthy and ready for normal operation"
curl -s http://localhost:8080/api/v1/health | \
    grep -o '"status":"[^"]*"' | sed 's/"status":/Trading Gateway Status: /'
echo

echo "🎉 Demo Complete!"
echo "=================="
echo "Key Observations:"
echo "• Circuit breaker protects against cascading failures"
echo "• Fast failure responses when services are down"
echo "• Automatic fallback mechanisms maintain service availability"
echo "• System recovers gracefully when dependencies are restored"
echo
echo "🔗 API Endpoints for further testing:"
echo "  Trading Gateway:     http://localhost:8080"
echo "  Circuit Breaker:     http://localhost:8080/api/v1/circuit-breaker/status"
echo "  Prometheus Metrics:  http://localhost:8080/metrics"
