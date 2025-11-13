#!/usr/bin/env sh
set -e
SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
ROOT_DIR=$(dirname "$SCRIPT_DIR")
cd "$ROOT_DIR"

echo "Running go vet..."
go vet ./...

echo "Running tests..."
go test ./...
