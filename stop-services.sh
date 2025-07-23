#!/bin/bash

# Zplus SaaS Platform - Service Stop Script
# This script stops all running microservices

echo "🛑 Stopping Zplus SaaS Platform Microservices..."

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Function to stop service by PID file
stop_service() {
    local service_name=$1
    local pid_file="logs/${service_name}.pid"
    
    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file")
        if kill -0 "$pid" 2>/dev/null; then
            kill "$pid"
            echo -e "${GREEN}✅ Stopped $service_name (PID: $pid)${NC}"
        else
            echo -e "${YELLOW}⚠️  $service_name process not found${NC}"
        fi
        rm -f "$pid_file"
    else
        echo -e "${YELLOW}⚠️  No PID file for $service_name${NC}"
    fi
}

# Function to stop service by port
stop_by_port() {
    local port=$1
    local service_name=$2
    
    local pids=$(lsof -ti:$port)
    if [ -n "$pids" ]; then
        echo -e "${GREEN}🔫 Killing processes on port $port ($service_name)${NC}"
        echo "$pids" | xargs kill -9
    else
        echo -e "${YELLOW}ℹ️  No processes running on port $port${NC}"
    fi
}

echo -e "${YELLOW}🔍 Stopping services by PID files...${NC}"
stop_service "Checkin-Service"
stop_service "Payment-Service" 
stop_service "File-Service"

echo -e "${YELLOW}🔍 Stopping any remaining processes by port...${NC}"
stop_by_port 8086 "Checkin-Service"
stop_by_port 8087 "Payment-Service"
stop_by_port 8088 "File-Service"

echo -e "${GREEN}✅ All services stopped!${NC}"

# Optionally stop Docker infrastructure
read -p "Do you want to stop Docker infrastructure as well? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}🐳 Stopping Docker containers...${NC}"
    docker-compose down
    echo -e "${GREEN}✅ Docker infrastructure stopped!${NC}"
fi
