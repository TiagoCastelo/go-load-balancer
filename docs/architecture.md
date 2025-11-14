Architecture

This project implements a simple HTTP reverse-proxy load balancer in Go with round-robin and least-connections algorithms.
Core components live under `internal/` (backend management, balancing, config, dashboard), and the entrypoint is under `cmd/loadbalancer`.
A web dashboard and JSON stats API are exposed by the same HTTP server.