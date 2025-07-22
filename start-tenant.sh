#!/bin/bash

# Start Tenant Service Script
# Description: Script to start the Tenant Service for development

echo "🏢 Starting Zplus SaaS Tenant Service..."

# Set working directory
cd "$(dirname "$0")/apps/backend/tenant-service"

# Set environment variables for development
export PORT=8089
export DATABASE_URL="postgres://postgres:postgres123@localhost:5432/zplus_saas?sslmode=disable"
export REDIS_URL="redis://localhost:6379"
export MONGODB_URL="mongodb://admin:admin123@localhost:27017"
export JWT_SECRET="your-super-secret-jwt-key"
export ENVIRONMENT="development"
export DEBUG="true"
export TENANT_SERVICE_PORT="8089"

echo "📡 Environment: $ENVIRONMENT"
echo "🔗 Database: $DATABASE_URL"
echo "🚀 Starting on port: $PORT"

# Run the service
go run cmd/main.go
