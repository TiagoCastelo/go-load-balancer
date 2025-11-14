package config

import (
	"os"
	"testing"
)

func TestParseDefaults(t *testing.T) {
	// Clear any env vars
	os.Clearenv()

	// Reset flags for testing
	os.Args = []string{"cmd", "-backends=http://localhost:8081,http://localhost:8082"}

	cfg := Parse()

	if len(cfg.Backends) != 2 {
		t.Errorf("Expected 2 backends, got %d", len(cfg.Backends))
	}

	if cfg.Port != 8080 {
		t.Errorf("Expected default port 8080, got %d", cfg.Port)
	}

	if cfg.Algorithm != "round-robin" {
		t.Errorf("Expected default algorithm round-robin, got %s", cfg.Algorithm)
	}

	if cfg.HealthCheckSeconds != 30 {
		t.Errorf("Expected default health check 30s, got %d", cfg.HealthCheckSeconds)
	}
}

func TestParseEnvironmentVariables(t *testing.T) {
	os.Setenv("LB_BACKENDS", "http://localhost:9001")
	os.Setenv("LB_PORT", "9000")
	os.Setenv("LB_ALGORITHM", "least-conn")
	os.Setenv("LB_HEALTH_INTERVAL", "60")
	defer os.Clearenv()

	os.Args = []string{"cmd"}
	cfg := Parse()

	if len(cfg.Backends) != 1 {
		t.Errorf("Expected 1 backend from env, got %d", len(cfg.Backends))
	}

	if cfg.Port != 9000 {
		t.Errorf("Expected port 9000 from env, got %d", cfg.Port)
	}

	if cfg.Algorithm != "least-conn" {
		t.Errorf("Expected algorithm least-conn from env, got %s", cfg.Algorithm)
	}
}
