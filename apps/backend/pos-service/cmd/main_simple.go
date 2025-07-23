package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"

	"zplus-saas/apps/backend/shared/config"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database connection
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Convert to sqlx.DB for better query handling
	sqlxDB := sqlx.NewDb(db, "postgres")

	// For now, just initialize without repositories to avoid compilation errors
	_ = sqlxDB

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
				"code":    code,
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-Tenant-ID,X-User-ID",
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success":   true,
			"message":   "POS Service is running",
			"service":   "pos-service",
			"version":   "1.0.0",
			"timestamp": time.Now(),
		})
	})

	// API routes
	api := app.Group("/api/v1")
	pos := api.Group("/pos")

	// POS Health check
	pos.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "POS Service is healthy",
			"service": "pos-service",
			"features": []string{
				"Product Catalog",
				"Order Management",
				"Inventory Tracking",
				"Sales Analytics",
				"Category Management",
			},
		})
	})

	// Basic stub endpoints for testing
	pos.Get("/products", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Get products endpoint",
			"data":    []interface{}{},
		})
	})

	pos.Post("/products", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Create product endpoint - Ready for database integration",
		})
	})

	pos.Get("/categories", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Get categories endpoint",
			"data":    []interface{}{},
		})
	})

	pos.Post("/categories", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Create category endpoint - Ready for database integration",
		})
	})

	pos.Get("/orders", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Get orders endpoint",
			"data":    []interface{}{},
		})
	})

	pos.Post("/orders", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Create order endpoint - Ready for database integration",
		})
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	fmt.Printf("🚀 POS Service starting on port %s\n", port)
	fmt.Println("📋 Available endpoints:")
	fmt.Println("  GET  /health - Service health check")
	fmt.Println("  GET  /api/v1/pos/health - POS health check")
	fmt.Println("  GET  /api/v1/pos/products - List products")
	fmt.Println("  POST /api/v1/pos/products - Create product")
	fmt.Println("  GET  /api/v1/pos/orders - List orders")
	fmt.Println("  POST /api/v1/pos/orders - Create order")
	fmt.Println("  GET  /api/v1/pos/categories - List categories")
	fmt.Println("  POST /api/v1/pos/categories - Create category")

	log.Fatal(app.Listen(":" + port))
}
