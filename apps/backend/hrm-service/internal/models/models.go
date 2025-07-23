package models

import (
	"time"
)

// Department represents a department in the organization
type Department struct {
	ID          int       `json:"id" db:"id"`
	TenantID    string    `json:"tenant_id" db:"tenant_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	ManagerID   *int      `json:"manager_id" db:"manager_id"`
	Budget      float64   `json:"budget" db:"budget"`
	Location    string    `json:"location" db:"location"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Relations
	Manager   *Employee  `json:"manager,omitempty"`
	Employees []Employee `json:"employees,omitempty"`
}

// Employee represents an employee in the system
type Employee struct {
	ID             int        `json:"id" db:"id"`
	TenantID       string     `json:"tenant_id" db:"tenant_id"`
	EmployeeCode   string     `json:"employee_code" db:"employee_code"`
	FirstName      string     `json:"first_name" db:"first_name"`
	LastName       string     `json:"last_name" db:"last_name"`
	Email          string     `json:"email" db:"email"`
	Phone          string     `json:"phone" db:"phone"`
	DepartmentID   int        `json:"department_id" db:"department_id"`
	Position       string     `json:"position" db:"position"`
	HireDate       time.Time  `json:"hire_date" db:"hire_date"`
	Salary         float64    `json:"salary" db:"salary"`
	Status         string     `json:"status" db:"status"` // active, inactive, terminated
	ManagerID      *int       `json:"manager_id" db:"manager_id"`
	Address        string     `json:"address" db:"address"`
	DateOfBirth    *time.Time `json:"date_of_birth" db:"date_of_birth"`
	Gender         string     `json:"gender" db:"gender"`
	EmergencyName  string     `json:"emergency_name" db:"emergency_name"`
	EmergencyPhone string     `json:"emergency_phone" db:"emergency_phone"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`

	// Relations
	Department *Department `json:"department,omitempty"`
	Manager    *Employee   `json:"manager,omitempty"`
}

// EmployeeStatus constants
const (
	EmployeeStatusActive     = "active"
	EmployeeStatusInactive   = "inactive"
	EmployeeStatusTerminated = "terminated"
)

// Leave represents a leave request/record
type Leave struct {
	ID         int        `json:"id" db:"id"`
	TenantID   string     `json:"tenant_id" db:"tenant_id"`
	EmployeeID int        `json:"employee_id" db:"employee_id"`
	LeaveType  string     `json:"leave_type" db:"leave_type"`
	StartDate  time.Time  `json:"start_date" db:"start_date"`
	EndDate    time.Time  `json:"end_date" db:"end_date"`
	Days       int        `json:"days" db:"days"`
	Reason     string     `json:"reason" db:"reason"`
	Status     string     `json:"status" db:"status"` // pending, approved, rejected
	ApprovedBy *int       `json:"approved_by" db:"approved_by"`
	ApprovedAt *time.Time `json:"approved_at" db:"approved_at"`
	Comments   string     `json:"comments" db:"comments"`
	IsActive   bool       `json:"is_active" db:"is_active"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`

	// Relations
	Employee *Employee `json:"employee,omitempty"`
	Approver *Employee `json:"approver,omitempty"`
}

// Leave types and statuses
const (
	LeaveTypeAnnual    = "annual"
	LeaveTypeSick      = "sick"
	LeaveTypeMaternity = "maternity"
	LeaveTypePaternity = "paternity"
	LeaveTypePersonal  = "personal"
	LeaveTypeEmergency = "emergency"

	LeaveStatusPending  = "pending"
	LeaveStatusApproved = "approved"
	LeaveStatusRejected = "rejected"
)

// Performance represents performance evaluation records
type Performance struct {
	ID            int       `json:"id" db:"id"`
	TenantID      string    `json:"tenant_id" db:"tenant_id"`
	EmployeeID    int       `json:"employee_id" db:"employee_id"`
	ReviewerID    int       `json:"reviewer_id" db:"reviewer_id"`
	Period        string    `json:"period" db:"period"`                 // Q1-2025, H1-2025, 2025
	ReviewType    string    `json:"review_type" db:"review_type"`       // quarterly, annual, probation
	OverallRating float64   `json:"overall_rating" db:"overall_rating"` // 1-5 scale
	Goals         string    `json:"goals" db:"goals"`
	Achievements  string    `json:"achievements" db:"achievements"`
	Strengths     string    `json:"strengths" db:"strengths"`
	Areas         string    `json:"areas_for_improvement" db:"areas_for_improvement"`
	Comments      string    `json:"comments" db:"comments"`
	Status        string    `json:"status" db:"status"` // draft, submitted, completed
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`

	// Relations
	Employee *Employee `json:"employee,omitempty"`
	Reviewer *Employee `json:"reviewer,omitempty"`
}

// Performance review types and statuses
const (
	ReviewTypeQuarterly = "quarterly"
	ReviewTypeAnnual    = "annual"
	ReviewTypeProbation = "probation"

	PerformanceStatusDraft     = "draft"
	PerformanceStatusSubmitted = "submitted"
	PerformanceStatusCompleted = "completed"
)

// HRM Statistics for dashboard
type HRMStats struct {
	TotalEmployees    int     `json:"total_employees"`
	ActiveEmployees   int     `json:"active_employees"`
	TotalDepartments  int     `json:"total_departments"`
	PendingLeaves     int     `json:"pending_leaves"`
	AvgPerformance    float64 `json:"avg_performance"`
	NewHiresThisMonth int     `json:"new_hires_this_month"`
	TurnoverRate      float64 `json:"turnover_rate"`
}

// Employee creation/update request
type EmployeeRequest struct {
	EmployeeCode   string     `json:"employee_code" validate:"required"`
	FirstName      string     `json:"first_name" validate:"required"`
	LastName       string     `json:"last_name" validate:"required"`
	Email          string     `json:"email" validate:"required,email"`
	Phone          string     `json:"phone"`
	DepartmentID   int        `json:"department_id" validate:"required"`
	Position       string     `json:"position" validate:"required"`
	HireDate       time.Time  `json:"hire_date" validate:"required"`
	Salary         float64    `json:"salary" validate:"required,min=0"`
	Status         string     `json:"status" validate:"required,oneof=active inactive terminated"`
	ManagerID      *int       `json:"manager_id"`
	Address        string     `json:"address"`
	DateOfBirth    *time.Time `json:"date_of_birth"`
	Gender         string     `json:"gender"`
	EmergencyName  string     `json:"emergency_name"`
	EmergencyPhone string     `json:"emergency_phone"`
}

// Leave request
type LeaveRequest struct {
	EmployeeID int       `json:"employee_id" validate:"required"`
	LeaveType  string    `json:"leave_type" validate:"required,oneof=annual sick maternity paternity personal emergency"`
	StartDate  time.Time `json:"start_date" validate:"required"`
	EndDate    time.Time `json:"end_date" validate:"required"`
	Reason     string    `json:"reason" validate:"required"`
}
