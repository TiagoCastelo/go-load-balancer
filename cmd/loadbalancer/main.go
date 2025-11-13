package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-load-balancer/internal/backend"
	"go-load-balancer/internal/balancer"
	"go-load-balancer/internal/config"
	"go-load-balancer/internal/dashboard"
)

func main() {
	cfg := config.Parse()

	if len(cfg.Backends) == 0 {
		log.Fatal("Please provide one or more backends to load balance")
	}

	// Build backend list
	var backends []*backend.Backend
	for _, b := range cfg.Backends {
		be := backend.NewBackend(b)
		if be != nil {
			backends = append(backends, be)
			log.Printf("Configured backend: %s", b)
		}
	}
	if len(backends) == 0 {
		log.Fatal("No valid backends after parsing inputs")
	}

	// Load balancer
	lb := balancer.NewLoadBalancer(cfg.Algorithm, backends)

	// Health checks
	go lb.HealthCheck(time.Duration(cfg.HealthCheckSeconds) * time.Second)

	// Create HTTP mux for multiple routes
	mux := http.NewServeMux()

	// Dashboard routes
	mux.HandleFunc("/dashboard", dashboard.DashboardHandler)
	mux.HandleFunc("/api/stats", dashboard.StatsAPIHandler)

	// Main load balancer route
	mux.HandleFunc("/", lb.ServeHTTP)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Load Balancer started at :%d", cfg.Port)
	log.Printf("Dashboard available at: http://localhost:%d/dashboard\n", cfg.Port)
	log.Printf("Algorithm: %s", cfg.Algorithm)
	log.Printf("Forwarding to backends: %v", cfg.Backends)

	// Graceful shutdown handling using context
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited cleanly")
}

// Keep flags visible in help for convenience (common pattern in CLI)
func init() {
	// Ensure default flag.Usage prints our usage from config package as well
	flag.Usage = config.Usage
}
