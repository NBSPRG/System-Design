package circuitbreaker

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

// State represents the circuit breaker state
type State int32

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

// Config holds circuit breaker configuration
type Config struct {
	Name                 string        `yaml:"name"`
	MaxRequests          uint32        `yaml:"max_requests"`           // Max requests allowed in half-open state
	Interval             time.Duration `yaml:"interval"`               // Statistical window duration
	Timeout              time.Duration `yaml:"timeout"`                // Time to wait before transitioning to half-open
	FailureThreshold     uint32        `yaml:"failure_threshold"`      // Number of failures to open circuit
	SuccessThreshold     uint32        `yaml:"success_threshold"`      // Number of successes to close circuit
	FailureRateThreshold float64       `yaml:"failure_rate_threshold"` // Failure rate (0.0-1.0) to open circuit
	MinimumRequests      uint32        `yaml:"minimum_requests"`       // Minimum requests before considering failure rate
}

// DefaultConfig returns a default configuration
func DefaultConfig(name string) Config {
	return Config{
		Name:                 name,
		MaxRequests:          5,
		Interval:             time.Minute,
		Timeout:              time.Minute,
		FailureThreshold:     10,
		SuccessThreshold:     3,
		FailureRateThreshold: 0.5,
		MinimumRequests:      5,
	}
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	config Config
	state  int32
	mutex  sync.RWMutex
	logger *zap.Logger

	// Counters
	requests  uint32
	failures  uint32
	successes uint32

	// Timing
	lastFailureTime time.Time
	lastStateChange time.Time

	// Half-open state tracking
	halfOpenRequests uint32

	// Metrics
	metrics *Metrics
}

// Metrics holds Prometheus metrics for the circuit breaker
type Metrics struct {
	requestsTotal   *prometheus.CounterVec
	failuresTotal   *prometheus.CounterVec
	stateChanges    *prometheus.CounterVec
	currentState    *prometheus.GaugeVec
	requestDuration *prometheus.HistogramVec
}

var (
	metricsOnce   sync.Once
	globalMetrics *Metrics
)

// NewCircuitBreaker creates a new circuit breaker instance
func NewCircuitBreaker(config Config, logger *zap.Logger) *CircuitBreaker {
	cb := &CircuitBreaker{
		config:          config,
		state:           int32(StateClosed),
		logger:          logger,
		lastStateChange: time.Now(),
		metrics:         getGlobalMetrics(),
	}

	cb.metrics.currentState.WithLabelValues(config.Name).Set(float64(StateClosed))

	return cb
}

func getGlobalMetrics() *Metrics {
	metricsOnce.Do(func() {
		globalMetrics = &Metrics{
			requestsTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "circuit_breaker_requests_total",
					Help: "Total number of requests processed by circuit breaker",
				},
				[]string{"service", "result"},
			),
			failuresTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "circuit_breaker_failures_total",
					Help: "Total number of failures",
				},
				[]string{"service"},
			),
			stateChanges: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "circuit_breaker_state_changes_total",
					Help: "Total number of state changes",
				},
				[]string{"service", "from", "to"},
			),
			currentState: promauto.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "circuit_breaker_state",
					Help: "Current state of the circuit breaker",
				},
				[]string{"service"},
			),
			requestDuration: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name: "circuit_breaker_request_duration_seconds",
					Help: "Request duration in seconds",
				},
				[]string{"service", "result"},
			),
		}
	})
	return globalMetrics
}

// Execute runs the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
	start := time.Now()
	// allowRequest determines if a request should be allowed through
	if !cb.allowRequest() {
		cb.metrics.requestsTotal.WithLabelValues(cb.config.Name, "rejected").Inc()
		return nil, errors.New("circuit breaker is open")
	}

	// Execute the function
	result, err := fn()
	duration := time.Since(start)
	// Record the result
	if err != nil {
		cb.onFailure()
		cb.metrics.requestsTotal.WithLabelValues(cb.config.Name, "failure").Inc()
		cb.metrics.failuresTotal.WithLabelValues(cb.config.Name).Inc()
		cb.metrics.requestDuration.WithLabelValues(cb.config.Name, "failure").Observe(duration.Seconds())
	} else {
		cb.onSuccess()
		cb.metrics.requestsTotal.WithLabelValues(cb.config.Name, "success").Inc()
		cb.metrics.requestDuration.WithLabelValues(cb.config.Name, "success").Observe(duration.Seconds())
	}

	return result, err
}

// allowRequest determines if a request should be allowed through
func (cb *CircuitBreaker) allowRequest() bool {
	currentState := State(atomic.LoadInt32(&cb.state))

	switch currentState {
	case StateClosed:
		return true
	case StateOpen:
		return cb.shouldAttemptReset()
	case StateHalfOpen:
		return cb.allowHalfOpenRequest()
	default:
		return false
	}
}

// shouldAttemptReset checks if enough time has passed to attempt reset
func (cb *CircuitBreaker) shouldAttemptReset() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	return time.Since(cb.lastFailureTime) >= cb.config.Timeout
}

// allowHalfOpenRequest checks if request is allowed in half-open state
func (cb *CircuitBreaker) allowHalfOpenRequest() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	return cb.halfOpenRequests < cb.config.MaxRequests
}

// onSuccess records a successful request
func (cb *CircuitBreaker) onSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	currentState := State(atomic.LoadInt32(&cb.state))

	switch currentState {
	case StateClosed:
		cb.resetCounters()
	case StateHalfOpen:
		cb.successes++
		cb.halfOpenRequests++

		if cb.successes >= cb.config.SuccessThreshold {
			cb.setState(StateClosed)
			cb.resetCounters()
		}
	}
}

// onFailure records a failed request
func (cb *CircuitBreaker) onFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	currentState := State(atomic.LoadInt32(&cb.state))
	cb.failures++
	cb.requests++
	cb.lastFailureTime = time.Now()

	switch currentState {
	case StateClosed:
		if cb.shouldOpenCircuit() {
			cb.setState(StateOpen)
		}
	case StateHalfOpen:
		cb.setState(StateOpen)
		cb.halfOpenRequests++
	}
}

// shouldOpenCircuit determines if the circuit should be opened
func (cb *CircuitBreaker) shouldOpenCircuit() bool {
	// Check failure threshold
	if cb.failures >= cb.config.FailureThreshold {
		return true
	}

	// Check failure rate
	if cb.requests >= cb.config.MinimumRequests {
		failureRate := float64(cb.failures) / float64(cb.requests)
		return failureRate >= cb.config.FailureRateThreshold
	}

	return false
}

// setState changes the circuit breaker state
func (cb *CircuitBreaker) setState(newState State) {
	oldState := State(atomic.LoadInt32(&cb.state))
	atomic.StoreInt32(&cb.state, int32(newState))

	cb.lastStateChange = time.Now()

	// Reset counters based on new state
	if newState == StateHalfOpen {
		cb.halfOpenRequests = 0
		cb.successes = 0
	}

	// Log state change
	cb.logger.Info("Circuit breaker state changed",
		zap.String("name", cb.config.Name),
		zap.String("from", oldState.String()),
		zap.String("to", newState.String()),
	)
	// Update metrics
	cb.metrics.stateChanges.WithLabelValues(cb.config.Name, oldState.String(), newState.String()).Inc()
	cb.metrics.currentState.WithLabelValues(cb.config.Name).Set(float64(newState))

	// Attempt transition to half-open if we're opening the circuit
	if newState == StateOpen {
		go cb.scheduleReset()
	}
}

// scheduleReset schedules a transition to half-open state
func (cb *CircuitBreaker) scheduleReset() {
	time.Sleep(cb.config.Timeout)

	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	currentState := State(atomic.LoadInt32(&cb.state))
	if currentState == StateOpen {
		cb.setState(StateHalfOpen)
	}
}

// resetCounters resets the failure and request counters
func (cb *CircuitBreaker) resetCounters() {
	cb.failures = 0
	cb.requests = 0
	cb.successes = 0
	cb.halfOpenRequests = 0
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() State {
	return State(atomic.LoadInt32(&cb.state))
}

// GetStats returns current statistics
func (cb *CircuitBreaker) GetStats() map[string]interface{} {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	return map[string]interface{}{
		"name":             cb.config.Name,
		"state":            cb.GetState().String(),
		"failures":         cb.failures,
		"requests":         cb.requests,
		"successes":        cb.successes,
		"halfOpenRequests": cb.halfOpenRequests,
		"lastStateChange":  cb.lastStateChange,
		"lastFailureTime":  cb.lastFailureTime,
	}
}
