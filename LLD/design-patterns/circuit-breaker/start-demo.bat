@echo off
REM Circuit Breaker Demo Startup Script for Windows
REM This script starts all services and runs a demo

echo 🚀 Starting Circuit Breaker Demo...
echo ====================================

REM Function to check if a service is running
:check_service
set "url=%1"
set "name=%2"
set "max_attempts=30"
set "attempt=1"

echo ⏳ Waiting for %name% to start...

:check_loop
curl -s -f "%url%" >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ %name% is ready!
    exit /b 0
)

echo    Attempt %attempt%/%max_attempts%...
timeout /t 2 /nobreak >nul
set /a attempt+=1

if %attempt% leq %max_attempts% goto check_loop

echo ❌ %name% failed to start after %max_attempts% attempts
exit /b 1

:start_services
echo 📦 Building binaries...
make build

if not exist "bin\" (
    echo ❌ Build failed
    exit /b 1
)

echo 🚀 Starting all services...

echo 🔄 Starting Market Data Service...
start /b bin\market-data-service.exe

echo 🔄 Starting Portfolio Service...
start /b bin\portfolio-service.exe

echo 🔄 Starting Trading Gateway...
start /b bin\trading-gateway.exe

echo.
echo ⏳ Waiting for all services to be ready...

call :check_service "http://localhost:8082/api/v1/health" "Market Data Service"
if %errorlevel% neq 0 exit /b 1

call :check_service "http://localhost:8081/api/v1/health" "Portfolio Service"
if %errorlevel% neq 0 exit /b 1

call :check_service "http://localhost:8080/api/v1/health" "Trading Gateway"
if %errorlevel% neq 0 exit /b 1

echo.
echo 🎉 All services are ready!
echo.
echo 📊 Available Endpoints:
echo    Trading Gateway:    http://localhost:8080
echo    Portfolio Service:  http://localhost:8081
echo    Market Data:        http://localhost:8082
echo    Metrics:            http://localhost:8080/metrics
echo.
echo 🧪 Circuit Breaker Status: http://localhost:8080/api/v1/circuit-breaker/status
echo.
echo 💡 Example API calls:
echo.
echo    # Execute a trade
echo    curl -X POST http://localhost:8080/api/v1/trades ^
echo      -H "Content-Type: application/json" ^
echo      -d "{\"userId\":\"user123\",\"symbol\":\"AAPL\",\"quantity\":10,\"orderType\":\"BUY\",\"price\":150.00}"
echo.
echo    # Get portfolio
echo    curl http://localhost:8080/api/v1/portfolio/user123
echo.
echo    # Get market data
echo    curl http://localhost:8080/api/v1/market-data/AAPL
echo.
echo    # Simulate market data failures
echo    curl -X POST http://localhost:8082/api/v1/simulate/failure ^
echo      -H "Content-Type: application/json" ^
echo      -d "{\"failure_rate\":0.7,\"is_healthy\":true}"
echo.

REM Ask user if they want to run the demo
set /p response="🎬 Would you like to run the circuit breaker demo? (y/n): "

if /i "%response%"=="y" (
    echo.
    echo 🎬 Running Circuit Breaker Demo...
    echo ==================================
    go run scripts\demo.go
) else (
    echo.
    echo ✅ Services are running. Use Postman or curl to test the APIs.
    echo 📖 See docs\POSTMAN_GUIDE.md for detailed API examples.
    echo.
    echo Press Ctrl+C to stop all services...
    pause
)

call :start_services
