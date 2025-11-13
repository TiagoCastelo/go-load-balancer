package backend

import (
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "sync"
    "time"
)

// Backend represents a backend server
type Backend struct {
    URL          *url.URL
    Alive        bool
    mux          sync.RWMutex
    ReverseProxy *httputil.ReverseProxy
    Connections  int
}

// NewBackend creates a new backend instance
func NewBackend(urlStr string) *Backend {
    u, err := url.Parse(urlStr)
    if err != nil {
        log.Printf("Error parsing URL %s: %v", urlStr, err)
        return nil
    }

    return &Backend{
        URL:          u,
        Alive:        true,
        ReverseProxy: httputil.NewSingleHostReverseProxy(u),
        Connections:  0,
    }
}

// SetAlive sets the alive status of the backend
func (b *Backend) SetAlive(alive bool) {
    b.mux.Lock()
    b.Alive = alive
    b.mux.Unlock()
}

// IsAlive returns true if backend is alive
func (b *Backend) IsAlive() bool {
    b.mux.RLock()
    alive := b.Alive
    b.mux.RUnlock()
    return alive
}

// IncrementConnections increments the connection count
func (b *Backend) IncrementConnections() {
    b.mux.Lock()
    b.Connections++
    b.mux.Unlock()
}

// DecrementConnections decrements the connection count
func (b *Backend) DecrementConnections() {
    b.mux.Lock()
    b.Connections--
    if b.Connections < 0 {
        b.Connections = 0
    }
    b.mux.Unlock()
}

// GetConnections returns the current connection count
func (b *Backend) GetConnections() int {
    b.mux.RLock()
    connections := b.Connections
    b.mux.RUnlock()
    return connections
}

// IsHealthy checks if backend is healthy
func (b *Backend) IsHealthy() bool {
    timeout := 5 * time.Second
    client := &http.Client{Timeout: timeout}

    resp, err := client.Get(b.URL.String())
    if err != nil {
        log.Printf("Backend %s is down: %v", b.URL.String(), err)
        return false
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 200 && resp.StatusCode < 500 {
        return true
    }

    log.Printf("Backend %s returned status code %d", b.URL.String(), resp.StatusCode)
    return false
}
