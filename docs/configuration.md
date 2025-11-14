Configuration

Flags (highest precedence):
- -backends  Comma-separated backend URLs (e.g., http://backend1:8080,http://backend2:8080)
- -port      Port to listen on inside the container/process (default 8080)
- -algorithm Load balancing algorithm: round-robin or least-conn (default round-robin)
- -health    Health check interval in seconds (default 30)

Environment variables (set defaults for flags; flags still win):
- LB_BACKENDS        Comma-separated backend URLs
- LB_PORT            Port to listen on
- LB_ALGORITHM       round-robin or least-conn
- LB_HEALTH_INTERVAL Health check interval in seconds