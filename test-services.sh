#!/bin/bash

# Comprehensive Service Test Suite
# Tests all microservices endpoints

echo "🧪 Running Comprehensive Service Test Suite..."

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Test counters
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

# Function to run test
run_test() {
    local test_name="$1"
    local command="$2"
    local expected_pattern="$3"
    
    echo -e "${BLUE}🔍 Testing: $test_name${NC}"
    ((TESTS_RUN++))
    
    response=$(eval "$command" 2>/dev/null)
    
    if echo "$response" | grep -q "$expected_pattern"; then
        echo -e "${GREEN}✅ PASS: $test_name${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}❌ FAIL: $test_name${NC}"
        echo -e "${YELLOW}Expected pattern: $expected_pattern${NC}"
        echo -e "${YELLOW}Actual response: $response${NC}"
        ((TESTS_FAILED++))
    fi
    echo ""
}

echo -e "${YELLOW}🔥 Starting Service Tests...${NC}"
echo ""

# Health Check Tests
echo -e "${BLUE}=========================${NC}"
echo -e "${BLUE}🏥 HEALTH CHECK TESTS${NC}"
echo -e "${BLUE}=========================${NC}"

run_test "Checkin Service Health" \
    "curl -s http://localhost:8086/health" \
    "checkin-service"

run_test "Payment Service Health" \
    "curl -s http://localhost:8087/health" \
    "payment-service"

run_test "File Service Health" \
    "curl -s http://localhost:8088/health" \
    "file-service"

# Functional Tests
echo -e "${BLUE}=========================${NC}"
echo -e "${BLUE}⚙️  FUNCTIONAL TESTS${NC}"
echo -e "${BLUE}=========================${NC}"

# Checkin Service Tests
run_test "Checkin Service - Create Check-in" \
    "curl -s -X POST http://localhost:8086/api/v1/checkin -H 'Content-Type: application/json' -d '{\"employee_id\": 999, \"checkin_type\": \"checkin\", \"location\": \"Test Office\", \"notes\": \"Automated test\"}'" \
    "Checkin recorded successfully"

run_test "Checkin Service - Get Check-ins" \
    "curl -s http://localhost:8086/api/v1/checkin" \
    "error.*false"

# File Service Tests  
echo "Creating test file for upload..."
echo "This is a test file created by automated test suite $(date)" > /tmp/test_upload.txt

run_test "File Service - Upload File" \
    "curl -s -X POST -F 'files=@/tmp/test_upload.txt' -H 'X-Tenant-ID: e04d1d3f-45dd-4f16-aa31-ce98d467ed6f' -H 'X-User-ID: 4c8a2684-97d1-4a66-85ea-f9466ebf61d6' http://localhost:8088/api/v1/files/upload" \
    "Successfully uploaded"

run_test "File Service - List Files" \
    "curl -s -H 'X-Tenant-ID: e04d1d3f-45dd-4f16-aa31-ce98d467ed6f' http://localhost:8088/api/v1/files" \
    "success"

# Database Connectivity Tests
echo -e "${BLUE}=========================${NC}"
echo -e "${BLUE}💾 DATABASE TESTS${NC}"
echo -e "${BLUE}=========================${NC}"

run_test "PostgreSQL Connectivity" \
    "docker exec zplus-postgres psql -U postgres -d zplus_saas -c 'SELECT 1 as test;' -t" \
    "1"

run_test "Checkin Records Count" \
    "docker exec zplus-postgres psql -U postgres -d zplus_saas -c 'SELECT COUNT(*) FROM checkin_records;' -t" \
    "[0-9]"

run_test "Files Table Count" \
    "docker exec zplus-postgres psql -U postgres -d zplus_saas -c 'SELECT COUNT(*) FROM files;' -t" \
    "[0-9]"

# API Gateway Integration Tests (if running)
echo -e "${BLUE}=========================${NC}"
echo -e "${BLUE}🌐 API GATEWAY TESTS${NC}"
echo -e "${BLUE}=========================${NC}"

# Check if API Gateway is running
if curl -s http://localhost:8080/health > /dev/null; then
    run_test "API Gateway - Checkin Route" \
        "curl -s http://localhost:8080/api/checkin/health" \
        "checkin-service"
        
    run_test "API Gateway - File Route" \
        "curl -s http://localhost:8080/api/file/health" \
        "file-service"
else
    echo -e "${YELLOW}⚠️  API Gateway not running - skipping gateway tests${NC}"
    echo ""
fi

# Cleanup
rm -f /tmp/test_upload.txt

# Test Results Summary
echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}📊 TEST RESULTS SUMMARY${NC}"
echo -e "${BLUE}================================================${NC}"
echo -e "${GREEN}✅ Tests Passed: $TESTS_PASSED${NC}"
echo -e "${RED}❌ Tests Failed: $TESTS_FAILED${NC}"
echo -e "${YELLOW}📋 Total Tests: $TESTS_RUN${NC}"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 ALL TESTS PASSED! Services are working correctly.${NC}"
    exit 0
else
    echo -e "${RED}⚠️  Some tests failed. Please check the service logs.${NC}"
    exit 1
fi
