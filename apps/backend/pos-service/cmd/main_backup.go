package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"

	"zplus-saas/apps/backend/shared/config"
	"zplus-saas/apps/backend/shared/database"
	"zplus-saas/apps/backend/pos-service/internal/handlers"
	"zplus-saas/apps/backend/pos-service/internal/repositories"
	"zplus-saas/apps/backend/pos-service/internal/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database connection
	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Convert to sqlx.DB for better query handling
	sqlxDB := sqlx.NewDb(db, "postgres")

	// Initialize repositories
	productRepo := repositories.NewProductRepository(sqlxDB)
	categoryRepo := repositories.NewCategoryRepository(sqlxDB)
	orderRepo := repositories.NewOrderRepository(sqlxDB)
	inventoryRepo := repositories.NewInventoryRepository(sqlxDB)

	// Initialize services
	posService := services.NewPOSService(productRepo, categoryRepo, orderRepo, inventoryRepo)

	// Initialize handlers
	posHandler := handlers.NewPOSHandler(posService)

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
			"timestamp": "2025-07-23T00:00:00Z",
		})
	})

	// API routes
	api := app.Group("/api/v1")

	// POS endpoints
	pos := api.Group("/pos")

	// Health check for POS
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

	// Categories routes
	categories := pos.Group("/categories")
	categories.Post("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Create category endpoint - Implementation in progress",
		})
	})
	categories.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Get categories endpoint",
			"data":    []interface{}{},
		})
	})

	// Products routes
	products := pos.Group("/products")
	products.Post("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Create product endpoint - Implementation in progress",
		})
	})
	products.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Get products endpoint",
			"data":    []interface{}{},
			"meta": fiber.Map{
				"page":        1,
				"limit":       20,
				"total":       0,
				"total_pages": 1,
			},
		})
	})
	products.Get("/low-stock", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Get low stock products endpoint",
			"data":    []interface{}{},
		})
	})
	products.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Get product endpoint",
			"data": fiber.Map{
				"id":      id,
				"message": "Product details would be here",
			},
		})
	})

	// Orders routes
	orders := pos.Group("/orders")
	orders.Post("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Create order endpoint - Implementation in progress",
		})
	})
	orders.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Get orders endpoint",
			"data":    []interface{}{},
			"meta": fiber.Map{
				"page":        1,
				"limit":       20,
				"total":       0,
				"total_pages": 1,
			},
		})
	})
	orders.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Get order endpoint",
			"data": fiber.Map{
				"id":      id,
				"message": "Order details would be here",
			},
		})
	})

	// Analytics routes
	analytics := pos.Group("/analytics")
	analytics.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "POS Analytics",
			"data": fiber.Map{
				"total_sales":         0.0,
				"total_orders":        0,
				"average_order_value": 0.0,
				"total_products":      0,
				"low_stock_products":  0,
				"top_selling_products": []interface{}{},
				"sales_by_date":       []interface{}{},
			},
		})
	})
	analytics.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "POS Dashboard Analytics",
			"data": fiber.Map{
				"today_sales":    0.0,
				"today_orders":   0,
				"total_products": 0,
				"low_stock":      0,
				"recent_orders":  []interface{}{},
			},
		})
	})

	// Inventory routes
	inventory := pos.Group("/inventory")
	inventory.Get("/transactions", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Get inventory transactions",
			"data":    []interface{}{},
		})
	})
	inventory.Post("/adjust", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Inventory adjustment endpoint - Implementation in progress",
		})
	})

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084" // Default POS service port
	}

	// Start server
	fmt.Printf("🚀 POS Service starting on port %s\n", port)
	fmt.Println("📋 Available endpoints:")
	fmt.Println("  GET  /health - Service health check")
	fmt.Println("  GET  /api/v1/pos/health - POS health check")
	fmt.Println("  GET  /api/v1/pos/products - List products")
	fmt.Println("  POST /api/v1/pos/products - Create product")
	fmt.Println("  GET  /api/v1/pos/orders - List orders")
	fmt.Println("  POST /api/v1/pos/orders - Create order")
	fmt.Println("  GET  /api/v1/pos/categories - List categories")
	fmt.Println("  GET  /api/v1/pos/analytics - Get analytics")
	fmt.Println("  GET  /api/v1/pos/inventory/transactions - Get inventory")

	log.Fatal(app.Listen(":" + port))
}
