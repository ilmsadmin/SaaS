package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Payment represents a payment record
type Payment struct {
	ID            int        `json:"id" db:"id"`
	TenantID      string     `json:"tenant_id" db:"tenant_id"`
	OrderID       string     `json:"order_id" db:"order_id"`
	CustomerID    *int       `json:"customer_id,omitempty" db:"customer_id"`
	CustomerName  string     `json:"customer_name" db:"customer_name"`
	Amount        float64    `json:"amount" db:"amount"`
	Currency      string     `json:"currency" db:"currency"`
	PaymentMethod string     `json:"payment_method" db:"payment_method"`
	Status        string     `json:"status" db:"status"`
	Gateway       string     `json:"gateway" db:"gateway"`
	GatewayTxnID  *string    `json:"gateway_txn_id,omitempty" db:"gateway_txn_id"`
	Description   *string    `json:"description,omitempty" db:"description"`
	Metadata      *string    `json:"metadata,omitempty" db:"metadata"`
	ProcessedAt   *time.Time `json:"processed_at,omitempty" db:"processed_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// PaymentRequest represents a payment request
type PaymentRequest struct {
	OrderID       string  `json:"order_id"`
	CustomerID    *int    `json:"customer_id,omitempty"`
	CustomerName  string  `json:"customer_name"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	PaymentMethod string  `json:"payment_method"`
	Gateway       string  `json:"gateway"`
	Description   string  `json:"description,omitempty"`
	Metadata      string  `json:"metadata,omitempty"`
}

// PaymentStats represents payment statistics
type PaymentStats struct {
	TotalPayments   int     `json:"total_payments"`
	TotalAmount     float64 `json:"total_amount"`
	SuccessfulCount int     `json:"successful_count"`
	PendingCount    int     `json:"pending_count"`
	FailedCount     int     `json:"failed_count"`
	SuccessRate     float64 `json:"success_rate"`
}

// PaymentHandler handles payment-related requests
type PaymentHandler struct {
	db *sqlx.DB
}

func NewPaymentHandler(db *sqlx.DB) *PaymentHandler {
	return &PaymentHandler{db: db}
}

// CreatePayment creates a new payment
func (h *PaymentHandler) CreatePayment(c *fiber.Ctx) error {
	tenantID := c.Get("X-Tenant-ID", "demo-tenant")

	var req PaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Validate request
	if req.OrderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Order ID is required",
		})
	}

	if req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Amount must be greater than 0",
		})
	}

	if req.Currency == "" {
		req.Currency = "USD"
	}

	validPaymentMethods := []string{"credit_card", "debit_card", "bank_transfer", "digital_wallet", "cash"}
	isValidMethod := false
	for _, method := range validPaymentMethods {
		if req.PaymentMethod == method {
			isValidMethod = true
			break
		}
	}
	if !isValidMethod {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid payment method",
		})
	}

	// Handle metadata - ensure it's valid JSON
	metadata := req.Metadata
	if metadata == "" {
		metadata = "{}"
	}

	// Create payment record
	payment := &Payment{
		TenantID:      tenantID,
		OrderID:       req.OrderID,
		CustomerID:    req.CustomerID,
		CustomerName:  req.CustomerName,
		Amount:        req.Amount,
		Currency:      req.Currency,
		PaymentMethod: req.PaymentMethod,
		Status:        "pending",
		Gateway:       req.Gateway,
		Description:   &req.Description,
		Metadata:      &metadata,
	}

	query := `
		INSERT INTO payments (
			tenant_id, order_id, customer_id, customer_name, amount, currency,
			payment_method, status, gateway, description, metadata
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) RETURNING id, created_at, updated_at
	`

	err := h.db.QueryRow(query, payment.TenantID, payment.OrderID, payment.CustomerID,
		payment.CustomerName, payment.Amount, payment.Currency, payment.PaymentMethod,
		payment.Status, payment.Gateway, payment.Description, payment.Metadata).
		Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)

	if err != nil {
		log.Printf("Failed to create payment: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": fmt.Sprintf("Failed to create payment: %v", err),
		})
	}

	// Simulate payment processing
	go h.processPayment(payment.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Payment created successfully",
		"data":    payment,
	})
}

// GetPayments gets payments with filters
func (h *PaymentHandler) GetPayments(c *fiber.Ctx) error {
	tenantID := c.Get("X-Tenant-ID", "demo-tenant")

	var payments []Payment
	query := `
		SELECT id, tenant_id, order_id, customer_id, customer_name, amount, currency,
			   payment_method, status, gateway, gateway_txn_id, description, metadata,
			   processed_at, created_at, updated_at
		FROM payments
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT 100
	`

	err := h.db.Select(&payments, query, tenantID)
	if err != nil {
		log.Printf("Failed to get payments: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": fmt.Sprintf("Failed to get payments: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  payments,
	})
}

// GetPaymentByID gets a payment by ID
func (h *PaymentHandler) GetPaymentByID(c *fiber.Ctx) error {
	tenantID := c.Get("X-Tenant-ID", "demo-tenant")

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid payment ID",
		})
	}

	var payment Payment
	query := `
		SELECT id, tenant_id, order_id, customer_id, customer_name, amount, currency,
			   payment_method, status, gateway, gateway_txn_id, description, metadata,
			   processed_at, created_at, updated_at
		FROM payments
		WHERE tenant_id = $1 AND id = $2
	`

	err = h.db.Get(&payment, query, tenantID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Payment not found",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  payment,
	})
}

// GetPaymentStats gets payment statistics
func (h *PaymentHandler) GetPaymentStats(c *fiber.Ctx) error {
	tenantID := c.Get("X-Tenant-ID", "demo-tenant")

	query := `
		SELECT 
			COUNT(*) as total_payments,
			COALESCE(SUM(amount), 0) as total_amount,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as successful_count,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_count,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_count
		FROM payments
		WHERE tenant_id = $1
	`

	var stats PaymentStats
	err := h.db.Get(&stats, query, tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to get payment stats",
		})
	}

	// Calculate success rate
	if stats.TotalPayments > 0 {
		stats.SuccessRate = float64(stats.SuccessfulCount) / float64(stats.TotalPayments) * 100
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  stats,
	})
}

// processPayment simulates payment processing
func (h *PaymentHandler) processPayment(paymentID int) {
	// Simulate processing delay
	time.Sleep(3 * time.Second)

	// Simulate random success/failure (90% success rate)
	status := "completed"
	gatewayTxnID := fmt.Sprintf("TXN_%d_%d", paymentID, time.Now().Unix())

	if time.Now().Unix()%10 == 0 { // 10% failure rate
		status = "failed"
		gatewayTxnID = ""
	}

	// Update payment status
	query := `
		UPDATE payments 
		SET status = $1, gateway_txn_id = $2, processed_at = NOW(), updated_at = NOW()
		WHERE id = $3
	`

	h.db.Exec(query, status, gatewayTxnID, paymentID)
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Database connection
	dbHost := os.Getenv("DATABASE_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DATABASE_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DATABASE_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres123"
	}
	dbName := os.Getenv("DATABASE_NAME")
	if dbName == "" {
		dbName = "zplus_saas"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Payment Service v1.0.0",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Tenant-ID",
	}))

	// Initialize handler
	paymentHandler := NewPaymentHandler(db)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"service":   "payment-service",
			"timestamp": time.Now(),
		})
	})

	// API routes
	api := app.Group("/api/v1/payment")
	api.Post("/", paymentHandler.CreatePayment)
	api.Get("/", paymentHandler.GetPayments)
	api.Get("/stats", paymentHandler.GetPaymentStats)
	api.Get("/:id", paymentHandler.GetPaymentByID)

	// Start server
	port := os.Getenv("PAYMENT_SERVICE_PORT")
	if port == "" {
		port = "8087"
	}

	log.Printf("Payment Service starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
