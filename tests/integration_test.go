package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
)

// IntegrationTestSuite provides common functionality for integration tests
type IntegrationTestSuite struct {
	suite.Suite
	app    *fiber.App
	client *http.Client
}

// SetupSuite runs once before all tests
func (suite *IntegrationTestSuite) SetupSuite() {
	// Initialize test app
	suite.app = fiber.New()
	suite.client = &http.Client{
		Timeout: time.Second * 10,
	}
}

// TearDownSuite runs once after all tests
func (suite *IntegrationTestSuite) TearDownSuite() {
	// Cleanup
}

// SetupTest runs before each test
func (suite *IntegrationTestSuite) SetupTest() {
	// Setup for each test
}

// TearDownTest runs after each test
func (suite *IntegrationTestSuite) TearDownTest() {
	// Cleanup after each test
}

// Helper method to make HTTP requests
func (suite *IntegrationTestSuite) makeRequest(method, url string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req := httptest.NewRequest(method, url, bytes.NewReader(reqBody))

	// Set default headers
	req.Header.Set("Content-Type", "application/json")

	// Set custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp := httptest.NewRecorder()
	suite.app.Test(req)

	return resp
}

// Helper method to parse JSON response
func (suite *IntegrationTestSuite) parseJSONResponse(resp *httptest.ResponseRecorder, target interface{}) error {
	return json.Unmarshal(resp.Body.Bytes(), target)
}

// AuthTestSuite for authentication service tests
type AuthTestSuite struct {
	IntegrationTestSuite
}

func (suite *AuthTestSuite) TestUserRegistration() {
	payload := map[string]interface{}{
		"email":     "test@example.com",
		"password":  "password123",
		"name":      "Test User",
		"tenant_id": "test-tenant",
	}

	resp := suite.makeRequest("POST", "/auth/register", payload, nil)

	suite.Equal(http.StatusCreated, resp.Code)

	var response map[string]interface{}
	err := suite.parseJSONResponse(resp, &response)
	suite.NoError(err)
	suite.Contains(response, "token")
	suite.Contains(response, "user")
}

func (suite *AuthTestSuite) TestUserLogin() {
	// First register a user
	registerPayload := map[string]interface{}{
		"email":     "login@example.com",
		"password":  "password123",
		"name":      "Login Test User",
		"tenant_id": "test-tenant",
	}
	suite.makeRequest("POST", "/auth/register", registerPayload, nil)

	// Then try to login
	loginPayload := map[string]interface{}{
		"email":    "login@example.com",
		"password": "password123",
	}

	resp := suite.makeRequest("POST", "/auth/login", loginPayload, nil)

	suite.Equal(http.StatusOK, resp.Code)

	var response map[string]interface{}
	err := suite.parseJSONResponse(resp, &response)
	suite.NoError(err)
	suite.Contains(response, "token")
	suite.Contains(response, "user")
}

func (suite *AuthTestSuite) TestInvalidLogin() {
	payload := map[string]interface{}{
		"email":    "nonexistent@example.com",
		"password": "wrongpassword",
	}

	resp := suite.makeRequest("POST", "/auth/login", payload, nil)

	suite.Equal(http.StatusUnauthorized, resp.Code)
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

// TenantTestSuite for tenant service tests
type TenantTestSuite struct {
	IntegrationTestSuite
}

func (suite *TenantTestSuite) TestCreateTenant() {
	payload := map[string]interface{}{
		"name":        "Test Company",
		"domain":      "testcompany",
		"email":       "admin@testcompany.com",
		"plan":        "basic",
		"description": "Test company description",
	}

	resp := suite.makeRequest("POST", "/tenant/create", payload, nil)

	suite.Equal(http.StatusCreated, resp.Code)

	var response map[string]interface{}
	err := suite.parseJSONResponse(resp, &response)
	suite.NoError(err)
	suite.Contains(response, "tenant")
	suite.Equal("Test Company", response["tenant"].(map[string]interface{})["name"])
}

func (suite *TenantTestSuite) TestGetTenants() {
	// Create a test tenant first
	payload := map[string]interface{}{
		"name":   "Get Test Company",
		"domain": "gettestcompany",
		"email":  "admin@gettestcompany.com",
		"plan":   "basic",
	}
	suite.makeRequest("POST", "/tenant/create", payload, nil)

	// Get all tenants
	resp := suite.makeRequest("GET", "/tenant/list", nil, nil)

	suite.Equal(http.StatusOK, resp.Code)

	var response map[string]interface{}
	err := suite.parseJSONResponse(resp, &response)
	suite.NoError(err)
	suite.Contains(response, "tenants")
}

func TestTenantTestSuite(t *testing.T) {
	suite.Run(t, new(TenantTestSuite))
}

// CRMTestSuite for CRM service tests
type CRMTestSuite struct {
	IntegrationTestSuite
}

func (suite *CRMTestSuite) TestCreateCustomer() {
	payload := map[string]interface{}{
		"name":      "John Doe",
		"email":     "john@example.com",
		"phone":     "+1234567890",
		"company":   "Acme Corp",
		"status":    "active",
		"tenant_id": "test-tenant",
	}

	resp := suite.makeRequest("POST", "/crm/customers", payload, nil)

	suite.Equal(http.StatusCreated, resp.Code)
}

func (suite *CRMTestSuite) TestCreateLead() {
	payload := map[string]interface{}{
		"name":      "Jane Smith",
		"email":     "jane@example.com",
		"phone":     "+1234567891",
		"company":   "Tech Startup",
		"source":    "website",
		"status":    "new",
		"score":     75,
		"tenant_id": "test-tenant",
	}

	resp := suite.makeRequest("POST", "/crm/leads", payload, nil)

	suite.Equal(http.StatusCreated, resp.Code)
}

func TestCRMTestSuite(t *testing.T) {
	suite.Run(t, new(CRMTestSuite))
}

// HRMTestSuite for HRM service tests
type HRMTestSuite struct {
	IntegrationTestSuite
}

func (suite *HRMTestSuite) TestCreateEmployee() {
	payload := map[string]interface{}{
		"name":          "Alice Johnson",
		"email":         "alice@company.com",
		"phone":         "+1234567892",
		"position":      "Software Developer",
		"department_id": "1",
		"salary":        75000,
		"hire_date":     "2024-01-15",
		"status":        "active",
		"tenant_id":     "test-tenant",
	}

	resp := suite.makeRequest("POST", "/hrm/employees", payload, nil)

	suite.Equal(http.StatusCreated, resp.Code)
}

func (suite *HRMTestSuite) TestCreateDepartment() {
	payload := map[string]interface{}{
		"name":        "Engineering",
		"description": "Software development team",
		"budget":      500000,
		"location":    "San Francisco",
		"tenant_id":   "test-tenant",
	}

	resp := suite.makeRequest("POST", "/hrm/departments", payload, nil)

	suite.Equal(http.StatusCreated, resp.Code)
}

func TestHRMTestSuite(t *testing.T) {
	suite.Run(t, new(HRMTestSuite))
}

// POSTestSuite for POS service tests
type POSTestSuite struct {
	IntegrationTestSuite
}

func (suite *POSTestSuite) TestCreateProduct() {
	payload := map[string]interface{}{
		"name":        "Test Product",
		"description": "A test product for unit testing",
		"price":       29.99,
		"category_id": "1",
		"sku":         "TEST-001",
		"stock":       100,
		"tenant_id":   "test-tenant",
	}

	resp := suite.makeRequest("POST", "/pos/products", payload, nil)

	suite.Equal(http.StatusCreated, resp.Code)
}

func (suite *POSTestSuite) TestCreateOrder() {
	// First create a product
	productPayload := map[string]interface{}{
		"name":        "Order Product",
		"price":       19.99,
		"category_id": "1",
		"sku":         "ORDER-001",
		"stock":       50,
		"tenant_id":   "test-tenant",
	}
	suite.makeRequest("POST", "/pos/products", productPayload, nil)

	// Then create an order
	orderPayload := map[string]interface{}{
		"customer_name": "Test Customer",
		"items": []map[string]interface{}{
			{
				"product_id": "1",
				"quantity":   2,
				"price":      19.99,
			},
		},
		"total":     39.98,
		"status":    "pending",
		"tenant_id": "test-tenant",
	}

	resp := suite.makeRequest("POST", "/pos/orders", orderPayload, nil)

	suite.Equal(http.StatusCreated, resp.Code)
}

func TestPOSTestSuite(t *testing.T) {
	suite.Run(t, new(POSTestSuite))
}

// Helper function to run all tests
func RunAllTests(t *testing.T) {
	fmt.Println("🧪 Running Zplus SaaS Platform Integration Tests")
	fmt.Println("===============================================")

	// Run all test suites
	suite.Run(t, new(AuthTestSuite))
	suite.Run(t, new(TenantTestSuite))
	suite.Run(t, new(CRMTestSuite))
	suite.Run(t, new(HRMTestSuite))
	suite.Run(t, new(POSTestSuite))

	fmt.Println("✅ All integration tests completed!")
}
