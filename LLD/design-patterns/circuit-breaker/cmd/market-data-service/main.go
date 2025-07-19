package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// MarketData represents market data for a symbol
type MarketData struct {
	Symbol        string    `json:"symbol"`
	Price         float64   `json:"price"`
	High          float64   `json:"high"`
	Low           float64   `json:"low"`
	Volume        int64     `json:"volume"`
	Change        float64   `json:"change"`
	ChangePercent float64   `json:"changePercent"`
	Timestamp     time.Time `json:"timestamp"`
}

// BatchMarketDataRequest represents a batch request
type BatchMarketDataRequest struct {
	Symbols []string `json:"symbols"`
}

// BatchMarketDataResponse represents a batch response
type BatchMarketDataResponse struct {
	Data      map[string]MarketData `json:"data"`
	Errors    map[string]string     `json:"errors,omitempty"`
	Timestamp time.Time             `json:"timestamp"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Service   string            `json:"service"`
	Version   string            `json:"version"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    map[string]string `json:"checks,omitempty"`
}

// MarketDataService simulates an external market data provider
type MarketDataService struct {
	prices       map[string]float64
	mutex        sync.RWMutex
	isHealthy    bool
	failureRate  float64
	responseTime time.Duration
}

// NewMarketDataService creates a new market data service
func NewMarketDataService() *MarketDataService {
	service := &MarketDataService{
		prices:       make(map[string]float64),
		isHealthy:    true,
		failureRate:  0.0,
		responseTime: 100 * time.Millisecond,
	}

	// Initialize with some sample stock prices
	symbols := []string{"AAPL", "GOOGL", "MSFT", "TSLA", "AMZN", "META", "NVDA", "BRK.A", "JPM", "JNJ"}
	basePrices := []float64{150.0, 2800.0, 300.0, 800.0, 3200.0, 350.0, 900.0, 450000.0, 150.0, 160.0}

	for i, symbol := range symbols {
		service.prices[symbol] = basePrices[i]
	}

	// Start price update goroutine
	go service.updatePrices()

	return service
}

// updatePrices simulates real-time price updates
func (s *MarketDataService) updatePrices() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.mutex.Lock()
		for symbol, currentPrice := range s.prices {
			// Random price change between -2% and +2%
			change := (rand.Float64() - 0.5) * 0.04
			newPrice := currentPrice * (1 + change)
			s.prices[symbol] = newPrice
		}
		s.mutex.Unlock()
	}
}

// simulateLatency adds artificial latency to simulate network delays
func (s *MarketDataService) simulateLatency() {
	s.mutex.RLock()
	latency := s.responseTime
	s.mutex.RUnlock()

	// Add random jitter
	jitter := time.Duration(rand.Float64() * float64(latency) * 0.5)
	time.Sleep(latency + jitter)
}

// simulateFailure randomly fails requests based on failure rate
func (s *MarketDataService) simulateFailure() bool {
	s.mutex.RLock()
	rate := s.failureRate
	healthy := s.isHealthy
	s.mutex.RUnlock()

	if !healthy {
		return true
	}

	return rand.Float64() < rate
}

// GetPrice handles individual price requests
func (s *MarketDataService) GetPrice(c *gin.Context) {
	symbol := c.Param("symbol")

	// Simulate latency
	s.simulateLatency()

	// Simulate failures
	if s.simulateFailure() {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Market data service temporarily unavailable",
			"timestamp": time.Now(),
		})
		return
	}

	s.mutex.RLock()
	currentPrice, exists := s.prices[symbol]
	s.mutex.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error":     "Symbol not found",
			"symbol":    symbol,
			"timestamp": time.Now(),
		})
		return
	}

	// Calculate some mock data
	high := currentPrice * (1 + rand.Float64()*0.05)
	low := currentPrice * (1 - rand.Float64()*0.05)
	volume := int64(rand.Intn(1000000) + 100000)
	change := (rand.Float64() - 0.5) * 10
	changePercent := change / currentPrice * 100

	marketData := MarketData{
		Symbol:        symbol,
		Price:         currentPrice,
		High:          high,
		Low:           low,
		Volume:        volume,
		Change:        change,
		ChangePercent: changePercent,
		Timestamp:     time.Now(),
	}

	c.JSON(http.StatusOK, marketData)
}

// GetBatchPrices handles batch price requests
func (s *MarketDataService) GetBatchPrices(c *gin.Context) {
	var request BatchMarketDataRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request format",
			"details":   err.Error(),
			"timestamp": time.Now(),
		})
		return
	}

	// Simulate latency
	s.simulateLatency()

	// Simulate failures
	if s.simulateFailure() {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Market data service temporarily unavailable",
			"timestamp": time.Now(),
		})
		return
	}

	response := BatchMarketDataResponse{
		Data:      make(map[string]MarketData),
		Errors:    make(map[string]string),
		Timestamp: time.Now(),
	}

	s.mutex.RLock()
	for _, symbol := range request.Symbols {
		if currentPrice, exists := s.prices[symbol]; exists {
			high := currentPrice * (1 + rand.Float64()*0.05)
			low := currentPrice * (1 - rand.Float64()*0.05)
			volume := int64(rand.Intn(1000000) + 100000)
			change := (rand.Float64() - 0.5) * 10
			changePercent := change / currentPrice * 100

			response.Data[symbol] = MarketData{
				Symbol:        symbol,
				Price:         currentPrice,
				High:          high,
				Low:           low,
				Volume:        volume,
				Change:        change,
				ChangePercent: changePercent,
				Timestamp:     time.Now(),
			}
		} else {
			response.Errors[symbol] = "Symbol not found"
		}
	}
	s.mutex.RUnlock()

	c.JSON(http.StatusOK, response)
}

// Health returns the health status of the service
func (s *MarketDataService) Health(c *gin.Context) {
	s.mutex.RLock()
	healthy := s.isHealthy
	s.mutex.RUnlock()

	status := "healthy"
	statusCode := http.StatusOK

	if !healthy {
		status = "unhealthy"
		statusCode = http.StatusServiceUnavailable
	}

	response := HealthResponse{
		Status:    status,
		Service:   "market-data-service",
		Version:   "1.0.0",
		Timestamp: time.Now(),
		Checks: map[string]string{
			"database": "healthy",
			"cache":    "healthy",
			"external": "healthy",
		},
	}

	c.JSON(statusCode, response)
}

// SimulateFailure allows external control of failure simulation
func (s *MarketDataService) SimulateFailure(c *gin.Context) {
	var request struct {
		FailureRate  float64 `json:"failure_rate"`
		IsHealthy    bool    `json:"is_healthy"`
		ResponseTime int     `json:"response_time_ms"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	s.mutex.Lock()
	s.failureRate = request.FailureRate
	s.isHealthy = request.IsHealthy
	if request.ResponseTime > 0 {
		s.responseTime = time.Duration(request.ResponseTime) * time.Millisecond
	}
	s.mutex.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message":       "Failure simulation updated",
		"failure_rate":  request.FailureRate,
		"is_healthy":    request.IsHealthy,
		"response_time": request.ResponseTime,
		"timestamp":     time.Now(),
	})
}

// GetStatus returns current service configuration
func (s *MarketDataService) GetStatus(c *gin.Context) {
	s.mutex.RLock()
	status := map[string]interface{}{
		"service":       "market-data-service",
		"version":       "1.0.0",
		"is_healthy":    s.isHealthy,
		"failure_rate":  s.failureRate,
		"response_time": s.responseTime.Milliseconds(),
		"symbols_count": len(s.prices),
		"timestamp":     time.Now(),
	}
	s.mutex.RUnlock()

	c.JSON(http.StatusOK, status)
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	service := NewMarketDataService()

	// Configure Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// API routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/prices/:symbol", service.GetPrice)
		v1.POST("/prices/batch", service.GetBatchPrices)
		v1.GET("/health", service.Health)
		v1.POST("/simulate/failure", service.SimulateFailure)
		v1.GET("/status", service.GetStatus)
	}

	// Start server
	port := 8082
	fmt.Printf("🚀 Market Data Service starting on port %d\n", port)
	fmt.Printf("📊 Available endpoints:\n")
	fmt.Printf("   GET  /api/v1/prices/{symbol}     - Get price for symbol\n")
	fmt.Printf("   POST /api/v1/prices/batch        - Get prices for multiple symbols\n")
	fmt.Printf("   GET  /api/v1/health              - Health check\n")
	fmt.Printf("   POST /api/v1/simulate/failure    - Simulate failures\n")
	fmt.Printf("   GET  /api/v1/status              - Service status\n")
	fmt.Printf("\n💡 Example: curl http://localhost:%d/api/v1/prices/AAPL\n", port)
	fmt.Printf("💡 Simulate failure: curl -X POST http://localhost:%d/api/v1/simulate/failure -H 'Content-Type: application/json' -d '{\"failure_rate\": 0.5, \"is_healthy\": true}'\n", port)

	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
