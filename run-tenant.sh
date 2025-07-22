#!/bin/bash

# Set working directory
cd /Users/toan/Documents/project/SaaS

# Export environment variables
export DATABASE_URL=postgres://postgres:postgres123@localhost:5432/zplus_saas?sslmode=disable
export TENANT_SERVICE_PORT=8089
export APP_NAME="Zplus SaaS"
export REDIS_URL=redis://localhost:6379
export CORS_ALLOW_ORIGINS="http://localhost:3000,http://localhost:8080"
export CORS_ALLOW_METHODS="GET,POST,PUT,DELETE,OPTIONS"
export CORS_ALLOW_HEADERS="Origin,Content-Type,Accept,Authorization"

echo "Starting Tenant Service..."
echo "DATABASE_URL: $DATABASE_URL"
echo "Port: $TENANT_SERVICE_PORT"

# Go to tenant service directory and run
cd apps/backend/tenant-service
go run cmd/main.go
