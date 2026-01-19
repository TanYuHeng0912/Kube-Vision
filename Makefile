.PHONY: help build run dev-backend dev-frontend test clean docker-build docker-run

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build-backend: ## Build backend
	cd backend && go build -o kubevision ./cmd/server

build-frontend: ## Build frontend
	cd frontend && npm install && npm run build

build: build-backend build-frontend ## Build both backend and frontend

run-backend: ## Run backend server
	cd backend && go run ./cmd/server

run-frontend: ## Run frontend dev server
	cd frontend && npm run dev

dev: ## Run both backend and frontend in development mode
	@echo "Starting development servers..."
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:5173"
	@make -j2 run-backend run-frontend

test-backend: ## Run backend tests
	cd backend && go test ./...

test-frontend: ## Run frontend tests
	cd frontend && npm test

test: test-backend test-frontend ## Run all tests

lint-backend: ## Lint backend code
	cd backend && golangci-lint run || true

lint-frontend: ## Lint frontend code
	cd frontend && npm run lint

lint: lint-backend lint-frontend ## Lint all code

docker-build: ## Build Docker image
	docker build -f docker/Dockerfile -t kubevision:latest .

docker-run: ## Run Docker container
	docker-compose -f docker/docker-compose.yml up

clean: ## Clean build artifacts
	rm -f backend/kubevision
	rm -rf frontend/dist
	rm -rf frontend/node_modules

install-backend: ## Install backend dependencies
	cd backend && go mod download

install-frontend: ## Install frontend dependencies
	cd frontend && npm install

install: install-backend install-frontend ## Install all dependencies

