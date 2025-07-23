package routes

import (
	"../handlers"
	"github.com/gofiber/fiber/v2"
)

func CheckinRoutes(app fiber.Router, handler *handlers.CheckinHandler) {
	checkin := app.Group("/checkin")

	// Checkin records
	checkin.Post("/", handler.CreateCheckin)
	checkin.Get("/", handler.GetCheckinRecords)
	checkin.Get("/:id", handler.GetCheckinRecordByID)
	checkin.Put("/:id", handler.UpdateCheckinRecord)
	checkin.Delete("/:id", handler.DeleteCheckinRecord)

	// Approval
	checkin.Post("/:id/approve", handler.ApproveCheckin)
	checkin.Post("/:id/reject", handler.RejectCheckin)

	// Employee specific
	checkin.Get("/employee/:employee_id/today", handler.GetTodayCheckinRecords)

	// Statistics
	checkin.Get("/stats", handler.GetAttendanceStats)
}
