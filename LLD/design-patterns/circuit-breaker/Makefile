.PHONY: build run-all run-market-data run-portfolio run-gateway clean test deps

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
MARKET_DATA_BINARY=market-data-service
PORTFOLIO_BINARY=portfolio-service
GATEWAY_BINARY=trading-gateway

# Build directory
BUILD_DIR=./bin

all: deps build

deps:
	$(GOMOD) download
	$(GOMOD) tidy

build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(MARKET_DATA_BINARY) ./cmd/market-data-service
	$(GOBUILD) -o $(BUILD_DIR)/$(PORTFOLIO_BINARY) ./cmd/portfolio-service
	$(GOBUILD) -o $(BUILD_DIR)/$(GATEWAY_BINARY) ./cmd/trading-gateway

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

test:
	$(GOTEST) -v ./...

# Run all services in background
run-all: build
	@echo "Starting all services..."
	./$(BUILD_DIR)/$(MARKET_DATA_BINARY) & echo $$! > .market-data.pid
	sleep 2
	./$(BUILD_DIR)/$(PORTFOLIO_BINARY) & echo $$! > .portfolio.pid
	sleep 2
	./$(BUILD_DIR)/$(GATEWAY_BINARY) & echo $$! > .gateway.pid
	@echo "All services started!"
	@echo "Trading Gateway: http://localhost:8080"
	@echo "Portfolio Service: http://localhost:8081"
	@echo "Market Data Service: http://localhost:8082"
	@echo "To stop all services, run: make stop-all"

# Stop all services
stop-all:
	@echo "Stopping all services..."
	-kill `cat .market-data.pid` 2>/dev/null || true
	-kill `cat .portfolio.pid` 2>/dev/null || true
	-kill `cat .gateway.pid` 2>/dev/null || true
	-rm -f .market-data.pid .portfolio.pid .gateway.pid
	@echo "All services stopped!"

# Run individual services
run-market-data: build
	./$(BUILD_DIR)/$(MARKET_DATA_BINARY)

run-portfolio: build
	./$(BUILD_DIR)/$(PORTFOLIO_BINARY)

run-gateway: build
	./$(BUILD_DIR)/$(GATEWAY_BINARY)

# Development mode (run with go run)
dev-market-data:
	$(GOCMD) run ./cmd/market-data-service/main.go

dev-portfolio:
	$(GOCMD) run ./cmd/portfolio-service/main.go

dev-gateway:
	$(GOCMD) run ./cmd/trading-gateway/main.go

# Docker targets
docker-build:
	docker build -t trading-gateway -f docker/Dockerfile.gateway .
	docker build -t market-data-service -f docker/Dockerfile.market-data .
	docker build -t portfolio-service -f docker/Dockerfile.portfolio .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

# Load testing
load-test:
	$(GOCMD) run ./scripts/load-test.go

# Circuit breaker demo
demo:
	$(GOCMD) run ./scripts/demo.go

help:
	@echo "Available targets:"
	@echo "  deps          - Download dependencies"
	@echo "  build         - Build all binaries"
	@echo "  run-all       - Start all services"
	@echo "  stop-all      - Stop all services"
	@echo "  run-gateway   - Run trading gateway"
	@echo "  run-portfolio - Run portfolio service"
	@echo "  run-market-data - Run market data service"
	@echo "  dev-*         - Run services in development mode"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  load-test     - Run load test"
	@echo "  demo          - Run circuit breaker demo"
	@echo "  help          - Show this help"
