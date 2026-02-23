package config

import "os"

// Config holds application configuration
type Config struct {
	Port   string
	DBPath string
}

// Load loads configuration from environment variables or defaults
func Load() *Config {
	return &Config{
		Port:   getEnv("PORT", "3000"),
		DBPath: getEnv("DB_PATH", "./petstore.db"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
