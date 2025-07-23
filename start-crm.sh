#!/bin/bash

# Start CRM Service Script
# Description: Script to start the CRM Service for development

echo "🎯 Starting Zplus SaaS CRM Service..."

# Set working directory
cd "$(dirname "$0")/apps/backend/crm-service"

# Set environment variables for development
export PORT=8082
export DATABASE_URL="postgres://postgres:postgres123@localhost:5432/zplus_saas?sslmode=disable"
export REDIS_URL="redis://localhost:6379"
export MONGODB_URL="mongodb://admin:admin123@localhost:27017"
export JWT_SECRET="your-super-secret-jwt-key"
export ENVIRONMENT="development"
export DEBUG="true"
export CRM_SERVICE_PORT="8082"

echo "📡 Environment: $ENVIRONMENT"
echo "🔗 Database: $DATABASE_URL"
echo "🚀 Starting on port: $PORT"

# Run the service
go run cmd/main.go
