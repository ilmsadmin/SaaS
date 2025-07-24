# Build and Development Scripts for Zplus SaaS Platform
.PHONY: help build dev test clean docker-build docker-up docker-down migrate seed lint format install start stop restart logs status health setup

# Default target
help: ## Show this help message
	@echo "🚀 Zplus SaaS Platform - Available Commands"  
	@echo "============================================"
	@echo 'Usage: make <target>'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Quick setup for new developers
setup: install build start ## Complete setup for new developers
	@echo "🎉 Setup completed! Platform is ready to use."
	@echo "Main App: http://localhost:3000"
	@echo "Admin: http://localhost:3001"

# Install dependencies
install: ## Install dependencies for all services
	@echo "📦 Installing dependencies..."
	@cd apps/frontend/web && npm install
	@cd apps/frontend/admin && npm install  
	@go mod download
	@echo "✅ Dependencies installed!"

# Start all services
start: ## Start all services with Docker
	@echo "🚀 Starting all services..."
	@./start-all.sh

# Stop all services
stop: ## Stop all services
	@echo "🛑 Stopping all services..."
	@docker-compose down
	@echo "✅ All services stopped!"

# Restart all services
restart: stop start ## Restart all services

# View logs
logs: ## View logs from all services
	@docker-compose logs -f

# Check service status
status: ## Check status of all services
	@echo "📊 Service Status:"
	@docker-compose ps

# Health check
health: ## Check health of all services
	@echo "🏥 Checking service health..."
	@curl -s http://localhost:8080/health && echo "API Gateway: ✅" || echo "API Gateway: ❌"
	@curl -s http://localhost:8081/health && echo "Auth Service: ✅" || echo "Auth Service: ❌"  
	@curl -s http://localhost:8082/health && echo "Tenant Service: ✅" || echo "Tenant Service: ❌"
	@curl -s http://localhost:8083/health && echo "CRM Service: ✅" || echo "CRM Service: ❌"
	@curl -s http://localhost:8084/health && echo "LMS Service: ✅" || echo "LMS Service: ❌"
	@curl -s http://localhost:8085/health && echo "POS Service: ✅" || echo "POS Service: ❌"
	@curl -s http://localhost:8086/health && echo "Checkin Service: ✅" || echo "Checkin Service: ❌"
	@curl -s http://localhost:8087/health && echo "Payment Service: ✅" || echo "Payment Service: ❌"
	@curl -s http://localhost:8088/health && echo "File Service: ✅" || echo "File Service: ❌"
	@curl -s http://localhost:8089/health && echo "HRM Service: ✅" || echo "HRM Service: ❌"

# Development modes
dev-web: ## Start web frontend in development mode
	@echo "🔥 Starting web frontend in dev mode..."
	@cd apps/frontend/web && npm run dev

dev-admin: ## Start admin frontend in development mode
	@echo "🔥 Starting admin frontend in dev mode..."
	@cd apps/frontend/admin && npm run dev

# Building
build: ## Build all services
	@echo "Building backend services..."
	$(MAKE) build-backend
	@echo "Building frontend applications..."
	$(MAKE) build-frontend

build-backend: ## Build all backend services
	@echo "Building API Gateway..."
	cd apps/backend/api-gateway && go build -o bin/api-gateway cmd/main.go
	@echo "Building Auth Service..."
	cd apps/backend/auth-service && go build -o bin/auth-service cmd/main.go
	@echo "Building CRM Service..."
	cd apps/backend/crm-service && go build -o bin/crm-service cmd/main.go
	@echo "Building LMS Service..."
	cd apps/backend/lms-service && go build -o bin/lms-service cmd/main.go
	@echo "Building POS Service..."
	cd apps/backend/pos-service && go build -o bin/pos-service cmd/main.go
	@echo "Building HRM Service..."
	cd apps/backend/hrm-service && go build -o bin/hrm-service cmd/main.go
	@echo "Building Checkin Service..."
	cd apps/backend/checkin-service && go build -o bin/checkin-service cmd/main.go
	@echo "Building File Service..."
	cd apps/backend/file-service && go build -o bin/file-service cmd/main.go
	@echo "Building Payment Service..."
	cd apps/backend/payment-service && go build -o bin/payment-service cmd/main.go

build-frontend: ## Build all frontend applications
	@echo "Building Web App..."
	cd apps/frontend/web && npm run build
	@echo "Building Mobile App..."
	cd apps/frontend/mobile && npm run build
	@echo "Building Admin App..."
	cd apps/frontend/admin && npm run build

# Docker
docker-build: ## Build all Docker images
	@echo "Building Docker images..."
	docker-compose build

docker-up: ## Start all services with Docker
	@echo "Starting all services with Docker..."
	docker-compose up -d
	@echo "Services started! Check http://localhost:3000"

docker-down: ## Stop all Docker services
	@echo "Stopping all Docker services..."
	docker-compose down

docker-logs: ## Show Docker logs
	docker-compose logs -f

docker-clean: ## Clean Docker containers and images
	@echo "Cleaning Docker containers and images..."
	docker-compose down --volumes --remove-orphans
	docker system prune -af

# Database
migrate: ## Run database migrations
	@echo "Running database migrations..."
	cd apps/backend/auth-service && go run cmd/migrate/main.go up
	cd apps/backend/crm-service && go run cmd/migrate/main.go up
	cd apps/backend/lms-service && go run cmd/migrate/main.go up
	cd apps/backend/pos-service && go run cmd/migrate/main.go up
	cd apps/backend/hrm-service && go run cmd/migrate/main.go up
	cd apps/backend/checkin-service && go run cmd/migrate/main.go up
	cd apps/backend/file-service && go run cmd/migrate/main.go up
	cd apps/backend/payment-service && go run cmd/migrate/main.go up

migrate-down: ## Rollback database migrations
	@echo "Rolling back database migrations..."
	cd apps/backend/auth-service && go run cmd/migrate/main.go down
	cd apps/backend/crm-service && go run cmd/migrate/main.go down
	cd apps/backend/lms-service && go run cmd/migrate/main.go down
	cd apps/backend/pos-service && go run cmd/migrate/main.go down
	cd apps/backend/hrm-service && go run cmd/migrate/main.go down
	cd apps/backend/checkin-service && go run cmd/migrate/main.go down
	cd apps/backend/file-service && go run cmd/migrate/main.go down
	cd apps/backend/payment-service && go run cmd/migrate/main.go down

seed: ## Seed database with sample data
	@echo "Seeding database with sample data..."
	go run scripts/seed/main.go

# Testing
test: ## Run all tests
	@echo "Running backend tests..."
	$(MAKE) test-backend
	@echo "Running frontend tests..."
	$(MAKE) test-frontend

test-backend: ## Run backend tests
	@echo "Running backend unit tests..."
	cd apps/backend/api-gateway && go test ./...
	cd apps/backend/auth-service && go test ./...
	cd apps/backend/crm-service && go test ./...
	cd apps/backend/lms-service && go test ./...
	cd apps/backend/pos-service && go test ./...
	cd apps/backend/hrm-service && go test ./...
	cd apps/backend/checkin-service && go test ./...
	cd apps/backend/file-service && go test ./...
	cd apps/backend/payment-service && go test ./...

test-frontend: ## Run frontend tests
	@echo "Running frontend tests..."
	cd apps/frontend/web && npm run test
	cd apps/frontend/mobile && npm run test
	cd apps/frontend/admin && npm run test

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	cd tests/integration && go test ./...

test-e2e: ## Run end-to-end tests
	@echo "Running E2E tests with Playwright..."
	cd tests/e2e/playwright && npm run test
	@echo "Running E2E tests with Cypress..."
	cd tests/e2e/cypress && npm run test

test-load: ## Run load tests
	@echo "Running load tests..."
	cd tests/load/k6 && k6 run api-load-test.js
	cd tests/load/artillery && artillery run load-test.yml

test-coverage: ## Generate test coverage reports
	@echo "Generating backend coverage..."
	cd apps/backend && go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Generating frontend coverage..."
	cd apps/frontend/web && npm run test:coverage

# Code Quality
lint: ## Run linters
	@echo "Running Go linter..."
	golangci-lint run ./...
	@echo "Running frontend linters..."
	cd apps/frontend/web && npm run lint
	cd apps/frontend/mobile && npm run lint
	cd apps/frontend/admin && npm run lint

format: ## Format code
	@echo "Formatting Go code..."
	gofmt -s -w .
	@echo "Formatting frontend code..."
	cd apps/frontend/web && npm run format
	cd apps/frontend/mobile && npm run format
	cd apps/frontend/admin && npm run format

# Security
security-scan: ## Run security scans
	@echo "Running security scans..."
	@echo "Scanning Go code..."
	gosec ./...
	@echo "Scanning dependencies..."
	go list -json -deps ./... | nancy sleuth
	@echo "Scanning frontend dependencies..."
	cd apps/frontend/web && npm audit

# Monitoring
monitoring-up: ## Start monitoring stack
	@echo "Starting monitoring stack..."
	cd infra/monitoring && docker-compose up -d
	@echo "Monitoring stack started!"
	@echo "Prometheus: http://localhost:9090"
	@echo "Grafana: http://localhost:3001 (admin/admin)"
	@echo "Jaeger: http://localhost:16686"

monitoring-down: ## Stop monitoring stack
	@echo "Stopping monitoring stack..."
	cd infra/monitoring && docker-compose down

# Deployment
deploy-staging: ## Deploy to staging environment
	@echo "Deploying to staging..."
	kubectl config use-context staging
	kubectl apply -f infra/k8s/namespaces/staging.yaml
	kubectl apply -f infra/k8s/configmaps/ -n zplus-staging
	kubectl apply -f infra/k8s/secrets/ -n zplus-staging
	kubectl apply -f infra/k8s/deployments/ -n zplus-staging
	kubectl apply -f infra/k8s/services/ -n zplus-staging
	kubectl apply -f infra/k8s/ingress/ -n zplus-staging

deploy-production: ## Deploy to production environment
	@echo "Deploying to production..."
	kubectl config use-context production
	kubectl apply -f infra/k8s/namespaces/production.yaml
	kubectl apply -f infra/k8s/configmaps/ -n zplus-production
	kubectl apply -f infra/k8s/secrets/ -n zplus-production
	kubectl apply -f infra/k8s/deployments/ -n zplus-production
	kubectl apply -f infra/k8s/services/ -n zplus-production
	kubectl apply -f infra/k8s/ingress/ -n zplus-production

# Utilities
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	find . -name "bin" -type d -exec rm -rf {} +
	find . -name "*.log" -type f -delete
	find . -name ".DS_Store" -type f -delete
	cd apps/frontend/web && rm -rf .next
	cd apps/frontend/mobile && rm -rf dist
	cd apps/frontend/admin && rm -rf dist

logs: ## Show application logs
	@echo "Showing application logs..."
	docker-compose logs -f api-gateway auth-service

backup: ## Backup databases
	@echo "Backing up databases..."
	./scripts/backup.sh

restore: ## Restore databases from backup
	@echo "Restoring databases from backup..."
	./scripts/restore.sh

# Health checks
health: ## Check service health
	@echo "Checking service health..."
	curl -f http://localhost:8080/health || exit 1
	curl -f http://localhost:8081/health || exit 1
	curl -f http://localhost:3000/api/health || exit 1

# Code generation
codegen: ## Generate code
	@echo "Generating GraphQL code..."
	cd apps/frontend/web && npm run codegen
	@echo "Generating GORM models..."
	cd tools/codegen && go run gorm-gen/main.go
	@echo "Generating OpenAPI clients..."
	cd tools/codegen && go run openapi-gen/main.go
