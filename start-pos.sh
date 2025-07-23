#!/bin/bash

# POS Service Startup Script
echo "🚀 Starting POS Service..."

# Set environment variables
export PORT=8084
export POS_SERVICE_URL="http://localhost:8084"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== POS Service Configuration ===${NC}"
echo -e "Port: ${GREEN}$PORT${NC}"
echo -e "Service URL: ${GREEN}$POS_SERVICE_URL${NC}"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed. Please install Go first.${NC}"
    exit 1
fi

# Navigate to POS service directory
cd "$(dirname "$0")/apps/backend/pos-service" || exit 1

echo -e "${YELLOW}📦 Installing dependencies...${NC}"
go mod tidy

echo -e "${YELLOW}🏗️  Building POS service...${NC}"
go build -o pos-service ./cmd/main.go

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Build successful!${NC}"
    echo -e "${BLUE}🚀 Starting POS Service on port $PORT...${NC}"
    echo -e "${BLUE}📋 Available endpoints:${NC}"
    echo -e "  • Health Check: ${GREEN}http://localhost:$PORT/health${NC}"
    echo -e "  • POS Health: ${GREEN}http://localhost:$PORT/api/v1/pos/health${NC}"
    echo -e "  • Products: ${GREEN}http://localhost:$PORT/api/v1/pos/products${NC}"
    echo -e "  • Orders: ${GREEN}http://localhost:$PORT/api/v1/pos/orders${NC}"
    echo -e "  • Categories: ${GREEN}http://localhost:$PORT/api/v1/pos/categories${NC}"
    echo -e "  • Analytics: ${GREEN}http://localhost:$PORT/api/v1/pos/analytics${NC}"
    echo ""
    echo -e "${YELLOW}Press Ctrl+C to stop the service${NC}"
    echo ""
    
    ./pos-service
else
    echo -e "${RED}❌ Build failed!${NC}"
    exit 1
fi
