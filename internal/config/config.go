package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Config stores all configuration of the application.
// Values are read from environment variables with sensible defaults.
type Config struct {
	Environment        string        `json:"environment"`
	DBSource           string        `json:"db_source"`
	HTTPServerAddress  string        `json:"http_server_address"`
	ReadTimeout        time.Duration `json:"read_timeout"`
	WriteTimeout       time.Duration `json:"write_timeout"`
	IdleTimeout        time.Duration `json:"idle_timeout"`
	ShutdownTimeout    time.Duration `json:"shutdown_timeout"`
	RunMigrationOnBoot bool          `json:"run_migration_on_boot"`
}

// LoadConfig reads configuration from environment variables or returns defaults.
func LoadConfig() (Config, error) {
	cfg := Config{
		Environment:        getEnv("ENVIRONMENT", "development"),
		DBSource:           getEnv("DB_SOURCE", "postgres://postgres:secretpassword@localhost:5432/ledger_db?sslmode=disable"),
		HTTPServerAddress:  getEnv("HTTP_SERVER_ADDRESS", ":8080"),
		ReadTimeout:        getEnvDuration("HTTP_READ_TIMEOUT", 5*time.Second),
		WriteTimeout:       getEnvDuration("HTTP_WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:        getEnvDuration("HTTP_IDLE_TIMEOUT", 60*time.Second),
		ShutdownTimeout:    getEnvDuration("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second),
		RunMigrationOnBoot: getEnvBool("RUN_MIGRATION_ON_BOOT", true),
	}

	if strings.TrimSpace(cfg.DBSource) == "" {
		return Config{}, fmt.Errorf("DB_SOURCE environment variable must not be empty")
	}
	if strings.TrimSpace(cfg.HTTPServerAddress) == "" {
		return Config{}, fmt.Errorf("HTTP_SERVER_ADDRESS must not be empty")
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && strings.TrimSpace(val) != "" {
		return strings.TrimSpace(val)
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val, ok := os.LookupEnv(key); ok && strings.TrimSpace(val) != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if val, ok := os.LookupEnv(key); ok && strings.TrimSpace(val) != "" {
		lower := strings.ToLower(strings.TrimSpace(val))
		return lower == "true" || lower == "1" || lower == "yes"
	}
	return defaultVal
}
