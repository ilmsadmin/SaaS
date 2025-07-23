#!/bin/bash

# LMS Service Startup Script
echo "🚀 Starting LMS Service..."

# Set environment variables
export PORT=8085
export LMS_SERVICE_URL="http://localhost:8085"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== LMS Service Configuration ===${NC}"
echo -e "Port: ${GREEN}$PORT${NC}"
echo -e "Service URL: ${GREEN}$LMS_SERVICE_URL${NC}"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed. Please install Go first.${NC}"
    exit 1
fi

# Navigate to LMS service directory
cd "$(dirname "$0")/apps/backend/lms-service" || exit 1

echo -e "${YELLOW}📦 Installing dependencies...${NC}"
go mod tidy

echo -e "${YELLOW}🏗️  Building LMS service...${NC}"
go build -o lms-service ./cmd/main.go

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Build successful!${NC}"
    echo -e "${BLUE}🚀 Starting LMS Service on port $PORT...${NC}"
    echo -e "${BLUE}📚 Available endpoints:${NC}"
    echo -e "  • Health Check: ${GREEN}http://localhost:$PORT/health${NC}"
    echo -e "  • LMS Health: ${GREEN}http://localhost:$PORT/api/v1/lms/health${NC}"
    echo -e "  • Courses: ${GREEN}http://localhost:$PORT/api/v1/lms/courses${NC}"
    echo -e "  • Enrollments: ${GREEN}http://localhost:$PORT/api/v1/lms/enrollments${NC}"
    echo -e "  • Progress: ${GREEN}http://localhost:$PORT/api/v1/lms/progress${NC}"
    echo -e "  • Quizzes: ${GREEN}http://localhost:$PORT/api/v1/lms/quizzes${NC}"
    echo -e "  • Assignments: ${GREEN}http://localhost:$PORT/api/v1/lms/assignments${NC}"
    echo -e "  • Analytics: ${GREEN}http://localhost:$PORT/api/v1/lms/analytics${NC}"
    echo ""
    echo -e "${YELLOW}Press Ctrl+C to stop the service${NC}"
    echo ""
    
    ./lms-service
else
    echo -e "${RED}❌ Build failed!${NC}"
    exit 1
fi
