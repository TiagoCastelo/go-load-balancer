package backend

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewBackend(t *testing.T) {
	backend := NewBackend("http://localhost:8080")
	if backend == nil {
		t.Fatal("Expected backend to be created, got nil")
	}
	if backend.URL.String() != "http://localhost:8080" {
		t.Errorf("Expected URL http://localhost:8080, got %s", backend.URL.String())
	}
	if !backend.IsAlive() {
		t.Error("Expected new backend to be alive")
	}
}

func TestBackendConnections(t *testing.T) {
	backend := NewBackend("http://localhost:8080")

	if backend.GetConnections() != 0 {
		t.Errorf("Expected 0 connections, got %d", backend.GetConnections())
	}

	backend.IncrementConnections()
	if backend.GetConnections() != 1 {
		t.Errorf("Expected 1 connection, got %d", backend.GetConnections())
	}

	backend.DecrementConnections()
	if backend.GetConnections() != 0 {
		t.Errorf("Expected 0 connections after decrement, got %d", backend.GetConnections())
	}
}

func TestBackendSetAlive(t *testing.T) {
	backend := NewBackend("http://localhost:8080")

	backend.SetAlive(false)
	if backend.IsAlive() {
		t.Error("Expected backend to be not alive")
	}

	backend.SetAlive(true)
	if !backend.IsAlive() {
		t.Error("Expected backend to be alive")
	}
}

func TestBackendIsHealthy(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	backend := NewBackend(server.URL)

	if !backend.IsHealthy() {
		t.Error("Expected backend to be healthy")
	}
}
