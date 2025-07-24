#!/bin/bash

echo "🚀 Starting Zplus SaaS Platform - All Services"
echo "=============================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker and try again."
    exit 1
fi

print_info "Starting infrastructure services..."
docker-compose up -d postgres mongodb redis minio

print_info "Waiting for databases to be ready..."
sleep 10

print_info "Starting backend services..."
docker-compose up -d api-gateway auth-service tenant-service crm-service lms-service pos-service hrm-service checkin-service payment-service file-service

print_info "Waiting for backend services to be ready..."
sleep 15

print_info "Starting frontend applications..."
docker-compose up -d web-app admin-app

print_info "All services are starting up..."
sleep 10

echo ""
echo "🎉 Zplus SaaS Platform is now running!"
echo "======================================"
echo ""
print_status "Main Web Application: http://localhost:3000"
print_status "Admin Dashboard: http://localhost:3001"
print_status "API Gateway: http://localhost:8080"
print_status "MinIO Console: http://localhost:9001"
echo ""
print_info "Service Status:"
print_status "✅ Auth Service (8081) - User authentication"
print_status "✅ Tenant Service (8082) - Multi-tenant management"
print_status "✅ CRM Service (8083) - Customer management"
print_status "✅ LMS Service (8084) - Learning management"
print_status "✅ POS Service (8085) - Point of sale"
print_status "✅ Checkin Service (8086) - Employee check-in"
print_status "✅ Payment Service (8087) - Payment processing"
print_status "✅ File Service (8088) - File management"
print_status "✅ HRM Service (8089) - Human resources"
echo ""
print_info "Infrastructure:"
print_status "✅ PostgreSQL (5432) - Main database"
print_status "✅ MongoDB (27017) - Document storage"
print_status "✅ Redis (6379) - Cache & sessions"
print_status "✅ MinIO (9000/9001) - File storage"
echo ""
print_warning "Note: It may take a few minutes for all services to be fully ready."
print_info "Use 'docker-compose ps' to check service status"
print_info "Use 'docker-compose logs [service-name]' to view logs"
print_info "Use './stop-dev.sh' to stop all services"
echo ""
print_status "Happy coding! 🚀"
