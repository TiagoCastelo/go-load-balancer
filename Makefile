APP_NAME=loadbalancer
BIN_DIR=dist

.PHONY: all build test run clean docker up down fmt vet

all: build

build:
	@echo Building $(APP_NAME)...
	go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)

test:
	go test ./...

run: build
	./$(BIN_DIR)/$(APP_NAME) -backends=http://127.0.0.1:8081,http://127.0.0.1:8082 -port=8080

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf $(BIN_DIR)

docker:
	docker build -f deployments/docker/Dockerfile.loadbalancer -t $(APP_NAME):latest .

up:
	cd deployments && docker compose up --build -d

down:
	cd deployments && docker compose down
