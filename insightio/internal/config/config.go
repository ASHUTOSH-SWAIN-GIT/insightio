package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	GRPCPort      int
	MetricsWindow int
	APIKey        string
	APIKeys       []string // Parsed API keys (supports comma-separated)
	Env           string
}

func Load() *Config {
	cfg := &Config{
		GRPCPort:      getEnvAsInt("INSIGHTIO_GRPC_PORT", 50051),
		MetricsWindow: getEnvAsInt("INSIGHTIO_METRICS_WINDOW", 60),
		APIKey:        getEnv("INSIGHTIO_API_KEY", ""),
		Env:           getEnv("INSIGHTIO_ENV", "dev"),
	}

	// Parse API keys - support comma-separated values
	if cfg.APIKey == "" {
		log.Fatal("INSIGHTIO_API_KEY must be set")
	}

	// Split by comma and trim whitespace
	keys := strings.Split(cfg.APIKey, ",")
	cfg.APIKeys = make([]string, 0, len(keys))
	for _, key := range keys {
		trimmed := strings.TrimSpace(key)
		if trimmed != "" {
			cfg.APIKeys = append(cfg.APIKeys, trimmed)
		}
	}

	if len(cfg.APIKeys) == 0 {
		log.Fatal("INSIGHTIO_API_KEY must contain at least one valid API key")
	}

	return cfg
}
func getEnv(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultVal
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Fatalf("Invalid value for %s", key)
	}

	return val
}
