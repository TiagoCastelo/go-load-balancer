package config

import (
	"flag"
	"os"
	"testing"
)

func TestParseDefaults(t *testing.T) {
	// Clear any env vars
	os.Clearenv()

	// Create a new FlagSet for this test to avoid conflicts
	oldCommandLine := flag.CommandLine
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	defer func() { flag.CommandLine = oldCommandLine }()

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
	// Create a new FlagSet for this test to avoid conflicts
	oldCommandLine := flag.CommandLine
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	defer func() { flag.CommandLine = oldCommandLine }()

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

	if cfg.HealthCheckSeconds != 60 {
		t.Errorf("Expected health check 60s from env, got %d", cfg.HealthCheckSeconds)
	}
}

func TestParseBackendsParsing(t *testing.T) {
	// Create a new FlagSet for this test
	oldCommandLine := flag.CommandLine
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	defer func() { flag.CommandLine = oldCommandLine }()

	os.Clearenv()
	os.Args = []string{"cmd", "-backends=http://localhost:8081, http://localhost:8082 , http://localhost:8083"}

	cfg := Parse()

	if len(cfg.Backends) != 3 {
		t.Errorf("Expected 3 backends after trimming spaces, got %d", len(cfg.Backends))
	}

	expected := []string{"http://localhost:8081", "http://localhost:8082", "http://localhost:8083"}
	for i, backend := range cfg.Backends {
		if backend != expected[i] {
			t.Errorf("Expected backend %s, got %s", expected[i], backend)
		}
	}
}
