# 🚀 Go Load Balancer

This repository provides a production-ready structure for a Go HTTP reverse-proxy load balancer with a real-time monitoring dashboard.

## ✨ Features

* **High-Performance Reverse Proxy:** Forwards HTTP requests to multiple backend servers.
* **Balancing Algorithms:** Easily extendable, includes `round-robin` and `least-connections`.
* **Active Health Checks:** Automatically monitors backend health and removes dead servers from the rotation.
* **Real-time Dashboard:** A built-in UI (using Go `embed`) to monitor backend status and request counts.

## 🐳 Quick Start (Docker - Recommended)

This is the fastest way to get the full system (load balancer + test backends) running.

1. Navigate to the deployments directory:

   ```bash
   cd deployments
   ````
2. Build and run the containers in detached mode:
   ```bash
    docker compose up --build -d
   ````

3. Access the services in your browser:

    * Load Balancer: http://localhost:9080 (Refresh to see it hit different backends)
    
    * Dashboard UI: http://localhost:9080/dashboard
    
    * Stats API: http://localhost:9080/api/stats

⚙️ Running Locally (Without Docker)

You can also run the load balancer as a standalone binary.

1. Build the executable:

    ```bash
    go build -o dist/loadbalancer ./cmd/loadbalancer
    ````

2. Run the load balancer, providing your backend URLs:
    ```bash
    ./dist/loadbalancer -backends=http://127.0.0.1:8081,http://127.0.0.1:8082 -port=8080
    ````

🧪 Testing the Load Balancer

Here is a simple PowerShell script to send 100 requests to your load balancer. Run this and watch your /dashboard page update in real-time.
```powershell
for ($i=1; $i -le 100; $i++) {
curl http://localhost:9080
Start-Sleep -Milliseconds 100
}
```