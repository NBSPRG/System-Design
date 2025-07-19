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

// Position represents a position in a user's portfolio
type Position struct {
	Symbol       string    `json:"symbol"`
	Quantity     int64     `json:"quantity"`
	AveragePrice float64   `json:"averagePrice"`
	CurrentPrice float64   `json:"currentPrice"`
	MarketValue  float64   `json:"marketValue"`
	PnL          float64   `json:"pnl"`
	PnLPercent   float64   `json:"pnlPercent"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// Portfolio represents a user's portfolio
type Portfolio struct {
	UserID      string     `json:"userId"`
	Positions   []Position `json:"positions"`
	CashBalance float64    `json:"cashBalance"`
	TotalValue  float64    `json:"totalValue"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// UpdatePositionRequest represents a request to update a position
type UpdatePositionRequest struct {
	Symbol   string  `json:"symbol" binding:"required"`
	Quantity int64   `json:"quantity" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
	Action   string  `json:"action" binding:"required"` // "BUY" or "SELL"
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Service   string            `json:"service"`
	Version   string            `json:"version"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    map[string]string `json:"checks,omitempty"`
}

// PortfolioService manages user portfolios
type PortfolioService struct {
	portfolios  map[string]*Portfolio
	mutex       sync.RWMutex
	isHealthy   bool
	failureRate float64
}

// NewPortfolioService creates a new portfolio service
func NewPortfolioService() *PortfolioService {
	service := &PortfolioService{
		portfolios:  make(map[string]*Portfolio),
		isHealthy:   true,
		failureRate: 0.0,
	}

	// Initialize with some sample portfolios
	service.initializeSamplePortfolios()

	return service
}

// initializeSamplePortfolios creates sample portfolios for testing
func (s *PortfolioService) initializeSamplePortfolios() {
	sampleUsers := []string{"user123", "user456", "user789"}

	for _, userID := range sampleUsers {
		portfolio := &Portfolio{
			UserID:      userID,
			Positions:   []Position{},
			CashBalance: 10000.0 + rand.Float64()*90000.0, // Random cash between 10k-100k
			UpdatedAt:   time.Now(),
		}
		// Add some random positions
		symbols := []string{"AAPL", "GOOGL", "MSFT", "TSLA"}
		for _, symbol := range symbols {
			if rand.Float64() > 0.5 { // 50% chance to have this position
				avgPrice := 100.0 + rand.Float64()*500.0
				quantity := int64(10 + rand.Intn(100))
				currentPrice := avgPrice * (0.8 + rand.Float64()*0.4) // ±20% from avg price

				position := Position{
					Symbol:       symbol,
					Quantity:     quantity,
					AveragePrice: avgPrice,
					CurrentPrice: currentPrice,
					MarketValue:  currentPrice * float64(quantity),
					PnL:          (currentPrice - avgPrice) * float64(quantity),
					UpdatedAt:    time.Now(),
				}
				position.PnLPercent = (position.PnL / (avgPrice * float64(quantity))) * 100

				portfolio.Positions = append(portfolio.Positions, position)
			}
		}

		// Calculate total value
		totalValue := portfolio.CashBalance
		for _, pos := range portfolio.Positions {
			totalValue += pos.MarketValue
		}
		portfolio.TotalValue = totalValue

		s.portfolios[userID] = portfolio
	}
}

// simulateFailure randomly fails requests based on failure rate
func (s *PortfolioService) simulateFailure() bool {
	s.mutex.RLock()
	rate := s.failureRate
	healthy := s.isHealthy
	s.mutex.RUnlock()

	if !healthy {
		return true
	}

	return rand.Float64() < rate
}

// GetPortfolio returns a user's portfolio
func (s *PortfolioService) GetPortfolio(c *gin.Context) {
	userID := c.Param("userId")

	// Simulate failures
	if s.simulateFailure() {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Portfolio service temporarily unavailable",
			"timestamp": time.Now(),
		})
		return
	}

	s.mutex.RLock()
	portfolio, exists := s.portfolios[userID]
	s.mutex.RUnlock()

	if !exists {
		// Create a new portfolio for new users
		s.mutex.Lock()
		portfolio = &Portfolio{
			UserID:      userID,
			Positions:   []Position{},
			CashBalance: 10000.0, // Default starting cash
			TotalValue:  10000.0,
			UpdatedAt:   time.Now(),
		}
		s.portfolios[userID] = portfolio
		s.mutex.Unlock()
	}

	c.JSON(http.StatusOK, portfolio)
}

// UpdatePosition updates a position in a user's portfolio
func (s *PortfolioService) UpdatePosition(c *gin.Context) {
	userID := c.Param("userId")

	var request UpdatePositionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Simulate failures
	if s.simulateFailure() {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Portfolio service temporarily unavailable",
			"timestamp": time.Now(),
		})
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	portfolio, exists := s.portfolios[userID]
	if !exists {
		portfolio = &Portfolio{
			UserID:      userID,
			Positions:   []Position{},
			CashBalance: 10000.0,
			UpdatedAt:   time.Now(),
		}
		s.portfolios[userID] = portfolio
	}

	// Find existing position or create new one
	var targetPosition *Position
	for i := range portfolio.Positions {
		if portfolio.Positions[i].Symbol == request.Symbol {
			targetPosition = &portfolio.Positions[i]
			break
		}
	}

	totalCost := request.Price * float64(request.Quantity)

	if request.Action == "BUY" {
		// Check if user has enough cash
		if portfolio.CashBalance < totalCost {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "Insufficient cash balance",
				"required":  totalCost,
				"available": portfolio.CashBalance,
			})
			return
		}

		if targetPosition == nil {
			// Create new position
			position := Position{
				Symbol:       request.Symbol,
				Quantity:     request.Quantity,
				AveragePrice: request.Price,
				CurrentPrice: request.Price,
				MarketValue:  totalCost,
				PnL:          0,
				PnLPercent:   0,
				UpdatedAt:    time.Now(),
			}
			portfolio.Positions = append(portfolio.Positions, position)
		} else {
			// Update existing position
			totalQuantity := targetPosition.Quantity + request.Quantity
			totalValue := (targetPosition.AveragePrice * float64(targetPosition.Quantity)) + totalCost

			targetPosition.AveragePrice = totalValue / float64(totalQuantity)
			targetPosition.Quantity = totalQuantity
			targetPosition.CurrentPrice = request.Price
			targetPosition.MarketValue = targetPosition.CurrentPrice * float64(targetPosition.Quantity)
			targetPosition.PnL = (targetPosition.CurrentPrice - targetPosition.AveragePrice) * float64(targetPosition.Quantity)
			targetPosition.PnLPercent = (targetPosition.PnL / (targetPosition.AveragePrice * float64(targetPosition.Quantity))) * 100
			targetPosition.UpdatedAt = time.Now()
		}

		portfolio.CashBalance -= totalCost

	} else if request.Action == "SELL" {
		if targetPosition == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "No position found for symbol",
				"symbol": request.Symbol,
			})
			return
		}

		if targetPosition.Quantity < request.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "Insufficient shares to sell",
				"available": targetPosition.Quantity,
				"requested": request.Quantity,
			})
			return
		}

		// Update position
		targetPosition.Quantity -= request.Quantity
		targetPosition.CurrentPrice = request.Price

		if targetPosition.Quantity == 0 {
			// Remove position if quantity is zero
			for i, pos := range portfolio.Positions {
				if pos.Symbol == request.Symbol {
					portfolio.Positions = append(portfolio.Positions[:i], portfolio.Positions[i+1:]...)
					break
				}
			}
		} else {
			targetPosition.MarketValue = targetPosition.CurrentPrice * float64(targetPosition.Quantity)
			targetPosition.PnL = (targetPosition.CurrentPrice - targetPosition.AveragePrice) * float64(targetPosition.Quantity)
			targetPosition.PnLPercent = (targetPosition.PnL / (targetPosition.AveragePrice * float64(targetPosition.Quantity))) * 100
			targetPosition.UpdatedAt = time.Now()
		}

		portfolio.CashBalance += totalCost
	}

	// Recalculate total portfolio value
	totalValue := portfolio.CashBalance
	for _, pos := range portfolio.Positions {
		totalValue += pos.MarketValue
	}
	portfolio.TotalValue = totalValue
	portfolio.UpdatedAt = time.Now()

	c.JSON(http.StatusOK, gin.H{
		"message":   "Position updated successfully",
		"portfolio": portfolio,
		"timestamp": time.Now(),
	})
}

// Health returns the health status of the service
func (s *PortfolioService) Health(c *gin.Context) {
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
		Service:   "portfolio-service",
		Version:   "1.0.0",
		Timestamp: time.Now(),
		Checks: map[string]string{
			"database": "healthy",
			"cache":    "healthy",
		},
	}

	c.JSON(statusCode, response)
}

// SimulateFailure allows external control of failure simulation
func (s *PortfolioService) SimulateFailure(c *gin.Context) {
	var request struct {
		FailureRate float64 `json:"failure_rate"`
		IsHealthy   bool    `json:"is_healthy"`
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
	s.mutex.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message":      "Failure simulation updated",
		"failure_rate": request.FailureRate,
		"is_healthy":   request.IsHealthy,
		"timestamp":    time.Now(),
	})
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	service := NewPortfolioService()

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
		v1.GET("/portfolio/:userId", service.GetPortfolio)
		v1.POST("/portfolio/:userId/positions", service.UpdatePosition)
		v1.GET("/health", service.Health)
		v1.POST("/simulate/failure", service.SimulateFailure)
	}

	// Start server
	port := 8081
	fmt.Printf("🚀 Portfolio Service starting on port %d\n", port)
	fmt.Printf("📊 Available endpoints:\n")
	fmt.Printf("   GET  /api/v1/portfolio/{userId}           - Get user portfolio\n")
	fmt.Printf("   POST /api/v1/portfolio/{userId}/positions - Update position\n")
	fmt.Printf("   GET  /api/v1/health                       - Health check\n")
	fmt.Printf("   POST /api/v1/simulate/failure             - Simulate failures\n")
	fmt.Printf("\n💡 Example: curl http://localhost:%d/api/v1/portfolio/user123\n", port)

	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
