package dashboard

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	_ "embed" // Required for the //go:embed directive
)

//go:embed dashboard.html
var dashboardHTML string

// Stats hold the state of the Load Balancer
type Stats struct {
	mu              sync.RWMutex
	TotalRequests   int64
	BackendRequests map[string]int64
	BackendStatus   map[string]bool
	Uptime          time.Time
}

// Initialises the stats struct globally
var stats = &Stats{
	BackendRequests: make(map[string]int64),
	BackendStatus:   make(map[string]bool),
	Uptime:          time.Now(),
}

// RecordRequest increments the counters in a thread-safe manner
func RecordRequest(backend string) {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.TotalRequests++
	stats.BackendRequests[backend]++
}

// UpdateBackendStatus updates whether the backend is Alive or Dead
func UpdateBackendStatus(backend string, isAlive bool) {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.BackendStatus[backend] = isAlive
}

// GetStats returns a safe copy of the data for the API
func GetStats() map[string]interface{} {
	stats.mu.RLock()
	defer stats.mu.RUnlock()

	// We make a deep copy of the maps.
	safeRequests := make(map[string]int64)
	for k, v := range stats.BackendRequests {
		safeRequests[k] = v
	}

	safeStatus := make(map[string]bool)
	for k, v := range stats.BackendStatus {
		safeStatus[k] = v
	}

	return map[string]interface{}{
		"total_requests":   stats.TotalRequests,
		"backend_requests": safeRequests,
		"backend_status":   safeStatus,
		"uptime":           time.Since(stats.Uptime).String(),
	}
}

// DashboardHandler serves the embedded HTML file
func DashboardHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	_, err := fmt.Fprint(w, dashboardHTML)
	if err != nil {
		// Log if we fail to write the HTML (e.g. the client disconnected)
		log.Printf("ERROR: Failed to write dashboard HTML: %v", err)
	}
}

// StatsAPIHandler serves the stat data as JSON for the JavaScript to consume
func StatsAPIHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// GetStats() is called *before* encoding
	statsData := GetStats()

	err := json.NewEncoder(w).Encode(statsData)
	if err != nil {
		// Log the error so you know it happened.
		// We can't reliably send an http.Error here because
		// headers may have already been written.
		log.Printf("ERROR: Failed to encode dashboard stats: %v", err)
	}
}
