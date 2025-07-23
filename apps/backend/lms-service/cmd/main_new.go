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
	_ "github.com/lib/pq" // PostgreSQL driver

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
			"message":   "LMS Service is running",
			"service":   "lms-service",
			"version":   "1.0.0",
			"timestamp": time.Now(),
		})
	})

	// API routes
	api := app.Group("/api/v1")
	lms := api.Group("/lms")

	// LMS Health check
	lms.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "LMS Service is healthy",
			"service": "lms-service",
			"features": []string{
				"Course Management",
				"Student Enrollment",
				"Progress Tracking",
				"Assessment System",
				"Learning Analytics",
			},
		})
	})

	// Basic stub endpoints for testing
	lms.Get("/courses", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Get courses endpoint",
			"data":    []interface{}{},
		})
	})

	lms.Post("/courses", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Create course endpoint - Ready for database integration",
		})
	})

	lms.Get("/enrollments", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Get enrollments endpoint",
			"data":    []interface{}{},
		})
	})

	lms.Post("/enrollments", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Create enrollment endpoint - Ready for database integration",
		})
	})

	lms.Get("/progress", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Get progress endpoint",
			"data":    []interface{}{},
		})
	})

	lms.Get("/analytics", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Get analytics endpoint",
			"data": map[string]interface{}{
				"total_courses":      0,
				"total_students":     0,
				"completion_rate":    0,
				"popular_courses":    []interface{}{},
				"recent_enrollments": []interface{}{},
				"progress_stats":     []interface{}{},
			},
		})
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	fmt.Printf("🚀 LMS Service starting on port %s\n", port)
	fmt.Println("📋 Available endpoints:")
	fmt.Println("  GET  /health - Service health check")
	fmt.Println("  GET  /api/v1/lms/health - LMS health check")
	fmt.Println("  GET  /api/v1/lms/courses - List courses")
	fmt.Println("  POST /api/v1/lms/courses - Create course")
	fmt.Println("  GET  /api/v1/lms/enrollments - List enrollments")
	fmt.Println("  POST /api/v1/lms/enrollments - Create enrollment")
	fmt.Println("  GET  /api/v1/lms/progress - Get progress")
	fmt.Println("  GET  /api/v1/lms/analytics - Get analytics")

	log.Fatal(app.Listen(":" + port))
}
