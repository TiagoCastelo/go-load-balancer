Go Load Balancer

This repository provides a production-ready structure for a Go HTTP reverse-proxy load balancer with round-robin and least-connections algorithms.

Quick start:
- Build: go build -o dist/loadbalancer ./cmd/loadbalancer
- Run:   ./dist/loadbalancer -backends=http://127.0.0.1:8081,http://127.0.0.1:8082 -port=8080

Docker demo:
- cd deployments
- docker compose up --build -d
- Open http://localhost:8080
