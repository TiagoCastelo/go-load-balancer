Configuration

Flags:
- -backends: Comma-separated backend URLs (e.g., http://backend1:80,http://backend2:80)
- -port: Port to listen on (default 8080)
- -algorithm: round-robin or least-conn (default round-robin)
- -health: Health check interval seconds (default 30)

Environment variables (override defaults, flags take precedence):
- LB_BACKENDS
- LB_PORT
- LB_ALGORITHM
- LB_HEALTH_INTERVAL
