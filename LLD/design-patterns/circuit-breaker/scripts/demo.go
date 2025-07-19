package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// TradeRequest represents a trading request
type TradeRequest struct {
	UserID    string  `json:"userId"`
	Symbol    string  `json:"symbol"`
	Quantity  int64   `json:"quantity"`
	OrderType string  `json:"orderType"`
	Price     float64 `json:"price"`
}

// CircuitBreakerDemo demonstrates circuit breaker functionality
type CircuitBreakerDemo struct {
	gatewayURL    string
	marketDataURL string
	portfolioURL  string
	client        *http.Client
}

// NewDemo creates a new demo instance
func NewDemo() *CircuitBreakerDemo {
	return &CircuitBreakerDemo{
		gatewayURL:    "http://localhost:8080",
		marketDataURL: "http://localhost:8082",
		portfolioURL:  "http://localhost:8081",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Run executes the circuit breaker demonstration
func (d *CircuitBreakerDemo) Run() {
	fmt.Println("🔄 Circuit Breaker Demonstration Starting...")
	fmt.Println("==========================================")

	// Wait for services to start
	d.waitForServices()

	// Phase 1: Normal operation
	fmt.Println("\n📊 Phase 1: Normal Operation")
	fmt.Println("-----------------------------")
	d.normalOperation()

	// Phase 2: Introduce failures
	fmt.Println("\n⚠️  Phase 2: Introducing Market Data Service Failures")
	fmt.Println("-----------------------------------------------------")
	d.introduceMarketDataFailures()

	// Phase 3: Circuit breaker opens
	fmt.Println("\n🚫 Phase 3: Circuit Breaker Opens - Fast Failures")
	fmt.Println("--------------------------------------------------")
	d.testCircuitBreakerOpen()

	// Phase 4: Recovery attempt
	fmt.Println("\n🔄 Phase 4: Recovery Attempt - Half-Open State")
	fmt.Println("----------------------------------------------")
	d.testRecovery()

	// Phase 5: Full recovery
	fmt.Println("\n✅ Phase 5: Full Recovery - Circuit Closed")
	fmt.Println("------------------------------------------")
	d.fullRecovery()

	fmt.Println("\n🎉 Demo completed! Check the metrics at http://localhost:8080/metrics")
}

// waitForServices waits for all services to be available
func (d *CircuitBreakerDemo) waitForServices() {
	fmt.Println("⏳ Waiting for services to start...")

	services := map[string]string{
		"Trading Gateway":   d.gatewayURL + "/api/v1/health",
		"Market Data":       d.marketDataURL + "/api/v1/health",
		"Portfolio Service": d.portfolioURL + "/api/v1/health",
	}

	for name, url := range services {
		for {
			resp, err := d.client.Get(url)
			if err == nil && resp.StatusCode == 200 {
				resp.Body.Close()
				fmt.Printf("✅ %s is ready\n", name)
				break
			}
			if resp != nil {
				resp.Body.Close()
			}
			fmt.Printf("⏳ Waiting for %s...\n", name)
			time.Sleep(2 * time.Second)
		}
	}

	time.Sleep(2 * time.Second) // Give services a moment to fully initialize
}

// normalOperation demonstrates normal trading operations
func (d *CircuitBreakerDemo) normalOperation() {
	trades := []TradeRequest{
		{"user123", "AAPL", 10, "BUY", 150.00},
		{"user456", "GOOGL", 5, "BUY", 2800.00},
		{"user789", "MSFT", 20, "BUY", 300.00},
		{"user123", "TSLA", 15, "BUY", 800.00},
	}

	fmt.Println("🔄 Executing normal trades...")

	for i, trade := range trades {
		result := d.executeTrade(trade)
		status := "✅ SUCCESS"
		if !result {
			status = "❌ FAILED"
		}
		fmt.Printf("Trade %d: %s %d %s @ $%.2f - %s\n",
			i+1, trade.OrderType, trade.Quantity, trade.Symbol, trade.Price, status)
		time.Sleep(500 * time.Millisecond)
	}

	// Show circuit breaker status
	d.showCircuitBreakerStatus()
}

// introduceMarketDataFailures simulates market data service failures
func (d *CircuitBreakerDemo) introduceMarketDataFailures() {
	// Configure market data service to fail 70% of requests
	failureConfig := map[string]interface{}{
		"failure_rate": 0.7,
		"is_healthy":   true,
	}

	d.configureServiceFailures(d.marketDataURL+"/api/v1/simulate/failure", failureConfig)
	fmt.Println("📉 Market Data Service configured to fail 70% of requests")

	// Continue trading to trigger failures
	trades := []TradeRequest{
		{"user123", "NVDA", 8, "BUY", 900.00},
		{"user456", "META", 12, "BUY", 350.00},
		{"user789", "AMZN", 3, "BUY", 3200.00},
		{"user123", "AAPL", 25, "SELL", 155.00},
		{"user456", "GOOGL", 2, "SELL", 2850.00},
	}

	fmt.Println("🔄 Executing trades while market data service is failing...")

	successCount := 0
	for i, trade := range trades {
		result := d.executeTrade(trade)
		status := "✅ SUCCESS"
		if result {
			successCount++
		} else {
			status = "❌ FAILED"
		}
		fmt.Printf("Trade %d: %s %d %s @ $%.2f - %s\n",
			i+1, trade.OrderType, trade.Quantity, trade.Symbol, trade.Price, status)
		time.Sleep(800 * time.Millisecond)
	}

	fmt.Printf("📊 Success rate: %d/%d (%.1f%%)\n", successCount, len(trades), float64(successCount)/float64(len(trades))*100)
	d.showCircuitBreakerStatus()
}

// testCircuitBreakerOpen tests behavior when circuit breaker is open
func (d *CircuitBreakerDemo) testCircuitBreakerOpen() {
	// Make the market data service completely unavailable
	failureConfig := map[string]interface{}{
		"failure_rate": 1.0,
		"is_healthy":   false,
	}

	d.configureServiceFailures(d.marketDataURL+"/api/v1/simulate/failure", failureConfig)
	fmt.Println("💥 Market Data Service is now completely unavailable")

	// Wait for circuit breaker to open
	time.Sleep(3 * time.Second)

	// Try rapid trades - should fail fast
	fmt.Println("🔄 Attempting rapid trades (should fail fast)...")

	start := time.Now()
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(tradeNum int) {
			defer wg.Done()
			trade := TradeRequest{
				UserID:    fmt.Sprintf("user%d", tradeNum),
				Symbol:    "AAPL",
				Quantity:  1,
				OrderType: "BUY",
				Price:     150.00,
			}

			tradeStart := time.Now()
			result := d.executeTrade(trade)
			duration := time.Since(tradeStart)

			status := "❌ FAILED"
			if result {
				status = "✅ SUCCESS"
			}
			fmt.Printf("Rapid Trade %d: %s (%.2fms)\n", tradeNum+1, status, float64(duration.Nanoseconds())/1e6)
		}(i)

		time.Sleep(100 * time.Millisecond) // Small delay between requests
	}

	wg.Wait()
	totalDuration := time.Since(start)
	fmt.Printf("⚡ Total time for 10 trades: %.2fs (avg: %.2fms per trade)\n",
		totalDuration.Seconds(), float64(totalDuration.Nanoseconds())/1e6/10)

	d.showCircuitBreakerStatus()
}

// testRecovery tests the half-open state and recovery
func (d *CircuitBreakerDemo) testRecovery() {
	// Restore market data service partially
	failureConfig := map[string]interface{}{
		"failure_rate": 0.3,
		"is_healthy":   true,
	}

	d.configureServiceFailures(d.marketDataURL+"/api/v1/simulate/failure", failureConfig)
	fmt.Println("🛠️  Market Data Service partially restored (30% failure rate)")

	// Wait for circuit breaker to attempt recovery
	fmt.Println("⏳ Waiting for circuit breaker timeout...")
	time.Sleep(35 * time.Second) // Wait for timeout to trigger half-open state

	// Try limited trades in half-open state
	fmt.Println("🔄 Testing recovery with limited trades...")

	for i := 0; i < 5; i++ {
		trade := TradeRequest{
			UserID:    "recovery_user",
			Symbol:    "AAPL",
			Quantity:  1,
			OrderType: "BUY",
			Price:     150.00,
		}

		result := d.executeTrade(trade)
		status := "✅ SUCCESS"
		if !result {
			status = "❌ FAILED"
		}
		fmt.Printf("Recovery Trade %d: %s\n", i+1, status)
		time.Sleep(2 * time.Second)

		d.showCircuitBreakerStatus()
	}
}

// fullRecovery demonstrates full service recovery
func (d *CircuitBreakerDemo) fullRecovery() {
	// Fully restore market data service
	failureConfig := map[string]interface{}{
		"failure_rate": 0.0,
		"is_healthy":   true,
	}

	d.configureServiceFailures(d.marketDataURL+"/api/v1/simulate/failure", failureConfig)
	fmt.Println("🎉 Market Data Service fully restored!")

	// Execute successful trades to close circuit breaker
	trades := []TradeRequest{
		{"user123", "AAPL", 5, "BUY", 150.00},
		{"user456", "GOOGL", 2, "BUY", 2800.00},
		{"user789", "MSFT", 10, "BUY", 300.00},
	}

	fmt.Println("🔄 Executing trades to verify full recovery...")

	successCount := 0
	for i, trade := range trades {
		result := d.executeTrade(trade)
		status := "✅ SUCCESS"
		if result {
			successCount++
		} else {
			status = "❌ FAILED"
		}
		fmt.Printf("Recovery Trade %d: %s %d %s @ $%.2f - %s\n",
			i+1, trade.OrderType, trade.Quantity, trade.Symbol, trade.Price, status)
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("🎯 Final success rate: %d/%d (%.1f%%)\n", successCount, len(trades), float64(successCount)/float64(len(trades))*100)
	d.showCircuitBreakerStatus()
}

// executeTrade executes a single trade request
func (d *CircuitBreakerDemo) executeTrade(trade TradeRequest) bool {
	jsonData, _ := json.Marshal(trade)

	resp, err := d.client.Post(
		d.gatewayURL+"/api/v1/trades",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}

// configureServiceFailures configures failure simulation for a service
func (d *CircuitBreakerDemo) configureServiceFailures(url string, config map[string]interface{}) {
	jsonData, _ := json.Marshal(config)

	resp, err := d.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err == nil {
		resp.Body.Close()
	}
}

// showCircuitBreakerStatus displays current circuit breaker status
func (d *CircuitBreakerDemo) showCircuitBreakerStatus() {
	resp, err := d.client.Get(d.gatewayURL + "/api/v1/circuit-breaker/status")
	if err != nil {
		fmt.Printf("❌ Failed to get circuit breaker status: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Failed to read circuit breaker status: %v\n", err)
		return
	}

	var status map[string]interface{}
	if err := json.Unmarshal(body, &status); err != nil {
		fmt.Printf("❌ Failed to parse circuit breaker status: %v\n", err)
		return
	}

	fmt.Println("\n📊 Circuit Breaker Status:")
	fmt.Println("-------------------------")

	// Extract market data service status
	if marketData, ok := status["market_data_service"].(map[string]interface{}); ok {
		state := marketData["state"].(string)
		failures := marketData["failures"].(float64)
		requests := marketData["requests"].(float64)

		stateIcon := "🔴"
		if state == "CLOSED" {
			stateIcon = "🟢"
		} else if state == "HALF_OPEN" {
			stateIcon = "🟡"
		}

		fmt.Printf("Market Data Service: %s %s (Failures: %.0f, Requests: %.0f)\n",
			stateIcon, state, failures, requests)
	}

	// Extract portfolio service status
	if portfolio, ok := status["portfolio_service"].(map[string]interface{}); ok {
		state := portfolio["state"].(string)
		failures := portfolio["failures"].(float64)
		requests := portfolio["requests"].(float64)

		stateIcon := "🔴"
		if state == "CLOSED" {
			stateIcon = "🟢"
		} else if state == "HALF_OPEN" {
			stateIcon = "🟡"
		}

		fmt.Printf("Portfolio Service:   %s %s (Failures: %.0f, Requests: %.0f)\n",
			stateIcon, state, failures, requests)
	}

	fmt.Println()
}

func main() {
	demo := NewDemo()
	demo.Run()
}
