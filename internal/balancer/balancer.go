package balancer

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"go-load-balancer/internal/backend"
	"go-load-balancer/internal/dashboard"
)

// LoadBalancer manages backend servers
type LoadBalancer struct {
	backends  []*backend.Backend
	current   uint64
	algorithm string
}

// NewLoadBalancer creates a new load balancer
func NewLoadBalancer(algorithm string, backends []*backend.Backend) *LoadBalancer {
	return &LoadBalancer{
		backends:  backends,
		current:   0,
		algorithm: algorithm,
	}
}

// NextBackend returns the next available backend based on the algorithm
func (lb *LoadBalancer) NextBackend() *backend.Backend {
	switch lb.algorithm {
	case "least-conn":
		return lb.leastConnections()
	default:
		return lb.roundRobin()
	}
}

// roundRobin selects backends in round-robin fashion
func (lb *LoadBalancer) roundRobin() *backend.Backend {
	numBackends := len(lb.backends)
	if numBackends == 0 {
		return nil
	}
	for i := 0; i < numBackends; i++ {
		idx := atomic.AddUint64(&lb.current, 1) % uint64(numBackends)
		be := lb.backends[idx]
		if be.IsAlive() {
			return be
		}
	}
	return nil
}

// leastConnections selects backend with least connections
func (lb *LoadBalancer) leastConnections() *backend.Backend {
	var selected *backend.Backend
	minConnections := -1
	for _, be := range lb.backends {
		if !be.IsAlive() {
			continue
		}
		c := be.GetConnections()
		if minConnections == -1 || c < minConnections {
			minConnections = c
			selected = be
		}
	}
	return selected
}

// ServeHTTP handles incoming requests
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	be := lb.NextBackend()
	if be == nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		log.Println("No available backends")
		return
	}
	be.IncrementConnections()
	defer be.DecrementConnections()
	// Record stats for dashboard
	dashboard.RecordRequest(be.URL.String())
	log.Printf("Forwarding request to %s (connections: %d)", be.URL.String(), be.GetConnections())
	be.ReverseProxy.ServeHTTP(w, r)
}

// HealthCheck periodically checks backend health
func (lb *LoadBalancer) HealthCheck(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		log.Println("Starting health check...")
		for _, be := range lb.backends {
			alive := be.IsHealthy()
			be.SetAlive(alive)

			// Update dashboard status
			dashboard.UpdateBackendStatus(be.URL.String(), alive) // ADD THIS

			status := "up"
			if !alive {
				status = "down"
			}
			log.Printf("Backend %s is %s", be.URL.String(), status)
		}
	}
}
