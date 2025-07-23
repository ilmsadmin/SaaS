package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"zplus-saas/apps/backend/hrm-service/internal/handlers"
	"zplus-saas/apps/backend/hrm-service/internal/repositories"
	"zplus-saas/apps/backend/hrm-service/internal/routes"
	"zplus-saas/apps/backend/hrm-service/internal/services"
	"zplus-saas/apps/backend/shared/config"
	"zplus-saas/apps/backend/shared/database"
	"zplus-saas/apps/backend/shared/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	employeeRepo := repositories.NewEmployeeRepository(db)
	departmentRepo := repositories.NewDepartmentRepository(db)
	leaveRepo := repositories.NewLeaveRepository(db)
	performanceRepo := repositories.NewPerformanceRepository(db)

	// Initialize services
	employeeService := services.NewEmployeeService(employeeRepo)
	departmentService := services.NewDepartmentService(departmentRepo)
	leaveService := services.NewLeaveService(leaveRepo)
	performanceService := services.NewPerformanceService(performanceRepo)

	// Initialize handlers
	employeeHandler := handlers.NewEmployeeHandler(employeeService)
	departmentHandler := handlers.NewDepartmentHandler(departmentService)
	leaveHandler := handlers.NewLeaveHandler(leaveService)
	performanceHandler := handlers.NewPerformanceHandler(performanceService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return ctx.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:8080",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Tenant-ID",
		AllowCredentials: true,
	}))

	// Apply tenant middleware
	app.Use(middleware.TenantMiddleware())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "hrm-service",
			"version": "1.0.0",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Setup routes
	routes.SetupEmployeeRoutes(app, employeeHandler)
	routes.SetupDepartmentRoutes(app, departmentHandler)
	routes.SetupLeaveRoutes(app, leaveHandler)
	routes.SetupPerformanceRoutes(app, performanceHandler)

	// Start server
	port := os.Getenv("HRM_SERVICE_PORT")
	if port == "" {
		port = "8083"
	}

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatal("Failed to start HRM service:", err)
		}
	}()

	log.Printf("🚀 HRM Service started on port %s", port)
	log.Printf("📊 Health check: http://localhost:%s/health", port)

	<-c
	log.Println("🔄 Gracefully shutting down HRM service...")
	_ = app.ShutdownWithContext(context.Background())
	log.Println("✅ HRM service stopped")
}
