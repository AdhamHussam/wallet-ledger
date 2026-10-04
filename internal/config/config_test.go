package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfigDefaults(t *testing.T) {
	// Clear any env vars that might interfere
	os.Unsetenv("ENVIRONMENT")
	os.Unsetenv("DB_SOURCE")
	os.Unsetenv("HTTP_SERVER_ADDRESS")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error loading defaults: %v", err)
	}

	if cfg.Environment != "development" {
		t.Errorf("expected development, got %s", cfg.Environment)
	}
	if cfg.HTTPServerAddress != ":8080" {
		t.Errorf("expected :8080, got %s", cfg.HTTPServerAddress)
	}
	if cfg.ReadTimeout != 5*time.Second {
		t.Errorf("expected 5s read timeout, got %v", cfg.ReadTimeout)
	}
	if !cfg.RunMigrationOnBoot {
		t.Errorf("expected RunMigrationOnBoot=true by default")
	}
}

func TestLoadConfigCustomEnv(t *testing.T) {
	t.Setenv("ENVIRONMENT", "production")
	t.Setenv("DB_SOURCE", "postgres://custom:pass@customhost:5432/customdb?sslmode=require")
	t.Setenv("HTTP_SERVER_ADDRESS", "0.0.0.0:9090")
	t.Setenv("HTTP_READ_TIMEOUT", "15s")
	t.Setenv("RUN_MIGRATION_ON_BOOT", "false")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Environment != "production" {
		t.Errorf("expected production, got %s", cfg.Environment)
	}
	if cfg.HTTPServerAddress != "0.0.0.0:9090" {
		t.Errorf("expected 0.0.0.0:9090, got %s", cfg.HTTPServerAddress)
	}
	if cfg.ReadTimeout != 15*time.Second {
		t.Errorf("expected 15s read timeout, got %v", cfg.ReadTimeout)
	}
	if cfg.RunMigrationOnBoot {
		t.Errorf("expected RunMigrationOnBoot=false")
	}
}
