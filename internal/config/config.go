package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Backends           []string
	Port               int
	Algorithm          string
	HealthCheckSeconds int
}

// Parse reads flags and environment variables to produce configuration.
// Env vars (if set) override defaults but not explicit flags.
//
// Supported env vars:
// - LB_BACKENDS (comma separated URLs)
// - LB_PORT
// - LB_ALGORITHM (round-robin|least-conn)
// - LB_HEALTH_INTERVAL (seconds)
func Parse() *Config {
	var (
		backendsStr string
		port        int
		algorithm   string
		healthSec   int
	)

	// Defaults
	defaultBackends := getEnv("LB_BACKENDS", "")
	defaultPort := getEnvInt("LB_PORT", 8080)
	defaultAlgorithm := getEnv("LB_ALGORITHM", "round-robin")
	defaultHealth := getEnvInt("LB_HEALTH_INTERVAL", 30)

	flag.StringVar(&backendsStr, "backends", defaultBackends, "Load balanced backends, use commas to separate")
	flag.IntVar(&port, "port", defaultPort, "Port to serve")
	flag.StringVar(&algorithm, "algorithm", defaultAlgorithm, "Load balancing algorithm (round-robin, least-conn)")
	flag.IntVar(&healthSec, "health", defaultHealth, "Health check interval in seconds")
	flag.Parse()

	var backends []string
	for _, s := range strings.Split(backendsStr, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			backends = append(backends, s)
		}
	}

	return &Config{
		Backends:           backends,
		Port:               port,
		Algorithm:          algorithm,
		HealthCheckSeconds: healthSec,
	}
}

// Usage prints helpful usage text, can be assigned to flag.Usage in main
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: loadbalancer -backends=http://127.0.0.1:8081,http://127.0.0.1:8082 -port=8080 -algorithm=round-robin -health=30\n")
	fmt.Fprintln(os.Stderr, "Environment variables:")
	fmt.Fprintln(os.Stderr, "  LB_BACKENDS           Comma-separated backend URLs")
	fmt.Fprintln(os.Stderr, "  LB_PORT               Port to listen on")
	fmt.Fprintln(os.Stderr, "  LB_ALGORITHM          round-robin|least-conn")
	fmt.Fprintln(os.Stderr, "  LB_HEALTH_INTERVAL    Health check interval in seconds")
	flag.PrintDefaults()
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
