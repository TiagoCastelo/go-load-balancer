package balancer

import (
	"go-load-balancer/internal/backend"
	"testing"
)

func TestNewLoadBalancer(t *testing.T) {
	backends := []*backend.Backend{
		backend.NewBackend("http://localhost:8081"),
		backend.NewBackend("http://localhost:8082"),
	}

	lb := NewLoadBalancer("round-robin", backends)

	if lb == nil {
		t.Fatal("Expected load balancer to be created, got nil")
	}

	if lb.algorithm != "round-robin" {
		t.Errorf("Expected algorithm round-robin, got %s", lb.algorithm)
	}

	if len(lb.backends) != 2 {
		t.Errorf("Expected 2 backends, got %d", len(lb.backends))
	}
}

func TestRoundRobin(t *testing.T) {
	backends := []*backend.Backend{
		backend.NewBackend("http://localhost:8081"),
		backend.NewBackend("http://localhost:8082"),
		backend.NewBackend("http://localhost:8083"),
	}

	lb := NewLoadBalancer("round-robin", backends)

	// Test round-robin distribution
	first := lb.NextBackend()
	second := lb.NextBackend()
	third := lb.NextBackend()
	fourth := lb.NextBackend() // Should wrap back to first

	if first == nil || second == nil || third == nil || fourth == nil {
		t.Fatal("NextBackend returned nil")
	}

	// Verify we cycle through all backends
	if first.URL.String() == second.URL.String() {
		t.Error("Round robin should select different backends")
	}
}

func TestLeastConnections(t *testing.T) {
	backends := []*backend.Backend{
		backend.NewBackend("http://localhost:8081"),
		backend.NewBackend("http://localhost:8082"),
	}

	// Simulate connections on first backend
	backends[0].IncrementConnections()
	backends[0].IncrementConnections()

	lb := NewLoadBalancer("least-conn", backends)

	// Should select backend with fewer connections (backend2)
	selected := lb.NextBackend()
	if selected.URL.String() != "http://localhost:8082" {
		t.Errorf("Expected least-conn to select backend2, got %s", selected.URL.String())
	}
}
