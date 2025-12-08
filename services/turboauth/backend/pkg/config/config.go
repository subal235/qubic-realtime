package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	// Server
	HTTPPort int
	GRPCPort int
	Env      string

	// Qubic
	QubicNodeURL      string
	QubicContractAddr string

	// Redis
	RedisURL      string
	RedisPassword string
	RedisDB       int

	// Cache
	CacheTTL       time.Duration
	UseMemoryCache bool

	// Logging
	LogLevel  string
	LogFormat string

	// Metrics
	MetricsEnabled bool
	MetricsPort    int
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		HTTPPort:          getEnvAsInt("HTTP_PORT", 8080),
		GRPCPort:          getEnvAsInt("GRPC_PORT", 9090),
		Env:               getEnv("ENV", "development"),
		QubicNodeURL:      getEnv("QUBIC_NODE_URL", "http://localhost:21841"),
		QubicContractAddr: getEnv("QUBIC_CONTRACT_ADDRESS", ""),
		RedisURL:          getEnv("REDIS_URL", "localhost:6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		RedisDB:           getEnvAsInt("REDIS_DB", 0),
		CacheTTL:          time.Duration(getEnvAsInt("CACHE_TTL_SECONDS", 300)) * time.Second,
		UseMemoryCache:    getEnvAsBool("USE_MEMORY_CACHE", true),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		LogFormat:         getEnv("LOG_FORMAT", "json"),
		MetricsEnabled:    getEnvAsBool("METRICS_ENABLED", true),
		MetricsPort:       getEnvAsInt("METRICS_PORT", 2112),
	}
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func getEnvAsBool(key string, defaultVal bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultVal
}
