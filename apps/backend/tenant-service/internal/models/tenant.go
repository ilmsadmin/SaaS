package models

import (
	"time"

	"github.com/google/uuid"
)

// Tenant represents a tenant in the system
type Tenant struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Subdomain string    `json:"subdomain" db:"subdomain"`
	Domain    *string   `json:"domain" db:"domain"`
	Logo      *string   `json:"logo" db:"logo"`
	Status    string    `json:"status" db:"status"` // active, suspended, trial
	Settings  string    `json:"settings" db:"settings"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Plan represents a subscription plan
type Plan struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Description  *string   `json:"description" db:"description"`
	Price        float64   `json:"price" db:"price"`
	Currency     string    `json:"currency" db:"currency"`
	BillingCycle string    `json:"billing_cycle" db:"billing_cycle"` // monthly, yearly
	MaxUsers     *int      `json:"max_users" db:"max_users"`
	MaxStorage   *int64    `json:"max_storage" db:"max_storage"` // in bytes
	Features     string    `json:"features" db:"features"`       // JSON array
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Subscription represents a tenant's subscription to a plan
type Subscription struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	TenantID           uuid.UUID  `json:"tenant_id" db:"tenant_id"`
	PlanID             uuid.UUID  `json:"plan_id" db:"plan_id"`
	Status             string     `json:"status" db:"status"` // active, cancelled, expired
	TrialEndAt         *time.Time `json:"trial_end_at" db:"trial_end_at"`
	CurrentPeriodStart time.Time  `json:"current_period_start" db:"current_period_start"`
	CurrentPeriodEnd   time.Time  `json:"current_period_end" db:"current_period_end"`
	CancelAtPeriodEnd  bool       `json:"cancel_at_period_end" db:"cancel_at_period_end"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`

	// Relationships
	Tenant *Tenant `json:"tenant,omitempty"`
	Plan   *Plan   `json:"plan,omitempty"`
}

// TenantModule represents modules enabled for a tenant
type TenantModule struct {
	ID          uuid.UUID `json:"id" db:"id"`
	TenantID    uuid.UUID `json:"tenant_id" db:"tenant_id"`
	ModuleID    uuid.UUID `json:"module_id" db:"module_id"`
	IsEnabled   bool      `json:"is_enabled" db:"is_enabled"`
	Config      string    `json:"config" db:"config"` // JSON config
	InstalledAt time.Time `json:"installed_at" db:"installed_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Module represents available modules in the system
type Module struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	DisplayName  string    `json:"display_name" db:"display_name"`
	Description  *string   `json:"description" db:"description"`
	Version      string    `json:"version" db:"version"`
	Category     string    `json:"category" db:"category"`
	Icon         *string   `json:"icon" db:"icon"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	Dependencies string    `json:"dependencies" db:"dependencies"` // JSON array
	Permissions  string    `json:"permissions" db:"permissions"`   // JSON array
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// TenantConfiguration represents advanced tenant configuration
type TenantConfiguration struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	TenantID           uuid.UUID `json:"tenant_id" db:"tenant_id"`
	CustomDomain       *string   `json:"custom_domain" db:"custom_domain"`
	SSLEnabled         bool      `json:"ssl_enabled" db:"ssl_enabled"`
	CustomCSS          *string   `json:"custom_css" db:"custom_css"`
	CustomJavaScript   *string   `json:"custom_javascript" db:"custom_javascript"`
	BrandingConfig     string    `json:"branding_config" db:"branding_config"`         // JSON
	SecurityConfig     string    `json:"security_config" db:"security_config"`         // JSON
	NotificationConfig string    `json:"notification_config" db:"notification_config"` // JSON
	IntegrationConfig  string    `json:"integration_config" db:"integration_config"`   // JSON
	FeatureFlags       string    `json:"feature_flags" db:"feature_flags"`             // JSON
	DataRetentionDays  int       `json:"data_retention_days" db:"data_retention_days"`
	AllowedIPs         *string   `json:"allowed_ips" db:"allowed_ips"` // JSON array
	TwoFactorRequired  bool      `json:"two_factor_required" db:"two_factor_required"`
	PasswordPolicy     string    `json:"password_policy" db:"password_policy"` // JSON
	SessionTimeoutMins int       `json:"session_timeout_mins" db:"session_timeout_mins"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// ModuleDependency represents module dependencies
type ModuleDependency struct {
	ID            uuid.UUID `json:"id" db:"id"`
	ModuleID      uuid.UUID `json:"module_id" db:"module_id"`
	DependsOnID   uuid.UUID `json:"depends_on_id" db:"depends_on_id"`
	MinVersion    string    `json:"min_version" db:"min_version"`
	IsRequired    bool      `json:"is_required" db:"is_required"`
	ConflictsWith bool      `json:"conflicts_with" db:"conflicts_with"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// ModuleMarketplace represents modules available in marketplace
type ModuleMarketplace struct {
	ID             uuid.UUID `json:"id" db:"id"`
	ModuleID       uuid.UUID `json:"module_id" db:"module_id"`
	PublisherID    uuid.UUID `json:"publisher_id" db:"publisher_id"`
	Price          float64   `json:"price" db:"price"`
	Currency       string    `json:"currency" db:"currency"`
	IsFree         bool      `json:"is_free" db:"is_free"`
	Rating         float32   `json:"rating" db:"rating"`
	TotalDownloads int       `json:"total_downloads" db:"total_downloads"`
	Screenshots    string    `json:"screenshots" db:"screenshots"` // JSON array
	Documentation  *string   `json:"documentation" db:"documentation"`
	SupportEmail   *string   `json:"support_email" db:"support_email"`
	Homepage       *string   `json:"homepage" db:"homepage"`
	Repository     *string   `json:"repository" db:"repository"`
	License        string    `json:"license" db:"license"`
	Tags           string    `json:"tags" db:"tags"` // JSON array
	IsVerified     bool      `json:"is_verified" db:"is_verified"`
	IsFeatured     bool      `json:"is_featured" db:"is_featured"`
	PublishedAt    time.Time `json:"published_at" db:"published_at"`
	LastUpdatedAt  time.Time `json:"last_updated_at" db:"last_updated_at"`
}

// ModuleInstallation represents module installation for tenants
type ModuleInstallation struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	TenantID      uuid.UUID  `json:"tenant_id" db:"tenant_id"`
	ModuleID      uuid.UUID  `json:"module_id" db:"module_id"`
	Version       string     `json:"version" db:"version"`
	Status        string     `json:"status" db:"status"`             // installing, installed, failed, updating, uninstalling
	Config        string     `json:"config" db:"config"`             // JSON
	InstallData   string     `json:"install_data" db:"install_data"` // JSON
	ErrorMessage  *string    `json:"error_message" db:"error_message"`
	InstalledBy   uuid.UUID  `json:"installed_by" db:"installed_by"`
	InstalledAt   time.Time  `json:"installed_at" db:"installed_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	UninstalledAt *time.Time `json:"uninstalled_at" db:"uninstalled_at"`
}

// ModulePermission represents permissions for modules
type ModulePermission struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ModuleID    uuid.UUID `json:"module_id" db:"module_id"`
	Permission  string    `json:"permission" db:"permission"`
	Description string    `json:"description" db:"description"`
	Category    string    `json:"category" db:"category"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Request/Response DTOs
type CreateTenantRequest struct {
	Name      string  `json:"name" validate:"required,min=2,max=255"`
	Subdomain string  `json:"subdomain" validate:"required,min=3,max=63,alphanum"`
	Domain    *string `json:"domain,omitempty"`
	PlanID    *string `json:"plan_id,omitempty"`
}

type UpdateTenantRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	Domain   *string `json:"domain,omitempty"`
	Logo     *string `json:"logo,omitempty"`
	Settings *string `json:"settings,omitempty"`
}

type CreatePlanRequest struct {
	Name         string  `json:"name" validate:"required,min=2,max=100"`
	Description  *string `json:"description,omitempty"`
	Price        float64 `json:"price" validate:"required,min=0"`
	Currency     string  `json:"currency" validate:"required,len=3"`
	BillingCycle string  `json:"billing_cycle" validate:"required,oneof=monthly yearly"`
	MaxUsers     *int    `json:"max_users,omitempty" validate:"omitempty,min=1"`
	MaxStorage   *int64  `json:"max_storage,omitempty" validate:"omitempty,min=1"`
	Features     *string `json:"features,omitempty"`
}

type UpdatePlanRequest struct {
	Name         *string  `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description  *string  `json:"description,omitempty"`
	Price        *float64 `json:"price,omitempty" validate:"omitempty,min=0"`
	Currency     *string  `json:"currency,omitempty" validate:"omitempty,len=3"`
	BillingCycle *string  `json:"billing_cycle,omitempty" validate:"omitempty,oneof=monthly yearly"`
	MaxUsers     *int     `json:"max_users,omitempty" validate:"omitempty,min=1"`
	MaxStorage   *int64   `json:"max_storage,omitempty" validate:"omitempty,min=1"`
	Features     *string  `json:"features,omitempty"`
	IsActive     *bool    `json:"is_active,omitempty"`
}

type CreateSubscriptionRequest struct {
	PlanID             string     `json:"plan_id" validate:"required,uuid"`
	TrialEndAt         *time.Time `json:"trial_end_at,omitempty"`
	CurrentPeriodStart time.Time  `json:"current_period_start" validate:"required"`
	CurrentPeriodEnd   time.Time  `json:"current_period_end" validate:"required"`
}

type UpdateSubscriptionRequest struct {
	PlanID             *string    `json:"plan_id,omitempty" validate:"omitempty,uuid"`
	Status             *string    `json:"status,omitempty" validate:"omitempty,oneof=active cancelled expired"`
	TrialEndAt         *time.Time `json:"trial_end_at,omitempty"`
	CurrentPeriodStart *time.Time `json:"current_period_start,omitempty"`
	CurrentPeriodEnd   *time.Time `json:"current_period_end,omitempty"`
	CancelAtPeriodEnd  *bool      `json:"cancel_at_period_end,omitempty"`
}

// Request/Response DTOs for Module System
type InstallModuleRequest struct {
	ModuleID string `json:"module_id" validate:"required,uuid"`
	Version  string `json:"version" validate:"required"`
	Config   string `json:"config,omitempty"`
}

type UpdateModuleConfigRequest struct {
	Config string `json:"config" validate:"required"`
}

type UpdateTenantConfigRequest struct {
	CustomDomain       *string `json:"custom_domain,omitempty"`
	SSLEnabled         *bool   `json:"ssl_enabled,omitempty"`
	CustomCSS          *string `json:"custom_css,omitempty"`
	CustomJavaScript   *string `json:"custom_javascript,omitempty"`
	BrandingConfig     *string `json:"branding_config,omitempty"`
	SecurityConfig     *string `json:"security_config,omitempty"`
	NotificationConfig *string `json:"notification_config,omitempty"`
	IntegrationConfig  *string `json:"integration_config,omitempty"`
	FeatureFlags       *string `json:"feature_flags,omitempty"`
	DataRetentionDays  *int    `json:"data_retention_days,omitempty"`
	AllowedIPs         *string `json:"allowed_ips,omitempty"`
	TwoFactorRequired  *bool   `json:"two_factor_required,omitempty"`
	PasswordPolicy     *string `json:"password_policy,omitempty"`
	SessionTimeoutMins *int    `json:"session_timeout_mins,omitempty"`
}

// Constants
const (
	TenantStatusActive    = "active"
	TenantStatusSuspended = "suspended"
	TenantStatusTrial     = "trial"

	SubscriptionStatusActive    = "active"
	SubscriptionStatusCancelled = "cancelled"
	SubscriptionStatusExpired   = "expired"

	BillingCycleMonthly = "monthly"
	BillingCycleYearly  = "yearly"

	ModuleStatusInstalling   = "installing"
	ModuleStatusInstalled    = "installed"
	ModuleStatusFailed       = "failed"
	ModuleStatusUpdating     = "updating"
	ModuleStatusUninstalling = "uninstalling"

	ModuleCategoryCRM        = "crm"
	ModuleCategoryHRM        = "hrm"
	ModuleCategoryPOS        = "pos"
	ModuleCategoryLMS        = "lms"
	ModuleCategoryCheckin    = "checkin"
	ModuleCategoryPayment    = "payment"
	ModuleCategoryAccounting = "accounting"
	ModuleCategoryEcommerce  = "ecommerce"
)
