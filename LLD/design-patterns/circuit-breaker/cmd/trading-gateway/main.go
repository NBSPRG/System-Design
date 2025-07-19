package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"circuit-breaker-demo/pkg/circuitbreaker"
	"circuit-breaker-demo/pkg/httpclient"
	"circuit-breaker-demo/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// TradingGateway represents the main trading gateway service
type TradingGateway struct {
	logger               *zap.Logger
	marketDataClient     *httpclient.HTTPClient
	portfolioClient      *httpclient.HTTPClient
	riskManagementClient *httpclient.HTTPClient
	notificationClient   *httpclient.HTTPClient
	auditClient          *httpclient.HTTPClient
}

// NewTradingGateway creates a new trading gateway instance
func NewTradingGateway(logger *zap.Logger) *TradingGateway {
	// Create circuit breakers for each service
	marketDataCB := circuitbreaker.NewCircuitBreaker(
		circuitbreaker.Config{
			Name:                 "market-data-service",
			MaxRequests:          3,
			Interval:             time.Minute,
			Timeout:              30 * time.Second,
			FailureThreshold:     5,
			SuccessThreshold:     2,
			FailureRateThreshold: 0.5,
			MinimumRequests:      3,
		},
		logger,
	)

	portfolioCB := circuitbreaker.NewCircuitBreaker(
		circuitbreaker.Config{
			Name:                 "portfolio-service",
			MaxRequests:          5,
			Interval:             time.Minute,
			Timeout:              20 * time.Second,
			FailureThreshold:     3,
			SuccessThreshold:     2,
			FailureRateThreshold: 0.6,
			MinimumRequests:      2,
		},
		logger,
	)

	riskMgmtCB := circuitbreaker.NewCircuitBreaker(
		circuitbreaker.Config{
			Name:                 "risk-management-service",
			MaxRequests:          2,
			Interval:             time.Minute,
			Timeout:              15 * time.Second,
			FailureThreshold:     2,
			SuccessThreshold:     1,
			FailureRateThreshold: 0.3,
			MinimumRequests:      2,
		},
		logger,
	)

	notificationCB := circuitbreaker.NewCircuitBreaker(
		circuitbreaker.Config{
			Name:                 "notification-service",
			MaxRequests:          10,
			Interval:             time.Minute,
			Timeout:              10 * time.Second,
			FailureThreshold:     10,
			SuccessThreshold:     3,
			FailureRateThreshold: 0.8,
			MinimumRequests:      5,
		},
		logger,
	)

	auditCB := circuitbreaker.NewCircuitBreaker(
		circuitbreaker.Config{
			Name:                 "audit-service",
			MaxRequests:          15,
			Interval:             time.Minute,
			Timeout:              5 * time.Second,
			FailureThreshold:     15,
			SuccessThreshold:     5,
			FailureRateThreshold: 0.9,
			MinimumRequests:      10,
		},
		logger,
	)

	return &TradingGateway{
		logger:               logger,
		marketDataClient:     httpclient.NewHTTPClient("http://localhost:8082", 5*time.Second, marketDataCB, logger),
		portfolioClient:      httpclient.NewHTTPClient("http://localhost:8081", 5*time.Second, portfolioCB, logger),
		riskManagementClient: httpclient.NewHTTPClient("http://localhost:8083", 3*time.Second, riskMgmtCB, logger),
		notificationClient:   httpclient.NewHTTPClient("http://localhost:8084", 2*time.Second, notificationCB, logger),
		auditClient:          httpclient.NewHTTPClient("http://localhost:8085", 3*time.Second, auditCB, logger),
	}
}

// ExecuteTrade handles trade execution requests
func (tg *TradingGateway) ExecuteTrade(c *gin.Context) {
	var request models.TradeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		tg.logger.Error("Invalid trade request", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:     "Invalid request format",
			Message:   err.Error(),
			Code:      "INVALID_REQUEST",
			Timestamp: time.Now(),
		})
		return
	}

	request.Timestamp = time.Now()
	ctx := context.Background()

	tg.logger.Info("Processing trade request",
		zap.String("userId", request.UserID),
		zap.String("symbol", request.Symbol),
		zap.Int64("quantity", request.Quantity),
		zap.String("orderType", string(request.OrderType)),
		zap.Float64("price", request.Price),
	)

	// Step 1: Get current market data
	var marketData models.MarketData
	err := tg.marketDataClient.GetJSON(ctx, fmt.Sprintf("/api/v1/prices/%s", request.Symbol), &marketData)
	if err != nil {
		tg.logger.Error("Failed to get market data", zap.Error(err))

		// Fallback: use the price from the request
		tg.logger.Warn("Using fallback price due to market data service failure")
		marketData = models.MarketData{
			Symbol:    request.Symbol,
			Price:     request.Price,
			Timestamp: time.Now(),
		}
	}

	// Step 2: Risk management check
	totalValue := request.Price * float64(request.Quantity)
	riskRequest := models.RiskCheckRequest{
		UserID:     request.UserID,
		Symbol:     request.Symbol,
		Quantity:   request.Quantity,
		OrderType:  request.OrderType,
		Price:      request.Price,
		TotalValue: totalValue,
	}

	var riskResponse models.RiskCheckResponse
	err = tg.riskManagementClient.PostJSON(ctx, "/api/v1/risk/check", riskRequest, &riskResponse)
	if err != nil {
		tg.logger.Error("Risk management check failed", zap.Error(err))

		// Fallback: apply basic risk rules
		riskResponse = tg.fallbackRiskCheck(riskRequest)
	}

	if !riskResponse.Approved {
		tg.logger.Warn("Trade rejected by risk management",
			zap.String("reason", riskResponse.Reason),
			zap.Float64("riskScore", riskResponse.RiskScore),
		)

		c.JSON(http.StatusBadRequest, models.TradeResponse{
			UserID:    request.UserID,
			Symbol:    request.Symbol,
			Quantity:  request.Quantity,
			OrderType: request.OrderType,
			Price:     request.Price,
			Status:    models.OrderStatusRejected,
			Message:   fmt.Sprintf("Trade rejected: %s", riskResponse.Reason),
		})
		return
	}

	// Step 3: Update portfolio
	portfolioUpdateRequest := map[string]interface{}{
		"symbol":   request.Symbol,
		"quantity": request.Quantity,
		"price":    marketData.Price, // Use current market price
		"action":   string(request.OrderType),
	}

	var portfolioResponse map[string]interface{}
	err = tg.portfolioClient.PostJSON(ctx, fmt.Sprintf("/api/v1/portfolio/%s/positions", request.UserID), portfolioUpdateRequest, &portfolioResponse)
	if err != nil {
		tg.logger.Error("Failed to update portfolio", zap.Error(err))

		c.JSON(http.StatusInternalServerError, models.TradeResponse{
			UserID:    request.UserID,
			Symbol:    request.Symbol,
			Quantity:  request.Quantity,
			OrderType: request.OrderType,
			Price:     request.Price,
			Status:    models.OrderStatusRejected,
			Message:   "Failed to update portfolio",
		})
		return
	}

	// Step 4: Generate trade ID and response
	tradeID := fmt.Sprintf("TXN_%d_%s", time.Now().Unix(), request.UserID)

	tradeResponse := models.TradeResponse{
		TradeID:    tradeID,
		UserID:     request.UserID,
		Symbol:     request.Symbol,
		Quantity:   request.Quantity,
		OrderType:  request.OrderType,
		Price:      marketData.Price,
		Status:     models.OrderStatusExecuted,
		Message:    "Trade executed successfully",
		ExecutedAt: time.Now(),
		TotalValue: marketData.Price * float64(request.Quantity),
	}

	// Step 5: Send notification (async, non-blocking)
	go func() {
		notificationRequest := models.NotificationRequest{
			UserID:  request.UserID,
			Type:    "TRADE_EXECUTED",
			Message: fmt.Sprintf("Trade executed: %s %d shares of %s at $%.2f", request.OrderType, request.Quantity, request.Symbol, marketData.Price),
			Data: map[string]interface{}{
				"tradeId":    tradeID,
				"symbol":     request.Symbol,
				"quantity":   request.Quantity,
				"orderType":  request.OrderType,
				"price":      marketData.Price,
				"totalValue": tradeResponse.TotalValue,
			},
		}

		var notificationResponse models.NotificationResponse
		if err := tg.notificationClient.PostJSON(context.Background(), "/api/v1/notifications", notificationRequest, &notificationResponse); err != nil {
			tg.logger.Error("Failed to send notification", zap.Error(err))
		}
	}()

	// Step 6: Audit log (async, non-blocking)
	go func() {
		auditEvent := models.AuditEvent{
			EventID:  fmt.Sprintf("AUDIT_%d", time.Now().UnixNano()),
			UserID:   request.UserID,
			Action:   "TRADE_EXECUTED",
			Resource: "trading-gateway",
			Details: map[string]interface{}{
				"tradeId":    tradeID,
				"symbol":     request.Symbol,
				"quantity":   request.Quantity,
				"orderType":  request.OrderType,
				"price":      marketData.Price,
				"totalValue": tradeResponse.TotalValue,
				"riskScore":  riskResponse.RiskScore,
			},
			Timestamp: time.Now(),
			IPAddress: c.ClientIP(),
		}

		if err := tg.auditClient.PostJSON(context.Background(), "/api/v1/audit", auditEvent, nil); err != nil {
			tg.logger.Error("Failed to log audit event", zap.Error(err))
		}
	}()

	tg.logger.Info("Trade executed successfully",
		zap.String("tradeId", tradeID),
		zap.String("userId", request.UserID),
		zap.String("symbol", request.Symbol),
		zap.Float64("executedPrice", marketData.Price),
	)

	c.JSON(http.StatusOK, tradeResponse)
}

// fallbackRiskCheck provides basic risk checking when the risk service is unavailable
func (tg *TradingGateway) fallbackRiskCheck(request models.RiskCheckRequest) models.RiskCheckResponse {
	// Basic risk rules
	maxTradeValue := 50000.0 // Max $50k per trade

	if request.TotalValue > maxTradeValue {
		return models.RiskCheckResponse{
			Approved:  false,
			Reason:    fmt.Sprintf("Trade value exceeds maximum allowed ($%.2f)", maxTradeValue),
			RiskScore: 0.9,
		}
	}

	return models.RiskCheckResponse{
		Approved:  true,
		Reason:    "Approved by fallback risk check",
		RiskScore: 0.3,
	}
}

// GetPortfolio returns a user's portfolio
func (tg *TradingGateway) GetPortfolio(c *gin.Context) {
	userID := c.Param("userId")
	ctx := context.Background()

	var portfolio models.Portfolio
	err := tg.portfolioClient.GetJSON(ctx, fmt.Sprintf("/api/v1/portfolio/%s", userID), &portfolio)
	if err != nil {
		tg.logger.Error("Failed to get portfolio", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:     "Failed to retrieve portfolio",
			Message:   err.Error(),
			Code:      "PORTFOLIO_SERVICE_ERROR",
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

// GetMarketData returns market data for a symbol
func (tg *TradingGateway) GetMarketData(c *gin.Context) {
	symbol := c.Param("symbol")
	ctx := context.Background()

	var marketData models.MarketData
	err := tg.marketDataClient.GetJSON(ctx, fmt.Sprintf("/api/v1/prices/%s", symbol), &marketData)
	if err != nil {
		tg.logger.Error("Failed to get market data", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:     "Failed to retrieve market data",
			Message:   err.Error(),
			Code:      "MARKET_DATA_SERVICE_ERROR",
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, marketData)
}

// GetCircuitBreakerStatus returns the status of all circuit breakers
func (tg *TradingGateway) GetCircuitBreakerStatus(c *gin.Context) {
	status := map[string]interface{}{
		"market_data_service":     tg.marketDataClient.GetCircuitBreakerStats(),
		"portfolio_service":       tg.portfolioClient.GetCircuitBreakerStats(),
		"risk_management_service": tg.riskManagementClient.GetCircuitBreakerStats(),
		"notification_service":    tg.notificationClient.GetCircuitBreakerStats(),
		"audit_service":           tg.auditClient.GetCircuitBreakerStats(),
		"timestamp":               time.Now(),
	}

	c.JSON(http.StatusOK, status)
}

// Health returns the health status of the gateway
func (tg *TradingGateway) Health(c *gin.Context) {
	response := models.HealthResponse{
		Status:    "healthy",
		Service:   "trading-gateway",
		Version:   "1.0.0",
		Timestamp: time.Now(),
		Checks: map[string]string{
			"market_data_circuit_breaker":     tg.marketDataClient.GetCircuitBreakerState().String(),
			"portfolio_circuit_breaker":       tg.portfolioClient.GetCircuitBreakerState().String(),
			"risk_management_circuit_breaker": tg.riskManagementClient.GetCircuitBreakerState().String(),
			"notification_circuit_breaker":    tg.notificationClient.GetCircuitBreakerState().String(),
			"audit_circuit_breaker":           tg.auditClient.GetCircuitBreakerState().String(),
		},
	}

	c.JSON(http.StatusOK, response)
}

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Create trading gateway
	gateway := NewTradingGateway(logger)

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

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/trades", gateway.ExecuteTrade)
		v1.GET("/portfolio/:userId", gateway.GetPortfolio)
		v1.GET("/market-data/:symbol", gateway.GetMarketData)
		v1.GET("/circuit-breaker/status", gateway.GetCircuitBreakerStatus)
		v1.GET("/health", gateway.Health)
	}

	// Start server
	port := 8080
	fmt.Printf("🚀 Trading Gateway starting on port %d\n", port)
	fmt.Printf("📊 Available endpoints:\n")
	fmt.Printf("   POST /api/v1/trades                    - Execute trade\n")
	fmt.Printf("   GET  /api/v1/portfolio/{userId}        - Get portfolio\n")
	fmt.Printf("   GET  /api/v1/market-data/{symbol}      - Get market data\n")
	fmt.Printf("   GET  /api/v1/circuit-breaker/status    - Circuit breaker status\n")
	fmt.Printf("   GET  /api/v1/health                     - Health check\n")
	fmt.Printf("   GET  /metrics                           - Prometheus metrics\n")
	fmt.Printf("\n💡 Example trade: curl -X POST http://localhost:%d/api/v1/trades \\\n", port)
	fmt.Printf("   -H 'Content-Type: application/json' \\\n")
	fmt.Printf("   -d '{\"userId\":\"user123\",\"symbol\":\"AAPL\",\"quantity\":10,\"orderType\":\"BUY\",\"price\":150.00}'\n")

	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
