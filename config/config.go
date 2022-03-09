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
	CacheSize                  int           `envconfig:"CACHE_SIZE"`
	CacheTTL                   time.Duration `envconfig:"CACHE_TTL"`
	EnableCantabular           bool          `envconfig:"ENABLE_CANTABULAR"`
	CantabularURL              string        `envconfig:"CANT_URL"`
	CantabularUser             string        `envconfig:"CANT_USER"`
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
		CacheSize:                  200,            // memory cache size in MB
		CacheTTL:                   12 * time.Hour, // cache entry TTL
		// Cantabular defaults to disabled, so no defaults
	}

	return cfg, envconfig.Process("", cfg)
}
