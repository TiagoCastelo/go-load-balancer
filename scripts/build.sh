#!/usr/bin/env sh
set -e
SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
ROOT_DIR=$(dirname "$SCRIPT_DIR")
cd "$ROOT_DIR"

echo "Building load balancer..."
go build -o dist/loadbalancer ./cmd/loadbalancer
echo "Build complete: dist/loadbalancer"
