package models

import (
	"time"
)

// OrderType represents the type of trading order
type OrderType string

const (
	OrderTypeBuy  OrderType = "BUY"
	OrderTypeSell OrderType = "SELL"
)

// OrderStatus represents the status of a trading order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusExecuted  OrderStatus = "EXECUTED"
	OrderStatusRejected  OrderStatus = "REJECTED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// TradeRequest represents a request to execute a trade
type TradeRequest struct {
	UserID    string    `json:"userId" binding:"required"`
	Symbol    string    `json:"symbol" binding:"required"`
	Quantity  int64     `json:"quantity" binding:"required,min=1"`
	OrderType OrderType `json:"orderType" binding:"required"`
	Price     float64   `json:"price" binding:"required,min=0"`
	Timestamp time.Time `json:"timestamp"`
}

// TradeResponse represents the response after executing a trade
type TradeResponse struct {
	TradeID    string      `json:"tradeId"`
	UserID     string      `json:"userId"`
	Symbol     string      `json:"symbol"`
	Quantity   int64       `json:"quantity"`
	OrderType  OrderType   `json:"orderType"`
	Price      float64     `json:"price"`
	Status     OrderStatus `json:"status"`
	Message    string      `json:"message,omitempty"`
	ExecutedAt time.Time   `json:"executedAt"`
	TotalValue float64     `json:"totalValue"`
}

// Portfolio represents a user's portfolio
type Portfolio struct {
	UserID      string     `json:"userId"`
	Positions   []Position `json:"positions"`
	CashBalance float64    `json:"cashBalance"`
	TotalValue  float64    `json:"totalValue"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

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

// RiskCheckRequest represents a risk check request
type RiskCheckRequest struct {
	UserID     string    `json:"userId"`
	Symbol     string    `json:"symbol"`
	Quantity   int64     `json:"quantity"`
	OrderType  OrderType `json:"orderType"`
	Price      float64   `json:"price"`
	TotalValue float64   `json:"totalValue"`
}

// RiskCheckResponse represents a risk check response
type RiskCheckResponse struct {
	Approved  bool    `json:"approved"`
	Reason    string  `json:"reason,omitempty"`
	RiskScore float64 `json:"riskScore"`
}

// NotificationRequest represents a notification request
type NotificationRequest struct {
	UserID  string                 `json:"userId"`
	Type    string                 `json:"type"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// NotificationResponse represents a notification response
type NotificationResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"messageId,omitempty"`
	Error     string `json:"error,omitempty"`
}

// AuditEvent represents an audit event
type AuditEvent struct {
	EventID   string                 `json:"eventId"`
	UserID    string                 `json:"userId"`
	Action    string                 `json:"action"`
	Resource  string                 `json:"resource"`
	Details   map[string]interface{} `json:"details"`
	Timestamp time.Time              `json:"timestamp"`
	IPAddress string                 `json:"ipAddress,omitempty"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Service   string            `json:"service"`
	Version   string            `json:"version"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    map[string]string `json:"checks,omitempty"`
}

// CircuitBreakerStatus represents circuit breaker status
type CircuitBreakerStatus struct {
	Name             string                 `json:"name"`
	State            string                 `json:"state"`
	Failures         uint32                 `json:"failures"`
	Requests         uint32                 `json:"requests"`
	Successes        uint32                 `json:"successes"`
	HalfOpenRequests uint32                 `json:"halfOpenRequests"`
	LastStateChange  time.Time              `json:"lastStateChange"`
	LastFailureTime  time.Time              `json:"lastFailureTime"`
	Configuration    map[string]interface{} `json:"configuration"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error     string                 `json:"error"`
	Message   string                 `json:"message"`
	Code      string                 `json:"code,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// BatchMarketDataRequest represents a batch market data request
type BatchMarketDataRequest struct {
	Symbols []string `json:"symbols" binding:"required"`
}

// BatchMarketDataResponse represents a batch market data response
type BatchMarketDataResponse struct {
	Data      map[string]MarketData `json:"data"`
	Errors    map[string]string     `json:"errors,omitempty"`
	Timestamp time.Time             `json:"timestamp"`
}
