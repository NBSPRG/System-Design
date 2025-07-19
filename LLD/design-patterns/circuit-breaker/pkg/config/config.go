package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	Server struct {
		Port         int           `yaml:"port"`
		ReadTimeout  time.Duration `yaml:"read_timeout"`
		WriteTimeout time.Duration `yaml:"write_timeout"`
	} `yaml:"server"`

	CircuitBreaker struct {
		MaxRequests          uint32        `yaml:"max_requests"`
		Interval             time.Duration `yaml:"interval"`
		Timeout              time.Duration `yaml:"timeout"`
		FailureThreshold     uint32        `yaml:"failure_threshold"`
		SuccessThreshold     uint32        `yaml:"success_threshold"`
		FailureRateThreshold float64       `yaml:"failure_rate_threshold"`
		MinimumRequests      uint32        `yaml:"minimum_requests"`
	} `yaml:"circuit_breaker"`

	Services struct {
		MarketData struct {
			URL     string        `yaml:"url"`
			Timeout time.Duration `yaml:"timeout"`
		} `yaml:"market_data"`

		RiskManagement struct {
			URL     string        `yaml:"url"`
			Timeout time.Duration `yaml:"timeout"`
		} `yaml:"risk_management"`

		Notification struct {
			URL     string        `yaml:"url"`
			Timeout time.Duration `yaml:"timeout"`
		} `yaml:"notification"`

		Audit struct {
			URL     string        `yaml:"url"`
			Timeout time.Duration `yaml:"timeout"`
		} `yaml:"audit"`
	} `yaml:"services"`

	Logging struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
	} `yaml:"logging"`

	Metrics struct {
		Enabled bool   `yaml:"enabled"`
		Port    int    `yaml:"port"`
		Path    string `yaml:"path"`
	} `yaml:"metrics"`
}

// LoadConfig loads configuration from file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: struct {
			Port         int           `yaml:"port"`
			ReadTimeout  time.Duration `yaml:"read_timeout"`
			WriteTimeout time.Duration `yaml:"write_timeout"`
		}{
			Port:         8080,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		CircuitBreaker: struct {
			MaxRequests          uint32        `yaml:"max_requests"`
			Interval             time.Duration `yaml:"interval"`
			Timeout              time.Duration `yaml:"timeout"`
			FailureThreshold     uint32        `yaml:"failure_threshold"`
			SuccessThreshold     uint32        `yaml:"success_threshold"`
			FailureRateThreshold float64       `yaml:"failure_rate_threshold"`
			MinimumRequests      uint32        `yaml:"minimum_requests"`
		}{
			MaxRequests:          5,
			Interval:             time.Minute,
			Timeout:              30 * time.Second,
			FailureThreshold:     10,
			SuccessThreshold:     3,
			FailureRateThreshold: 0.6,
			MinimumRequests:      5,
		},
		Services: struct {
			MarketData struct {
				URL     string        `yaml:"url"`
				Timeout time.Duration `yaml:"timeout"`
			} `yaml:"market_data"`
			RiskManagement struct {
				URL     string        `yaml:"url"`
				Timeout time.Duration `yaml:"timeout"`
			} `yaml:"risk_management"`
			Notification struct {
				URL     string        `yaml:"url"`
				Timeout time.Duration `yaml:"timeout"`
			} `yaml:"notification"`
			Audit struct {
				URL     string        `yaml:"url"`
				Timeout time.Duration `yaml:"timeout"`
			} `yaml:"audit"`
		}{
			MarketData: struct {
				URL     string        `yaml:"url"`
				Timeout time.Duration `yaml:"timeout"`
			}{
				URL:     "http://localhost:8082",
				Timeout: 5 * time.Second,
			},
			RiskManagement: struct {
				URL     string        `yaml:"url"`
				Timeout time.Duration `yaml:"timeout"`
			}{
				URL:     "http://localhost:8083",
				Timeout: 3 * time.Second,
			},
			Notification: struct {
				URL     string        `yaml:"url"`
				Timeout time.Duration `yaml:"timeout"`
			}{
				URL:     "http://localhost:8084",
				Timeout: 2 * time.Second,
			},
			Audit: struct {
				URL     string        `yaml:"url"`
				Timeout time.Duration `yaml:"timeout"`
			}{
				URL:     "http://localhost:8085",
				Timeout: 3 * time.Second,
			},
		},
		Logging: struct {
			Level  string `yaml:"level"`
			Format string `yaml:"format"`
		}{
			Level:  "info",
			Format: "json",
		},
		Metrics: struct {
			Enabled bool   `yaml:"enabled"`
			Port    int    `yaml:"port"`
			Path    string `yaml:"path"`
		}{
			Enabled: true,
			Port:    9090,
			Path:    "/metrics",
		},
	}
}
