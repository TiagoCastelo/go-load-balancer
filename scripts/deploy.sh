#!/usr/bin/env sh
set -e
SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
ROOT_DIR=$(dirname "$SCRIPT_DIR")
cd "$ROOT_DIR/deployments"

echo "Starting docker compose stack..."
docker compose up --build -d
echo "Stack is up. Access load balancer at http://localhost:8080"
