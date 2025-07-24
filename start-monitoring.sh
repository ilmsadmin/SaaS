#!/bin/bash

echo "📊 Starting Zplus SaaS Platform Monitoring Stack"
echo "==============================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Create monitoring network if it doesn't exist
print_info "Creating monitoring network..."
docker network create zplus-network 2>/dev/null || true

print_info "Starting monitoring stack..."
docker-compose -f docker-compose.monitoring.yml up -d

print_info "Waiting for services to be ready..."
sleep 30

echo ""
print_status "🎉 Monitoring stack is now running!"
echo "=================================="
echo ""
print_status "Prometheus (Metrics): http://localhost:9090"
print_status "Grafana (Dashboards): http://localhost:3002"
print_status "Jaeger (Tracing): http://localhost:16686"
print_status "Kibana (Logs): http://localhost:5601"
print_status "Elasticsearch: http://localhost:9200"
print_status "Node Exporter: http://localhost:9100"
print_status "cAdvisor: http://localhost:8081"
echo ""
print_info "Default Credentials:"
print_info "Grafana: admin / admin123"
echo ""
print_warning "Note: It may take a few minutes for all monitoring services to be fully ready."
echo ""
print_status "Happy monitoring! 📊"
