# Postman API Testing Guide

This guide provides examples of how to test the circuit breaker implementation using Postman.

## Service Endpoints

### Trading Gateway (Port 8080)
- Base URL: `http://localhost:8080`

### Market Data Service (Port 8082)
- Base URL: `http://localhost:8082`

### Portfolio Service (Port 8081)
- Base URL: `http://localhost:8081`

## API Examples

### 1. Execute a Trade

**POST** `http://localhost:8080/api/v1/trades`

```json
{
  "userId": "user123",
  "symbol": "AAPL",
  "quantity": 10,
  "orderType": "BUY",
  "price": 150.00
}
```

**Expected Response:**
```json
{
  "tradeId": "TXN_1640995200_user123",
  "userId": "user123",
  "symbol": "AAPL",
  "quantity": 10,
  "orderType": "BUY",
  "price": 150.25,
  "status": "EXECUTED",
  "message": "Trade executed successfully",
  "executedAt": "2024-01-01T10:00:00Z",
  "totalValue": 1502.50
}
```

### 2. Get User Portfolio

**GET** `http://localhost:8080/api/v1/portfolio/user123`

**Expected Response:**
```json
{
  "userId": "user123",
  "positions": [
    {
      "symbol": "AAPL",
      "quantity": 10,
      "averagePrice": 150.00,
      "currentPrice": 150.25,
      "marketValue": 1502.50,
      "pnl": 2.50,
      "pnlPercent": 0.17,
      "updatedAt": "2024-01-01T10:00:00Z"
    }
  ],
  "cashBalance": 8497.50,
  "totalValue": 10000.00,
  "updatedAt": "2024-01-01T10:00:00Z"
}
```

### 3. Get Market Data

**GET** `http://localhost:8080/api/v1/market-data/AAPL`

**Expected Response:**
```json
{
  "symbol": "AAPL",
  "price": 150.25,
  "high": 152.00,
  "low": 148.50,
  "volume": 1500000,
  "change": 2.25,
  "changePercent": 1.52,
  "timestamp": "2024-01-01T10:00:00Z"
}
```

### 4. Check Circuit Breaker Status

**GET** `http://localhost:8080/api/v1/circuit-breaker/status`

**Expected Response:**
```json
{
  "market_data_service": {
    "name": "market-data-service",
    "state": "CLOSED",
    "failures": 0,
    "requests": 5,
    "successes": 0,
    "halfOpenRequests": 0,
    "lastStateChange": "2024-01-01T09:00:00Z",
    "lastFailureTime": "0001-01-01T00:00:00Z"
  },
  "portfolio_service": {
    "name": "portfolio-service",
    "state": "CLOSED",
    "failures": 0,
    "requests": 3,
    "successes": 0,
    "halfOpenRequests": 0,
    "lastStateChange": "2024-01-01T09:00:00Z",
    "lastFailureTime": "0001-01-01T00:00:00Z"
  },
  "timestamp": "2024-01-01T10:00:00Z"
}
```

### 5. Health Check

**GET** `http://localhost:8080/api/v1/health`

**Expected Response:**
```json
{
  "status": "healthy",
  "service": "trading-gateway",
  "version": "1.0.0",
  "timestamp": "2024-01-01T10:00:00Z",
  "checks": {
    "market_data_circuit_breaker": "CLOSED",
    "portfolio_circuit_breaker": "CLOSED",
    "risk_management_circuit_breaker": "CLOSED",
    "notification_circuit_breaker": "CLOSED",
    "audit_circuit_breaker": "CLOSED"
  }
}
```

## Testing Circuit Breaker Behavior

### Step 1: Normal Operation
1. Execute several trades using the execute trade endpoint
2. Check circuit breaker status - should show CLOSED state
3. Verify trades are successful

### Step 2: Simulate Market Data Service Failures

**POST** `http://localhost:8082/api/v1/simulate/failure`

```json
{
  "failure_rate": 0.7,
  "is_healthy": true,
  "response_time_ms": 2000
}
```

This configures the market data service to fail 70% of requests with 2-second delays.

### Step 3: Execute Trades During Failures
1. Execute multiple trades
2. Observe some trades failing
3. Check circuit breaker status - failures should increase

### Step 4: Trigger Circuit Breaker to Open

**POST** `http://localhost:8082/api/v1/simulate/failure`

```json
{
  "failure_rate": 1.0,
  "is_healthy": false
}
```

This makes the market data service completely unavailable.

### Step 5: Test Fast Failures
1. Execute trades immediately after service becomes unavailable
2. Trades should fail quickly (circuit breaker is OPEN)
3. Check circuit breaker status - should show OPEN state

### Step 6: Test Recovery

**POST** `http://localhost:8082/api/v1/simulate/failure`

```json
{
  "failure_rate": 0.0,
  "is_healthy": true
}
```

1. Wait 30+ seconds for circuit breaker timeout
2. Execute a few trades
3. Circuit breaker should transition to HALF_OPEN, then CLOSED
4. Verify trades are successful again

## Testing Different Scenarios

### High-Frequency Trading Test
Execute rapid trades to test circuit breaker under load:

```bash
# Use a tool like Apache Bench or create a simple script
for i in {1..20}; do
  curl -X POST http://localhost:8080/api/v1/trades \
    -H "Content-Type: application/json" \
    -d "{\"userId\":\"user$i\",\"symbol\":\"AAPL\",\"quantity\":1,\"orderType\":\"BUY\",\"price\":150.00}" &
done
wait
```

### Mixed Symbol Trading
Test with different symbols:

```json
// Trade 1
{
  "userId": "user123",
  "symbol": "AAPL",
  "quantity": 10,
  "orderType": "BUY",
  "price": 150.00
}

// Trade 2  
{
  "userId": "user123",
  "symbol": "GOOGL",
  "quantity": 5,
  "orderType": "BUY",
  "price": 2800.00
}

// Trade 3
{
  "userId": "user123",
  "symbol": "TSLA",
  "quantity": 8,
  "orderType": "SELL",
  "price": 800.00
}
```

### Error Cases Testing

**Invalid Trade Request:**
```json
{
  "userId": "",
  "symbol": "INVALID",
  "quantity": -5,
  "orderType": "INVALID",
  "price": -100.00
}
```

**Insufficient Funds:**
```json
{
  "userId": "user123",
  "symbol": "BRK.A",
  "quantity": 1000,
  "orderType": "BUY",
  "price": 450000.00
}
```

## Monitoring and Metrics

### Prometheus Metrics
Access metrics at: `http://localhost:8080/metrics`

Key metrics to monitor:
- `circuit_breaker_requests_total_*`
- `circuit_breaker_failures_total_*`
- `circuit_breaker_state_changes_total_*`
- `circuit_breaker_state_*`

### Real-time Monitoring
Use tools like Grafana to visualize:
- Request success/failure rates
- Circuit breaker state changes
- Response times
- Service health status

## Troubleshooting

### Common Issues

1. **Service Not Responding**
   - Check if all services are running on correct ports
   - Verify network connectivity
   - Check service logs

2. **Circuit Breaker Stuck Open**
   - Wait for timeout period (default 30s)
   - Check if underlying service is healthy
   - Restart services if needed

3. **Trades Always Failing**
   - Check portfolio service health
   - Verify user has sufficient cash balance
   - Check market data service availability

### Debugging Steps

1. Check service health endpoints
2. Monitor circuit breaker status
3. Review service logs
4. Test individual service endpoints
5. Verify configuration settings
