# Distributed Circuit Breaker - Go Implementation

A production-ready circuit breaker implementation for distributed systems in Go.

## Real-World Use Case: Financial Trading Platform

This implementation simulates a real financial trading platform with multiple backend microservices:

- **Trading API Gateway** - Entry point for all trading requests with circuit breaker protection
- **Portfolio Service** - Manages user portfolios and positions
- **Market Data Service** - Fetches real-time market prices (simulates external dependency)
- **Risk Management Service** - Validates trades against risk limits (simulated)
- **Notification Service** - Sends trade confirmations and alerts (simulated)
- **Audit Service** - Logs all trading activities for compliance (simulated)

The circuit breaker protects against cascading failures when external market data providers or internal services become unavailable.

## Features

- **Production-Ready**: Thread-safe, high-performance implementation
- **Configurable States**: CLOSED, OPEN, HALF_OPEN with configurable thresholds
- **Multiple Strategies**: Failure rate, consecutive failures, and timeout-based
- **Metrics & Monitoring**: Detailed metrics with Prometheus integration
- **Graceful Degradation**: Fallback mechanisms for service failures
- **Health Checks**: Automatic service health monitoring
- **Request Isolation**: Separate circuit breakers for each service dependency

## Circuit Breaker States

1. **CLOSED**: Normal operation, all requests pass through
2. **OPEN**: Service is failing, requests fail fast with fallback
3. **HALF_OPEN**: Limited requests allowed to test recovery

## Quick Start

### Prerequisites
- Go 1.21 or later
- Git
- curl (for testing)

### 1. Clone and Setup

```bash
# Initialize Go module (if not already done)
go mod init distributed-circuit-breaker
go mod tidy
```

### 2. Easy Start (Recommended)

**Windows:**
```bash
# Make executable and run
chmod +x start-demo.bat
start-demo.bat
```

**Linux/macOS:**
```bash
# Make executable and run
chmod +x start-demo.sh
./start-demo.sh
```

### 3. Manual Start

```bash
# Build all services
make build

# Start all services
make run-all

# Or start individually
make run-market-data  # Port 8082
make run-portfolio    # Port 8081  
make run-gateway      # Port 8080

# Stop all services
make stop-all
```

### 4. Quick Demo

```bash
# Run the interactive circuit breaker demo
make demo
```

## Service Endpoints

### Trading Gateway (Port 8080)
- `POST /api/v1/trades` - Execute a trade
- `GET /api/v1/portfolio/{userId}` - Get user portfolio
- `GET /api/v1/market-data/{symbol}` - Get market data
- `GET /api/v1/health` - Health check
- `GET /api/v1/circuit-breaker/status` - Circuit breaker status
- `GET /metrics` - Prometheus metrics

### Portfolio Service (Port 8081)
- `GET /api/v1/portfolio/{userId}` - Get portfolio
- `POST /api/v1/portfolio/{userId}/positions` - Update position
- `GET /api/v1/health` - Health check
- `POST /api/v1/simulate/failure` - Simulate failures

### Market Data Service (Port 8082)
- `GET /api/v1/prices/{symbol}` - Get current price
- `POST /api/v1/prices/batch` - Get multiple prices
- `POST /api/v1/simulate/failure` - Simulate service failure
- `GET /api/v1/health` - Health check

## API Examples (Postman/curl)

### Execute Trade
```bash
curl -X POST http://localhost:8080/api/v1/trades \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "user123",
    "symbol": "AAPL", 
    "quantity": 10,
    "orderType": "BUY",
    "price": 150.00
  }'
```

### Get Portfolio
```bash
curl http://localhost:8080/api/v1/portfolio/user123
```

### Get Market Data
```bash
curl http://localhost:8080/api/v1/market-data/AAPL
```

### Check Circuit Breaker Status
```bash
curl http://localhost:8080/api/v1/circuit-breaker/status
```

### Simulate Market Data Failures
```bash
curl -X POST http://localhost:8082/api/v1/simulate/failure \
  -H "Content-Type: application/json" \
  -d '{
    "failure_rate": 0.7,
    "is_healthy": true,
    "response_time_ms": 2000
  }'
```

## Testing Circuit Breaker Behavior

### 1. Normal Operation
```bash
# Execute several successful trades
curl -X POST http://localhost:8080/api/v1/trades \
  -H "Content-Type: application/json" \
  -d '{"userId":"user123","symbol":"AAPL","quantity":10,"orderType":"BUY","price":150.00}'

# Check status - should show CLOSED state
curl http://localhost:8080/api/v1/circuit-breaker/status
```

### 2. Introduce Failures
```bash
# Make market data service fail 70% of requests
curl -X POST http://localhost:8082/api/v1/simulate/failure \
  -H "Content-Type: application/json" \
  -d '{"failure_rate": 0.7, "is_healthy": true}'

# Execute trades - some will fail, circuit breaker accumulates failures
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/v1/trades \
    -H "Content-Type: application/json" \
    -d "{\"userId\":\"user$i\",\"symbol\":\"AAPL\",\"quantity\":1,\"orderType\":\"BUY\",\"price\":150.00}"
  sleep 1
done
```

### 3. Trigger Circuit Breaker Open
```bash
# Make service completely unavailable
curl -X POST http://localhost:8082/api/v1/simulate/failure \
  -H "Content-Type: application/json" \
  -d '{"failure_rate": 1.0, "is_healthy": false}'

# Trades should now fail fast (circuit breaker OPEN)
curl -X POST http://localhost:8080/api/v1/trades \
  -H "Content-Type: application/json" \
  -d '{"userId":"user123","symbol":"AAPL","quantity":1,"orderType":"BUY","price":150.00}'
```

### 4. Test Recovery
```bash
# Restore service
curl -X POST http://localhost:8082/api/v1/simulate/failure \
  -H "Content-Type: application/json" \
  -d '{"failure_rate": 0.0, "is_healthy": true}'

# Wait 30+ seconds for circuit breaker timeout, then trade should work
sleep 35
curl -X POST http://localhost:8080/api/v1/trades \
  -H "Content-Type: application/json" \
  -d '{"userId":"user123","symbol":"AAPL","quantity":1,"orderType":"BUY","price":150.00}'
```

## Configuration

Circuit breakers can be configured in `config/config.yaml`:

```yaml
circuit_breaker:
  max_requests: 5          # Max requests in half-open state
  interval: 60s            # Statistical window
  timeout: 30s             # Time before attempting recovery
  failure_threshold: 10    # Failures to open circuit
  success_threshold: 3     # Successes to close circuit
  failure_rate_threshold: 0.6  # Failure rate (0.0-1.0)
  minimum_requests: 5      # Min requests before considering rate
```

## Monitoring & Metrics

### Prometheus Metrics
Access at: `http://localhost:8080/metrics`

Key metrics:
- `circuit_breaker_requests_total_*` - Total requests
- `circuit_breaker_failures_total_*` - Total failures  
- `circuit_breaker_state_changes_total_*` - State transitions
- `circuit_breaker_state_*` - Current state
- `circuit_breaker_request_duration_seconds_*` - Request latency

### Real-time Status
Check circuit breaker status: `http://localhost:8080/api/v1/circuit-breaker/status`

## Project Structure

```
distributed-circuit-breaker/
├── cmd/                          # Service entry points
│   ├── trading-gateway/         
│   ├── portfolio-service/       
│   └── market-data-service/     
├── pkg/                         # Shared packages
│   ├── circuitbreaker/         # Circuit breaker implementation
│   ├── httpclient/             # HTTP client with CB integration
│   ├── config/                 # Configuration management
│   └── models/                 # Data models
├── config/                     # Configuration files
├── docs/                       # Documentation
├── scripts/                    # Demo and testing scripts
├── Makefile                    # Build automation
└── README.md                   # This file
```

## Advanced Usage

### Custom Circuit Breaker
```go
cb := circuitbreaker.NewCircuitBreaker(
    circuitbreaker.Config{
        Name:                 "my-service",
        MaxRequests:          3,
        Timeout:              30 * time.Second,
        FailureThreshold:     5,
        SuccessThreshold:     2,
        FailureRateThreshold: 0.5,
        MinimumRequests:      3,
    },
    logger,
)

result, err := cb.Execute(ctx, func() (interface{}, error) {
    return myService.DoWork()
})
```

### HTTP Client Integration
```go
client := httpclient.NewHTTPClient(
    "http://api.example.com",
    5*time.Second,
    circuitBreaker,
    logger,
)

var response MyResponse
err := client.GetJSON(ctx, "/api/data", &response)
```

## Documentation

- [Postman API Testing Guide](docs/POSTMAN_GUIDE.md) - Detailed API examples
- [Circuit Breaker Configuration](config/config.yaml) - Configuration options
- [Architecture Overview](pkg/circuitbreaker/circuitbreaker.go) - Implementation details

## Common Use Cases

1. **External API Protection** - Protect against third-party API failures
2. **Database Connection Pooling** - Prevent connection exhaustion
3. **Microservice Communication** - Graceful degradation between services
4. **Rate Limiting** - Control request flow to downstream services
5. **Cache Fallbacks** - Fallback to cache when primary data source fails

## Troubleshooting

### Services Won't Start
```bash
# Check if ports are available
netstat -an | grep :8080
netstat -an | grep :8081  
netstat -an | grep :8082

# Check Go installation
go version

# Rebuild
make clean && make build
```

### Circuit Breaker Stuck Open
- Wait for timeout period (default 30s)
- Check underlying service health
- Reset failure simulation: `{"failure_rate": 0.0, "is_healthy": true}`

### Build Issues
```bash
# Clean and rebuild
go clean -modcache
go mod download
go mod tidy
make build
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

## License

This project is licensed under the MIT License.
