package models

import (
	"time"
)

// CheckinRecord represents a checkin/checkout record
type CheckinRecord struct {
	ID           int        `json:"id" db:"id"`
	TenantID     string     `json:"tenant_id" db:"tenant_id"`
	EmployeeID   int        `json:"employee_id" db:"employee_id"`
	EmployeeName string     `json:"employee_name" db:"employee_name"`
	CheckinType  string     `json:"checkin_type" db:"checkin_type"` // "checkin", "checkout", "break_start", "break_end"
	Timestamp    time.Time  `json:"timestamp" db:"timestamp"`
	Location     string     `json:"location" db:"location"`
	Latitude     *float64   `json:"latitude,omitempty" db:"latitude"`
	Longitude    *float64   `json:"longitude,omitempty" db:"longitude"`
	IPAddress    string     `json:"ip_address" db:"ip_address"`
	DeviceInfo   string     `json:"device_info" db:"device_info"`
	Photo        string     `json:"photo,omitempty" db:"photo"`
	Notes        string     `json:"notes,omitempty" db:"notes"`
	Status       string     `json:"status" db:"status"` // "approved", "pending", "rejected"
	ApprovedBy   *int       `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt   *time.Time `json:"approved_at,omitempty" db:"approved_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// AttendancePolicy represents attendance policy rules
type AttendancePolicy struct {
	ID              int       `json:"id" db:"id"`
	TenantID        string    `json:"tenant_id" db:"tenant_id"`
	Name            string    `json:"name" db:"name"`
	WorkStartTime   string    `json:"work_start_time" db:"work_start_time"`
	WorkEndTime     string    `json:"work_end_time" db:"work_end_time"`
	BreakDuration   int       `json:"break_duration" db:"break_duration"`   // minutes
	LateThreshold   int       `json:"late_threshold" db:"late_threshold"`   // minutes
	EarlyThreshold  int       `json:"early_threshold" db:"early_threshold"` // minutes
	RequirePhoto    bool      `json:"require_photo" db:"require_photo"`
	RequireLocation bool      `json:"require_location" db:"require_location"`
	AllowedRadius   int       `json:"allowed_radius" db:"allowed_radius"` // meters
	WorkDays        string    `json:"work_days" db:"work_days"`           // JSON array of days
	IsActive        bool      `json:"is_active" db:"is_active"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// AttendanceSummary represents daily attendance summary
type AttendanceSummary struct {
	ID            int        `json:"id" db:"id"`
	TenantID      string     `json:"tenant_id" db:"tenant_id"`
	EmployeeID    int        `json:"employee_id" db:"employee_id"`
	EmployeeName  string     `json:"employee_name" db:"employee_name"`
	Date          time.Time  `json:"date" db:"date"`
	CheckinTime   *time.Time `json:"checkin_time,omitempty" db:"checkin_time"`
	CheckoutTime  *time.Time `json:"checkout_time,omitempty" db:"checkout_time"`
	WorkHours     float64    `json:"work_hours" db:"work_hours"`
	BreakHours    float64    `json:"break_hours" db:"break_hours"`
	OvertimeHours float64    `json:"overtime_hours" db:"overtime_hours"`
	Status        string     `json:"status" db:"status"` // "present", "absent", "late", "early_leave", "partial"
	Notes         string     `json:"notes,omitempty" db:"notes"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// CheckinRequest represents a checkin request
type CheckinRequest struct {
	EmployeeID  int      `json:"employee_id" validate:"required"`
	CheckinType string   `json:"checkin_type" validate:"required,oneof=checkin checkout break_start break_end"`
	Location    string   `json:"location,omitempty"`
	Latitude    *float64 `json:"latitude,omitempty"`
	Longitude   *float64 `json:"longitude,omitempty"`
	Photo       string   `json:"photo,omitempty"`
	Notes       string   `json:"notes,omitempty"`
}

// AttendanceStats represents attendance statistics
type AttendanceStats struct {
	TotalEmployees int     `json:"total_employees"`
	PresentToday   int     `json:"present_today"`
	AbsentToday    int     `json:"absent_today"`
	LateToday      int     `json:"late_today"`
	AvgWorkHours   float64 `json:"avg_work_hours"`
	AttendanceRate float64 `json:"attendance_rate"`
}

// Constants
const (
	CheckinTypeCheckin    = "checkin"
	CheckinTypeCheckout   = "checkout"
	CheckinTypeBreakStart = "break_start"
	CheckinTypeBreakEnd   = "break_end"

	StatusApproved = "approved"
	StatusPending  = "pending"
	StatusRejected = "rejected"

	AttendanceStatusPresent    = "present"
	AttendanceStatusAbsent     = "absent"
	AttendanceStatusLate       = "late"
	AttendanceStatusEarlyLeave = "early_leave"
	AttendanceStatusPartial    = "partial"
)
