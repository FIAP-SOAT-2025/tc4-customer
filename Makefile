.PHONY: help build run test test-unit test-integration clean docker-build docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	go build -o bin/customer-service ./api

run: ## Run the application
	go run ./api/main.go

test: ## Run all tests
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

test-unit: ## Run unit tests only
	go test -v -race -short ./...

test-integration: ## Run integration tests
	go test -v -race -run Integration ./...

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.txt

docker-build: ## Build Docker image
	docker-compose build

docker-up: ## Start services with docker-compose
	docker-compose up -d

docker-down: ## Stop services
	docker-compose down

docker-logs: ## View docker-compose logs
	docker-compose logs -f

lint: ## Run linter
	golangci-lint run

deps: ## Download dependencies
	go mod download
	go mod tidy
