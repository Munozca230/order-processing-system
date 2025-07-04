package configs

import (
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Logging  LoggingConfig  `json:"logging"`
	Features FeatureFlags   `json:"features"`
	Cache    CacheConfig    `json:"cache"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port            string        `json:"port"`
	Host            string        `json:"host"`
	ReadTimeout     time.Duration `json:"readTimeout"`
	WriteTimeout    time.Duration `json:"writeTimeout"`
	ShutdownTimeout time.Duration `json:"shutdownTimeout"`
	Environment     string        `json:"environment"`
	Version         string        `json:"version"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Type           string        `json:"type"`
	URL            string        `json:"url"`
	Database       string        `json:"database"`
	Collection     string        `json:"collection"`
	MaxConnections int           `json:"maxConnections"`
	Timeout        time.Duration `json:"timeout"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      string `json:"level"`
	Format     string `json:"format"`     // json, text
	Output     string `json:"output"`     // stdout, file
	RequestLog bool   `json:"requestLog"`
}

// FeatureFlags holds feature toggles
type FeatureFlags struct {
	EnableMetrics     bool    `json:"enableMetrics"`
	EnableTracing     bool    `json:"enableTracing"`
	SimulateLatency   bool    `json:"simulateLatency"`
	SimulateErrors    bool    `json:"simulateErrors"`
	ErrorRate         float64 `json:"errorRate"`
	MaxLatencyMs      int     `json:"maxLatencyMs"`
	EnableHealthCheck bool    `json:"enableHealthCheck"`
}

// CacheConfig holds caching configuration
type CacheConfig struct {
	Enabled bool          `json:"enabled"`
	TTL     time.Duration `json:"ttl"`
	MaxSize int           `json:"maxSize"`
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", "8080"),
			Host:            getEnv("HOST", "0.0.0.0"),
			ReadTimeout:     getDurationEnv("READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getDurationEnv("WRITE_TIMEOUT", 10*time.Second),
			ShutdownTimeout: getDurationEnv("SHUTDOWN_TIMEOUT", 30*time.Second),
			Environment:     getEnv("ENVIRONMENT", "development"),
			Version:         getEnv("VERSION", "1.0.0"),
		},
		Database: DatabaseConfig{
			Type:           getEnv("DATABASE_TYPE", "mongodb"),
			URL:            getEnv("DATABASE_URL", "mongodb://mongo:27017"),
			Database:       getEnv("DATABASE_NAME", "catalog"),
			Collection:     getEnv("DATABASE_COLLECTION", "products"),
			MaxConnections: getIntEnv("DATABASE_MAX_CONNECTIONS", 10),
			Timeout:        getDurationEnv("DATABASE_TIMEOUT", 5*time.Second),
		},
		Logging: LoggingConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			Format:     getEnv("LOG_FORMAT", "json"),
			Output:     getEnv("LOG_OUTPUT", "stdout"),
			RequestLog: getBoolEnv("LOG_REQUESTS", true),
		},
		Features: FeatureFlags{
			EnableMetrics:     getBoolEnv("ENABLE_METRICS", true),
			EnableTracing:     getBoolEnv("ENABLE_TRACING", false),
			SimulateLatency:   getBoolEnv("SIMULATE_LATENCY", false),
			SimulateErrors:    getBoolEnv("SIMULATE_ERRORS", false),
			ErrorRate:         getFloatEnv("ERROR_RATE", 0.0),
			MaxLatencyMs:      getIntEnv("MAX_LATENCY_MS", 200),
			EnableHealthCheck: getBoolEnv("ENABLE_HEALTH_CHECK", true),
		},
		Cache: CacheConfig{
			Enabled: getBoolEnv("CACHE_ENABLED", false),
			TTL:     getDurationEnv("CACHE_TTL", 5*time.Minute),
			MaxSize: getIntEnv("CACHE_MAX_SIZE", 1000),
		},
	}
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getFloatEnv(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}