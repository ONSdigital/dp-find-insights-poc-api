package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-find-insights-poc-api
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	EnableDatabase             bool          `envconfig:"ENABLE_DATABASE"`
	MaxMetrics                 int           `envconfig:"MAX_METRICS"`
	WriteTimeout               time.Duration `envconfig:"WRITE_TIMEOUT"`
	APIToken                   string        `envconfig:"API_TOKEN"`
	EnableHeaderAuth           bool          `envconfig:"ENABLE_HEADER_AUTH"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   "localhost:25252",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		MaxMetrics:                 200000,           // max number of rows to accept from "geo" table queries
		WriteTimeout:               30 * time.Second, // http WriteTimeout
		APIToken:                   "",
		EnableHeaderAuth:           false,
	}

	return cfg, envconfig.Process("", cfg)
}
