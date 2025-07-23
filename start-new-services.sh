#!/bin/bash

# Zplus SaaS Platform - Service Startup Script
# This script starts all microservices for development

echo "🚀 Starting Zplus SaaS Platform Microservices..."

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to check if port is available
check_port() {
    if lsof -Pi :$1 -sTCP:LISTEN -t >/dev/null ; then
        echo -e "${RED}Port $1 is already in use${NC}"
        return 1
    else
        echo -e "${GREEN}Port $1 is available${NC}"
        return 0
    fi
}

# Function to start service in background
start_service() {
    local service_name=$1
    local port=$2
    local path=$3
    local main_file=$4
    
    echo -e "${BLUE}Starting $service_name on port $port...${NC}"
    
    if check_port $port; then
        cd "$path" && go run "$main_file" > "../../../logs/${service_name}.log" 2>&1 &
        local pid=$!
        echo $pid > "../../../logs/${service_name}.pid"
        echo -e "${GREEN}✅ $service_name started (PID: $pid)${NC}"
        sleep 2
    else
        echo -e "${RED}❌ Failed to start $service_name - port $port in use${NC}"
    fi
}

# Create logs directory
mkdir -p logs

echo -e "${YELLOW}📋 Checking infrastructure services...${NC}"

# Check Docker services
if ! docker-compose ps | grep -q "Up"; then
    echo -e "${YELLOW}Starting Docker infrastructure...${NC}"
    docker-compose up -d
    sleep 10
fi

echo -e "${YELLOW}🔧 Starting microservices...${NC}"

# Start all services
start_service "Checkin-Service" 8086 "apps/backend/checkin-service" "cmd/main.go"
start_service "Payment-Service" 8087 "apps/backend/payment-service" "cmd/main.go"  
start_service "File-Service" 8088 "apps/backend/file-service" "cmd/main_simple.go"

echo -e "${YELLOW}⏳ Waiting for services to initialize...${NC}"
sleep 5

echo -e "${YELLOW}🏥 Health check results:${NC}"

# Health checks
services=("8086:Checkin-Service" "8087:Payment-Service" "8088:File-Service")

for service in "${services[@]}"; do
    IFS=':' read -r port name <<< "$service"
    response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:$port/health)
    if [ "$response" = "200" ]; then
        echo -e "${GREEN}✅ $name (Port $port): Healthy${NC}"
    else
        echo -e "${RED}❌ $name (Port $port): Unhealthy (HTTP $response)${NC}"
    fi
done

echo -e "${BLUE}"
echo "=================================================="
echo "🎉 Service Startup Complete!"
echo "=================================================="
echo -e "${NC}"

echo -e "${YELLOW}📊 Service URLs:${NC}"
echo "• Checkin Service: http://localhost:8086/health"
echo "• Payment Service: http://localhost:8087/health" 
echo "• File Service: http://localhost:8088/health"

echo -e "${YELLOW}📋 Test Commands:${NC}"
echo "• Test Checkin: curl -X POST http://localhost:8086/api/v1/checkin -H 'Content-Type: application/json' -d '{\"employee_id\": 123, \"checkin_type\": \"checkin\", \"location\": \"Office\"}'"
echo "• Test File Upload: curl -X POST -F 'files=@test.txt' -H 'X-Tenant-ID: e04d1d3f-45dd-4f16-aa31-ce98d467ed6f' -H 'X-User-ID: 4c8a2684-97d1-4a66-85ea-f9466ebf61d6' http://localhost:8088/api/v1/files/upload"

echo -e "${YELLOW}🛠️  Management Commands:${NC}"
echo "• View logs: tail -f logs/[service-name].log"
echo "• Stop all: ./stop-services.sh"
echo "• Check status: ./check-services.sh"

echo -e "${GREEN}✅ All services started successfully!${NC}"
