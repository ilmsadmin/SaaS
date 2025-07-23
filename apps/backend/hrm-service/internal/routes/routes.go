package routes

import (
	"github.com/gofiber/fiber/v2"

	"zplus-saas/apps/backend/hrm-service/internal/handlers"
)

func SetupEmployeeRoutes(app *fiber.App, handler *handlers.EmployeeHandler) {
	api := app.Group("/api/v1/employees")

	// Employee CRUD operations
	api.Post("/", handler.CreateEmployee)
	api.Get("/:id", handler.GetEmployee)
	api.Get("/", handler.GetAllEmployees)
	api.Put("/:id", handler.UpdateEmployee)
	api.Delete("/:id", handler.DeleteEmployee)

	// Additional employee operations
	api.Get("/search", handler.SearchEmployees)
	api.Get("/by-email", handler.GetEmployeeByEmail)

	// Statistics
	api.Get("/statistics/hrm", handler.GetHRMStatistics)
}

func SetupDepartmentRoutes(app *fiber.App, handler *handlers.DepartmentHandler) {
	api := app.Group("/api/v1/departments")

	// Department CRUD operations
	api.Post("/", handler.CreateDepartment)
	api.Get("/:id", handler.GetDepartment)
	api.Get("/", handler.GetAllDepartments)
	api.Put("/:id", handler.UpdateDepartment)
	api.Delete("/:id", handler.DeleteDepartment)

	// Additional department operations
	api.Get("/with-employee-count", handler.GetDepartmentsWithEmployeeCount)
}

func SetupLeaveRoutes(app *fiber.App, handler *handlers.LeaveHandler) {
	api := app.Group("/api/v1/leaves")

	// Leave CRUD operations
	api.Post("/", handler.CreateLeave)
	api.Get("/:id", handler.GetLeave)
	api.Get("/", handler.GetAllLeaves)
	api.Put("/:id", handler.UpdateLeave)
	api.Delete("/:id", handler.DeleteLeave)

	// Leave management operations
	api.Post("/:id/approve", handler.ApproveLeave)
	api.Post("/:id/reject", handler.RejectLeave)

	// Leave queries
	api.Get("/employee/:employee_id", handler.GetLeavesByEmployee)
	api.Get("/statistics/pending-count", handler.GetPendingLeavesCount)
	api.Get("/balance", handler.GetLeaveBalance)
}

func SetupPerformanceRoutes(app *fiber.App, handler *handlers.PerformanceHandler) {
	api := app.Group("/api/v1/performance")

	// Performance CRUD operations
	api.Post("/", handler.CreatePerformance)
	api.Get("/:id", handler.GetPerformance)
	api.Get("/", handler.GetAllPerformance)
	api.Put("/:id", handler.UpdatePerformance)
	api.Delete("/:id", handler.DeletePerformance)

	// Performance workflow operations
	api.Post("/:id/submit", handler.SubmitPerformance)
	api.Post("/:id/complete", handler.CompletePerformance)

	// Performance queries
	api.Get("/employee/:employee_id", handler.GetPerformanceByEmployee)
	api.Get("/statistics/average-rating", handler.GetAveragePerformanceRating)
	api.Get("/statistics/by-department", handler.GetPerformanceStatsByDepartment)
}
