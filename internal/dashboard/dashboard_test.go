package dashboard

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordRequest(t *testing.T) {
	// Reset stats
	stats = &Stats{
		BackendRequests: make(map[string]int64),
		BackendStatus:   make(map[string]bool),
	}

	RecordRequest("http://backend1:8080")
	RecordRequest("http://backend1:8080")
	RecordRequest("http://backend2:8080")

	if stats.TotalRequests != 3 {
		t.Errorf("Expected 3 total requests, got %d", stats.TotalRequests)
	}

	if stats.BackendRequests["http://backend1:8080"] != 2 {
		t.Errorf("Expected 2 requests to backend1, got %d", stats.BackendRequests["http://backend1:8080"])
	}
}

func TestUpdateBackendStatus(t *testing.T) {
	stats = &Stats{
		BackendRequests: make(map[string]int64),
		BackendStatus:   make(map[string]bool),
	}

	UpdateBackendStatus("http://backend1:8080", true)
	UpdateBackendStatus("http://backend2:8080", false)

	if !stats.BackendStatus["http://backend1:8080"] {
		t.Error("Expected backend1 to be alive")
	}

	if stats.BackendStatus["http://backend2:8080"] {
		t.Error("Expected backend2 to be down")
	}
}

func TestDashboardHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/dashboard", nil)
	w := httptest.NewRecorder()

	DashboardHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html" {
		t.Errorf("Expected Content-Type text/html, got %s", contentType)
	}
}

func TestStatsAPIHandler(t *testing.T) {
	stats = &Stats{
		BackendRequests: make(map[string]int64),
		BackendStatus:   make(map[string]bool),
	}

	RecordRequest("http://backend1:8080")

	req := httptest.NewRequest("GET", "/api/stats", nil)
	w := httptest.NewRecorder()

	StatsAPIHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	if result["total_requests"].(float64) != 1 {
		t.Error("Expected 1 total request in stats API")
	}
}
