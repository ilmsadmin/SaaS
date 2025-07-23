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
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

// CheckinRecord represents a checkin/checkout record
type CheckinRecord struct {
	ID            int       `json:"id" db:"id"`
	TenantID      string    `json:"tenant_id" db:"tenant_id"`
	EmployeeID    int       `json:"employee_id" db:"employee_id"`
	EmployeeName  string    `json:"employee_name" db:"employee_name"`
	CheckinType   string    `json:"checkin_type" db:"checkin_type"`
	Timestamp     time.Time `json:"timestamp" db:"timestamp"`
	Location      string    `json:"location" db:"location"`
	Latitude      *float64  `json:"latitude,omitempty" db:"latitude"`
	Longitude     *float64  `json:"longitude,omitempty" db:"longitude"`
	IPAddress     string    `json:"ip_address" db:"ip_address"`
	DeviceInfo    string    `json:"device_info" db:"device_info"`
	Photo         string    `json:"photo,omitempty" db:"photo"`
	Notes         string    `json:"notes,omitempty" db:"notes"`
	Status        string    `json:"status" db:"status"`
	ApprovedBy    *int      `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt    *time.Time `json:"approved_at,omitempty" db:"approved_at"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// CheckinRequest represents a checkin request
type CheckinRequest struct {
	EmployeeID  int      `json:"employee_id"`
	CheckinType string   `json:"checkin_type"`
	Location    string   `json:"location,omitempty"`
	Latitude    *float64 `json:"latitude,omitempty"`
	Longitude   *float64 `json:"longitude,omitempty"`
	Photo       string   `json:"photo,omitempty"`
	Notes       string   `json:"notes,omitempty"`
}

// CheckinHandler handles checkin-related requests
type CheckinHandler struct {
	db *sqlx.DB
}

func NewCheckinHandler(db *sqlx.DB) *CheckinHandler {
	return &CheckinHandler{db: db}
}

// CreateCheckin creates a new checkin record
func (h *CheckinHandler) CreateCheckin(c *fiber.Ctx) error {
	tenantID := c.Get("X-Tenant-ID", "demo-tenant")

	var req CheckinRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Validate request
	if req.EmployeeID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Employee ID is required",
		})
	}

	validTypes := []string{"checkin", "checkout", "break_start", "break_end"}
	isValidType := false
	for _, validType := range validTypes {
		if req.CheckinType == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid checkin type",
		})
	}

	// Create checkin record
	record := &CheckinRecord{
		TenantID:     tenantID,
		EmployeeID:   req.EmployeeID,
		EmployeeName: fmt.Sprintf("Employee_%d", req.EmployeeID),
		CheckinType:  req.CheckinType,
		Timestamp:    time.Now(),
		Location:     req.Location,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		IPAddress:    c.IP(),
		DeviceInfo:   c.Get("User-Agent"),
		Photo:        req.Photo,
		Notes:        req.Notes,
		Status:       "approved",
	}

	query := `
		INSERT INTO checkin_records (
			tenant_id, employee_id, employee_name, checkin_type, timestamp,
			location, latitude, longitude, ip_address, device_info, photo, notes, status
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		) RETURNING id, created_at, updated_at
	`

	err := h.db.QueryRow(query, record.TenantID, record.EmployeeID, record.EmployeeName,
		record.CheckinType, record.Timestamp, record.Location, record.Latitude,
		record.Longitude, record.IPAddress, record.DeviceInfo, record.Photo,
		record.Notes, record.Status).Scan(&record.ID, &record.CreatedAt, &record.UpdatedAt)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create checkin record",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Checkin recorded successfully",
		"data":    record,
	})
}

// GetCheckinRecords gets checkin records
func (h *CheckinHandler) GetCheckinRecords(c *fiber.Ctx) error {
	tenantID := c.Get("X-Tenant-ID", "demo-tenant")

	var records []CheckinRecord
	query := `
		SELECT id, tenant_id, employee_id, employee_name, checkin_type, timestamp,
			   location, latitude, longitude, ip_address, device_info, photo, notes,
			   status, approved_by, approved_at, created_at, updated_at
		FROM checkin_records
		WHERE tenant_id = $1
		ORDER BY timestamp DESC
		LIMIT 100
	`

	err := h.db.Select(&records, query, tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to get checkin records",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  records,
	})
}

// GetTodayCheckinRecords gets today's checkin records for an employee
func (h *CheckinHandler) GetTodayCheckinRecords(c *fiber.Ctx) error {
	tenantID := c.Get("X-Tenant-ID", "demo-tenant")

	employeeID, err := strconv.Atoi(c.Params("employee_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid employee ID",
		})
	}

	today := time.Now().Format("2006-01-02")
	var records []CheckinRecord
	query := `
		SELECT id, tenant_id, employee_id, employee_name, checkin_type, timestamp,
			   location, latitude, longitude, ip_address, device_info, photo, notes,
			   status, approved_by, approved_at, created_at, updated_at
		FROM checkin_records
		WHERE tenant_id = $1 AND employee_id = $2 AND DATE(timestamp) = $3
		ORDER BY timestamp ASC
	`

	err = h.db.Select(&records, query, tenantID, employeeID, today)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to get today's checkin records",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  records,
	})
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
		AppName: "Checkin Service v1.0.0",
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
	checkinHandler := NewCheckinHandler(db)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"service":   "checkin-service",
			"timestamp": time.Now(),
		})
	})

	// API routes
	api := app.Group("/api/v1/checkin")
	api.Post("/", checkinHandler.CreateCheckin)
	api.Get("/", checkinHandler.GetCheckinRecords)
	api.Get("/employee/:employee_id/today", checkinHandler.GetTodayCheckinRecords)

	// Start server
	port := os.Getenv("CHECKIN_SERVICE_PORT")
	if port == "" {
		port = "8086"
	}

	log.Printf("Checkin Service starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
