package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"circuit-breaker-demo/pkg/circuitbreaker"

	"go.uber.org/zap"
)

// HTTPClient wraps http.Client with circuit breaker functionality
type HTTPClient struct {
	client         *http.Client
	circuitBreaker *circuitbreaker.CircuitBreaker
	logger         *zap.Logger
	baseURL        string
}

// NewHTTPClient creates a new HTTP client with circuit breaker
func NewHTTPClient(baseURL string, timeout time.Duration, cb *circuitbreaker.CircuitBreaker, logger *zap.Logger) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
		circuitBreaker: cb,
		logger:         logger,
		baseURL:        baseURL,
	}
}

// Get performs a GET request with circuit breaker protection
func (c *HTTPClient) Get(ctx context.Context, path string) (*http.Response, error) {
	return c.Do(ctx, "GET", path, nil, nil)
}

// Post performs a POST request with circuit breaker protection
func (c *HTTPClient) Post(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	return c.Do(ctx, "POST", path, body, map[string]string{
		"Content-Type": "application/json",
	})
}

// Put performs a PUT request with circuit breaker protection
func (c *HTTPClient) Put(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	return c.Do(ctx, "PUT", path, body, map[string]string{
		"Content-Type": "application/json",
	})
}

// Delete performs a DELETE request with circuit breaker protection
func (c *HTTPClient) Delete(ctx context.Context, path string) (*http.Response, error) {
	return c.Do(ctx, "DELETE", path, nil, nil)
}

// Do performs an HTTP request with circuit breaker protection
func (c *HTTPClient) Do(ctx context.Context, method, path string, body interface{}, headers map[string]string) (*http.Response, error) {
	result, err := c.circuitBreaker.Execute(ctx, func() (interface{}, error) {
		return c.doRequest(ctx, method, path, body, headers)
	})

	if err != nil {
		return nil, err
	}

	response, ok := result.(*http.Response)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}

	return response, nil
}

// doRequest performs the actual HTTP request
func (c *HTTPClient) doRequest(ctx context.Context, method, path string, body interface{}, headers map[string]string) (*http.Response, error) {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Add default headers
	req.Header.Set("User-Agent", "trading-gateway/1.0")
	req.Header.Set("Accept", "application/json")

	c.logger.Debug("Making HTTP request",
		zap.String("method", method),
		zap.String("url", url),
		zap.Any("headers", headers),
	)

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("HTTP request failed",
			zap.String("method", method),
			zap.String("url", url),
			zap.Error(err),
		)
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Check for HTTP error status codes
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		c.logger.Error("HTTP request returned error status",
			zap.String("method", method),
			zap.String("url", url),
			zap.Int("status_code", resp.StatusCode),
			zap.String("response_body", string(body)),
		)

		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	c.logger.Debug("HTTP request successful",
		zap.String("method", method),
		zap.String("url", url),
		zap.Int("status_code", resp.StatusCode),
	)

	return resp, nil
}

// GetJSON performs a GET request and unmarshals JSON response
func (c *HTTPClient) GetJSON(ctx context.Context, path string, target interface{}) error {
	resp, err := c.Get(ctx, path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

// PostJSON performs a POST request and unmarshals JSON response
func (c *HTTPClient) PostJSON(ctx context.Context, path string, body interface{}, target interface{}) error {
	resp, err := c.Post(ctx, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

// PutJSON performs a PUT request and unmarshals JSON response
func (c *HTTPClient) PutJSON(ctx context.Context, path string, body interface{}, target interface{}) error {
	resp, err := c.Put(ctx, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

// GetCircuitBreakerStats returns circuit breaker statistics
func (c *HTTPClient) GetCircuitBreakerStats() map[string]interface{} {
	return c.circuitBreaker.GetStats()
}

// GetCircuitBreakerState returns the current circuit breaker state
func (c *HTTPClient) GetCircuitBreakerState() circuitbreaker.State {
	return c.circuitBreaker.GetState()
}
