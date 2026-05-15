// Package config handles loading and validation of application configuration.
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration values.
type Config struct {
	// Server settings
	Port     int
	Host     string
	LogLevel string

	// AI provider settings
	OpenAIAPIKey   string
	OpenAIModel    string
	MaxTokens      int
	Temperature    float64

	// Engram memory settings
	EngramPath     string
	MemoryEnabled  bool
}

// Load reads configuration from environment variables with sensible defaults.
func Load() (*Config, error) {
	cfg := &Config{
		Port:          getEnvInt("PORT", 8080),
		Host:          getEnv("HOST", "0.0.0.0"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		OpenAIAPIKey:  os.Getenv("OPENAI_API_KEY"),
		OpenAIModel:   getEnv("OPENAI_MODEL", "gpt-4o-mini"),
		MaxTokens:     getEnvInt("MAX_TOKENS", 2048),
		Temperature:   getEnvFloat("TEMPERATURE", 0.7),
		EngramPath:    getEnv("ENGRAM_PATH", ".engram"),
		MemoryEnabled: getEnvBool("MEMORY_ENABLED", true),
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// validate checks that required configuration values are present.
func (c *Config) validate() error {
	if c.OpenAIAPIKey == "" {
		return fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("PORT must be between 1 and 65535, got %d", c.Port)
	}
	if c.Temperature < 0.0 || c.Temperature > 2.0 {
		return fmt.Errorf("TEMPERATURE must be between 0.0 and 2.0, got %f", c.Temperature)
	}
	return nil
}

// Addr returns the full host:port address string.
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return n
}

func getEnvFloat(key string, defaultVal float64) float64 {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return defaultVal
	}
	return f
}

func getEnvBool(key string, defaultVal bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultVal
	}
	return b
}
