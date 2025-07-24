#!/bin/bash

# Zplus SaaS Development Startup Script
echo "🚀 Starting Zplus SaaS Development Environment..."

# Set working directory
cd /Users/toan/Documents/project/SaaS

# Start infrastructure services
echo "📦 Starting infrastructure services..."
docker-compose up -d postgres redis minio

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 5

# Set environment variables
export DATABASE_URL=postgres://postgres:postgres123@localhost:5432/zplus_saas?sslmode=disable
export TENANT_SERVICE_PORT=8089
export AUTH_SERVICE_PORT=8081
export API_GATEWAY_PORT=8080
export APP_NAME="Zplus SaaS"
export REDIS_URL=redis://localhost:6379
export CORS_ALLOW_ORIGINS="http://localhost:3000,http://localhost:8080"
export CORS_ALLOW_METHODS="GET,POST,PUT,DELETE,OPTIONS"
export CORS_ALLOW_HEADERS="Origin,Content-Type,Accept,Authorization"

echo "🔧 Environment variables set"

# Function to start service in background
start_service() {
    SERVICE_NAME=$1
    SERVICE_PATH=$2
    SERVICE_PORT=$3
    
    echo "🔄 Starting $SERVICE_NAME on port $SERVICE_PORT..."
    cd $SERVICE_PATH
    go run cmd/main.go > /tmp/${SERVICE_NAME}.log 2>&1 &
    echo $! > /tmp/${SERVICE_NAME}.pid
    cd /Users/toan/Documents/project/SaaS
}

# Start backend services
start_service "api-gateway" "apps/backend/api-gateway" $API_GATEWAY_PORT
start_service "auth-service" "apps/backend/auth-service" $AUTH_SERVICE_PORT
start_service "tenant-service" "apps/backend/tenant-service" $TENANT_SERVICE_PORT

# Start frontend
echo "🎨 Starting web frontend on port 3000..."
cd apps/frontend/web
npm run dev > /tmp/frontend-web.log 2>&1 &
echo $! > /tmp/frontend-web.pid
cd /Users/toan/Documents/project/SaaS

# Start admin frontend
echo "🎨 Starting admin frontend on port 3001..."
cd apps/frontend/admin
npm run dev > /tmp/frontend-admin.log 2>&1 &
echo $! > /tmp/frontend-admin.pid
cd /Users/toan/Documents/project/SaaS

echo ""
echo "✅ Zplus SaaS Development Environment Started!"
echo ""
echo "🌐 Services:"
echo "   - Web Frontend: http://localhost:3000"
echo "   - Admin Panel:  http://localhost:3001"
echo "   - API Gateway:  http://localhost:8080"
echo "   - Auth Service: http://localhost:8081" 
echo "   - Tenant Service: http://localhost:8089"
echo ""
echo "📊 Infrastructure:"
echo "   - PostgreSQL:   localhost:5432"
echo "   - Redis:        localhost:6379"
echo "   - MinIO:        http://localhost:9001"
echo ""
echo "📝 Logs:"
echo "   - API Gateway:  tail -f /tmp/api-gateway.log"
echo "   - Auth Service: tail -f /tmp/auth-service.log"
echo "   - Tenant Service: tail -f /tmp/tenant-service.log"
echo "   - Web Frontend: tail -f /tmp/frontend-web.log"
echo "   - Admin Panel:  tail -f /tmp/frontend-admin.log"
echo ""
echo "🛑 To stop all services: ./stop-dev.sh"
echo ""
echo "👤 To create admin user: ./create-admin.sh"
echo "   (Run this after services are started)"
