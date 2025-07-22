package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"zplus-saas/apps/backend/tenant-service/internal/models"
)

type TenantConfigService struct {
	db *sqlx.DB
}

func NewTenantConfigService(db *sqlx.DB) *TenantConfigService {
	return &TenantConfigService{db: db}
}

// GetTenantConfiguration gets tenant configuration
func (s *TenantConfigService) GetTenantConfiguration(tenantID uuid.UUID) (*models.TenantConfiguration, error) {
	var config models.TenantConfiguration
	err := s.db.Get(&config, "SELECT * FROM tenant_configurations WHERE tenant_id = $1", tenantID)
	if err != nil {
		// Create default configuration if not exists
		config = models.TenantConfiguration{
			ID:                 uuid.New(),
			TenantID:           tenantID,
			SSLEnabled:         false,
			BrandingConfig:     "{}",
			SecurityConfig:     "{}",
			NotificationConfig: "{}",
			IntegrationConfig:  "{}",
			FeatureFlags:       "{}",
			DataRetentionDays:  365,
			TwoFactorRequired:  false,
			PasswordPolicy:     `{"minLength": 8, "requireUppercase": true, "requireLowercase": true, "requireNumbers": true, "requireSpecialChars": false}`,
			SessionTimeoutMins: 480, // 8 hours
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}

		query := `
			INSERT INTO tenant_configurations (id, tenant_id, ssl_enabled, branding_config, security_config, 
				notification_config, integration_config, feature_flags, data_retention_days, 
				two_factor_required, password_policy, session_timeout_mins, created_at, updated_at)
			VALUES (:id, :tenant_id, :ssl_enabled, :branding_config, :security_config, 
				:notification_config, :integration_config, :feature_flags, :data_retention_days, 
				:two_factor_required, :password_policy, :session_timeout_mins, :created_at, :updated_at)`

		_, err = s.db.NamedExec(query, &config)
		if err != nil {
			return nil, fmt.Errorf("failed to create default configuration: %w", err)
		}
	}

	return &config, nil
}

// UpdateTenantConfiguration updates tenant configuration
func (s *TenantConfigService) UpdateTenantConfiguration(tenantID uuid.UUID, req *models.UpdateTenantConfigRequest) error {
	// Validate custom domain if provided
	if req.CustomDomain != nil && *req.CustomDomain != "" {
		err := s.validateCustomDomain(*req.CustomDomain)
		if err != nil {
			return fmt.Errorf("invalid custom domain: %w", err)
		}
	}

	// Validate JSON configurations
	if req.BrandingConfig != nil {
		err := s.validateJSONConfig(*req.BrandingConfig)
		if err != nil {
			return fmt.Errorf("invalid branding config: %w", err)
		}
	}

	if req.SecurityConfig != nil {
		err := s.validateJSONConfig(*req.SecurityConfig)
		if err != nil {
			return fmt.Errorf("invalid security config: %w", err)
		}
	}

	if req.NotificationConfig != nil {
		err := s.validateJSONConfig(*req.NotificationConfig)
		if err != nil {
			return fmt.Errorf("invalid notification config: %w", err)
		}
	}

	if req.IntegrationConfig != nil {
		err := s.validateJSONConfig(*req.IntegrationConfig)
		if err != nil {
			return fmt.Errorf("invalid integration config: %w", err)
		}
	}

	if req.FeatureFlags != nil {
		err := s.validateJSONConfig(*req.FeatureFlags)
		if err != nil {
			return fmt.Errorf("invalid feature flags: %w", err)
		}
	}

	if req.PasswordPolicy != nil {
		err := s.validatePasswordPolicy(*req.PasswordPolicy)
		if err != nil {
			return fmt.Errorf("invalid password policy: %w", err)
		}
	}

	// Build update query dynamically
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.CustomDomain != nil {
		setParts = append(setParts, fmt.Sprintf("custom_domain = $%d", argIndex))
		args = append(args, req.CustomDomain)
		argIndex++
	}

	if req.SSLEnabled != nil {
		setParts = append(setParts, fmt.Sprintf("ssl_enabled = $%d", argIndex))
		args = append(args, *req.SSLEnabled)
		argIndex++
	}

	if req.CustomCSS != nil {
		setParts = append(setParts, fmt.Sprintf("custom_css = $%d", argIndex))
		args = append(args, req.CustomCSS)
		argIndex++
	}

	if req.CustomJavaScript != nil {
		setParts = append(setParts, fmt.Sprintf("custom_javascript = $%d", argIndex))
		args = append(args, req.CustomJavaScript)
		argIndex++
	}

	if req.BrandingConfig != nil {
		setParts = append(setParts, fmt.Sprintf("branding_config = $%d", argIndex))
		args = append(args, *req.BrandingConfig)
		argIndex++
	}

	if req.SecurityConfig != nil {
		setParts = append(setParts, fmt.Sprintf("security_config = $%d", argIndex))
		args = append(args, *req.SecurityConfig)
		argIndex++
	}

	if req.NotificationConfig != nil {
		setParts = append(setParts, fmt.Sprintf("notification_config = $%d", argIndex))
		args = append(args, *req.NotificationConfig)
		argIndex++
	}

	if req.IntegrationConfig != nil {
		setParts = append(setParts, fmt.Sprintf("integration_config = $%d", argIndex))
		args = append(args, *req.IntegrationConfig)
		argIndex++
	}

	if req.FeatureFlags != nil {
		setParts = append(setParts, fmt.Sprintf("feature_flags = $%d", argIndex))
		args = append(args, *req.FeatureFlags)
		argIndex++
	}

	if req.DataRetentionDays != nil {
		setParts = append(setParts, fmt.Sprintf("data_retention_days = $%d", argIndex))
		args = append(args, *req.DataRetentionDays)
		argIndex++
	}

	if req.AllowedIPs != nil {
		setParts = append(setParts, fmt.Sprintf("allowed_ips = $%d", argIndex))
		args = append(args, req.AllowedIPs)
		argIndex++
	}

	if req.TwoFactorRequired != nil {
		setParts = append(setParts, fmt.Sprintf("two_factor_required = $%d", argIndex))
		args = append(args, *req.TwoFactorRequired)
		argIndex++
	}

	if req.PasswordPolicy != nil {
		setParts = append(setParts, fmt.Sprintf("password_policy = $%d", argIndex))
		args = append(args, *req.PasswordPolicy)
		argIndex++
	}

	if req.SessionTimeoutMins != nil {
		setParts = append(setParts, fmt.Sprintf("session_timeout_mins = $%d", argIndex))
		args = append(args, *req.SessionTimeoutMins)
		argIndex++
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Add updated_at
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add tenant_id for WHERE clause
	args = append(args, tenantID)

	query := fmt.Sprintf("UPDATE tenant_configurations SET %s WHERE tenant_id = $%d",
		strings.Join(setParts, ", "), argIndex)

	result, err := s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update configuration: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("configuration not found")
	}

	return nil
}

// SetupCustomDomain sets up custom domain for tenant
func (s *TenantConfigService) SetupCustomDomain(tenantID uuid.UUID, domain string, sslEnabled bool) error {
	// Validate domain
	err := s.validateCustomDomain(domain)
	if err != nil {
		return fmt.Errorf("invalid domain: %w", err)
	}

	// Check if domain is already used by another tenant
	var existingTenantID uuid.UUID
	err = s.db.Get(&existingTenantID,
		"SELECT tenant_id FROM tenant_configurations WHERE custom_domain = $1 AND tenant_id != $2",
		domain, tenantID)
	if err == nil {
		return fmt.Errorf("domain is already in use by another tenant")
	}

	// Update tenant configuration
	req := &models.UpdateTenantConfigRequest{
		CustomDomain: &domain,
		SSLEnabled:   &sslEnabled,
	}

	return s.UpdateTenantConfiguration(tenantID, req)
}

// RemoveCustomDomain removes custom domain from tenant
func (s *TenantConfigService) RemoveCustomDomain(tenantID uuid.UUID) error {
	_, err := s.db.Exec("UPDATE tenant_configurations SET custom_domain = NULL, ssl_enabled = false, updated_at = $1 WHERE tenant_id = $2",
		time.Now(), tenantID)
	if err != nil {
		return fmt.Errorf("failed to remove custom domain: %w", err)
	}
	return nil
}

// GetFeatureFlags gets feature flags for tenant
func (s *TenantConfigService) GetFeatureFlags(tenantID uuid.UUID) (map[string]interface{}, error) {
	config, err := s.GetTenantConfiguration(tenantID)
	if err != nil {
		return nil, err
	}

	var featureFlags map[string]interface{}
	err = json.Unmarshal([]byte(config.FeatureFlags), &featureFlags)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feature flags: %w", err)
	}

	return featureFlags, nil
}

// UpdateFeatureFlag updates a specific feature flag
func (s *TenantConfigService) UpdateFeatureFlag(tenantID uuid.UUID, flagName string, value interface{}) error {
	featureFlags, err := s.GetFeatureFlags(tenantID)
	if err != nil {
		return err
	}

	featureFlags[flagName] = value

	flagsJSON, err := json.Marshal(featureFlags)
	if err != nil {
		return fmt.Errorf("failed to marshal feature flags: %w", err)
	}

	flagsString := string(flagsJSON)
	req := &models.UpdateTenantConfigRequest{
		FeatureFlags: &flagsString,
	}

	return s.UpdateTenantConfiguration(tenantID, req)
}

// validateCustomDomain validates custom domain format
func (s *TenantConfigService) validateCustomDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain cannot be empty")
	}

	// Basic domain validation
	if !strings.Contains(domain, ".") {
		return fmt.Errorf("invalid domain format")
	}

	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return fmt.Errorf("domain cannot start or end with a dot")
	}

	if len(domain) > 253 {
		return fmt.Errorf("domain too long")
	}

	return nil
}

// validateJSONConfig validates JSON configuration
func (s *TenantConfigService) validateJSONConfig(config string) error {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(config), &result)
	if err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}
	return nil
}

// validatePasswordPolicy validates password policy configuration
func (s *TenantConfigService) validatePasswordPolicy(policy string) error {
	var policyMap map[string]interface{}
	err := json.Unmarshal([]byte(policy), &policyMap)
	if err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	// Validate required fields
	if minLength, exists := policyMap["minLength"]; exists {
		if length, ok := minLength.(float64); ok && length < 6 {
			return fmt.Errorf("minimum password length must be at least 6")
		}
	}

	return nil
}
