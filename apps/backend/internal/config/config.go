package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds API runtime configuration loaded from environment variables.
type Config struct {
	Addr           string
	AdminToken     string
	DatabaseDSN    string
	DatabaseDriver string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

// LoadFromEnv builds Config from environment variables.
func LoadFromEnv() (Config, error) {
	cfg := Config{
		Addr:           getEnv("OPS_API_ADDR", ":8080"),
		AdminToken:     os.Getenv("OPS_API_ADMIN_TOKEN"),
		DatabaseDSN:    os.Getenv("OPS_API_DB_DSN"),
		DatabaseDriver: getEnv("OPS_API_DB_DRIVER", "postgres"),
		ReadTimeout:    time.Duration(getEnvInt("OPS_API_READ_TIMEOUT_SEC", 10)) * time.Second,
		WriteTimeout:   time.Duration(getEnvInt("OPS_API_WRITE_TIMEOUT_SEC", 10)) * time.Second,
	}

	if cfg.ReadTimeout <= 0 || cfg.WriteTimeout <= 0 {
		return Config{}, fmt.Errorf("timeouts must be greater than zero")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return parsed
}
