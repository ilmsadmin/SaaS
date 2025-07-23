package models

import (
	"time"
)

// Customer represents a customer in the CRM system
type Customer struct {
	ID        int       `json:"id" db:"id"`
	TenantID  string    `json:"tenant_id" db:"tenant_id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone" db:"phone"`
	Company   string    `json:"company" db:"company"`
	Address   string    `json:"address" db:"address"`
	City      string    `json:"city" db:"city"`
	State     string    `json:"state" db:"state"`
	Country   string    `json:"country" db:"country"`
	ZipCode   string    `json:"zip_code" db:"zip_code"`
	Status    string    `json:"status" db:"status"` // active, inactive, prospect
	Source    string    `json:"source" db:"source"` // website, referral, social, etc.
	Tags      []string  `json:"tags" db:"tags"`
	Notes     string    `json:"notes" db:"notes"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Lead represents a sales lead
type Lead struct {
	ID          int        `json:"id" db:"id"`
	TenantID    string     `json:"tenant_id" db:"tenant_id"`
	Name        string     `json:"name" db:"name"`
	Email       string     `json:"email" db:"email"`
	Phone       string     `json:"phone" db:"phone"`
	Company     string     `json:"company" db:"company"`
	Title       string     `json:"title" db:"title"`
	Source      string     `json:"source" db:"source"`
	Status      string     `json:"status" db:"status"`           // new, qualified, contacted, converted, lost
	Score       int        `json:"score" db:"score"`             // Lead scoring (0-100)
	AssignedTo  int        `json:"assigned_to" db:"assigned_to"` // User ID
	Value       float64    `json:"value" db:"value"`             // Potential deal value
	Notes       string     `json:"notes" db:"notes"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	ConvertedAt *time.Time `json:"converted_at" db:"converted_at"`
}

// Opportunity represents a sales opportunity
type Opportunity struct {
	ID           int        `json:"id" db:"id"`
	TenantID     string     `json:"tenant_id" db:"tenant_id"`
	CustomerID   int        `json:"customer_id" db:"customer_id"`
	Name         string     `json:"name" db:"name"`
	Description  string     `json:"description" db:"description"`
	Value        float64    `json:"value" db:"value"`
	Currency     string     `json:"currency" db:"currency"`
	Stage        string     `json:"stage" db:"stage"`             // prospecting, qualification, proposal, negotiation, closed-won, closed-lost
	Probability  int        `json:"probability" db:"probability"` // 0-100%
	Source       string     `json:"source" db:"source"`
	AssignedTo   int        `json:"assigned_to" db:"assigned_to"`
	ExpectedDate time.Time  `json:"expected_date" db:"expected_date"`
	ClosedDate   *time.Time `json:"closed_date" db:"closed_date"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// ContactActivity represents interactions with customers/leads
type ContactActivity struct {
	ID          int        `json:"id" db:"id"`
	TenantID    string     `json:"tenant_id" db:"tenant_id"`
	CustomerID  *int       `json:"customer_id" db:"customer_id"`
	LeadID      *int       `json:"lead_id" db:"lead_id"`
	Type        string     `json:"type" db:"type"` // call, email, meeting, note
	Subject     string     `json:"subject" db:"subject"`
	Description string     `json:"description" db:"description"`
	Duration    int        `json:"duration" db:"duration"` // in minutes
	UserID      int        `json:"user_id" db:"user_id"`
	ScheduledAt *time.Time `json:"scheduled_at" db:"scheduled_at"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// CustomerRequest represents common request structures
type CreateCustomerRequest struct {
	Name    string   `json:"name" validate:"required"`
	Email   string   `json:"email" validate:"required,email"`
	Phone   string   `json:"phone"`
	Company string   `json:"company"`
	Address string   `json:"address"`
	City    string   `json:"city"`
	State   string   `json:"state"`
	Country string   `json:"country"`
	ZipCode string   `json:"zip_code"`
	Source  string   `json:"source"`
	Tags    []string `json:"tags"`
	Notes   string   `json:"notes"`
}

type UpdateCustomerRequest struct {
	Name    *string  `json:"name"`
	Email   *string  `json:"email" validate:"omitempty,email"`
	Phone   *string  `json:"phone"`
	Company *string  `json:"company"`
	Address *string  `json:"address"`
	City    *string  `json:"city"`
	State   *string  `json:"state"`
	Country *string  `json:"country"`
	ZipCode *string  `json:"zip_code"`
	Status  *string  `json:"status"`
	Source  *string  `json:"source"`
	Tags    []string `json:"tags"`
	Notes   *string  `json:"notes"`
}

// Lead Request structures
type CreateLeadRequest struct {
	Name       string  `json:"name" validate:"required"`
	Email      string  `json:"email" validate:"required,email"`
	Phone      string  `json:"phone"`
	Company    string  `json:"company"`
	Title      string  `json:"title"`
	Source     string  `json:"source"`
	AssignedTo int     `json:"assigned_to"`
	Value      float64 `json:"value"`
	Notes      string  `json:"notes"`
}

type UpdateLeadRequest struct {
	Name       *string  `json:"name"`
	Email      *string  `json:"email" validate:"omitempty,email"`
	Phone      *string  `json:"phone"`
	Company    *string  `json:"company"`
	Title      *string  `json:"title"`
	Source     *string  `json:"source"`
	Status     *string  `json:"status"`
	Score      *int     `json:"score"`
	AssignedTo *int     `json:"assigned_to"`
	Value      *float64 `json:"value"`
	Notes      *string  `json:"notes"`
}

// Opportunity Request structures
type CreateOpportunityRequest struct {
	CustomerID   int     `json:"customer_id" validate:"required"`
	Name         string  `json:"name" validate:"required"`
	Description  string  `json:"description"`
	Value        float64 `json:"value" validate:"required"`
	Currency     string  `json:"currency"`
	Source       string  `json:"source"`
	AssignedTo   int     `json:"assigned_to"`
	ExpectedDate string  `json:"expected_date" validate:"required"` // ISO date format
}

type UpdateOpportunityRequest struct {
	Name         *string  `json:"name"`
	Description  *string  `json:"description"`
	Value        *float64 `json:"value"`
	Currency     *string  `json:"currency"`
	Stage        *string  `json:"stage"`
	Probability  *int     `json:"probability"`
	Source       *string  `json:"source"`
	AssignedTo   *int     `json:"assigned_to"`
	ExpectedDate *string  `json:"expected_date"`
}
