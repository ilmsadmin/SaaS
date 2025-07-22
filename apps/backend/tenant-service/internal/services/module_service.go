package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"zplus-saas/apps/backend/tenant-service/internal/models"
)

type ModuleService struct {
	db *sqlx.DB
}

func NewModuleService(db *sqlx.DB) *ModuleService {
	return &ModuleService{db: db}
}

// GetAvailableModules gets all available modules
func (s *ModuleService) GetAvailableModules() ([]models.Module, error) {
	var modules []models.Module
	err := s.db.Select(&modules, "SELECT * FROM modules WHERE is_active = true ORDER BY category, name")
	if err != nil {
		return nil, fmt.Errorf("failed to get modules: %w", err)
	}
	return modules, nil
}

// GetTenantModules gets modules for a specific tenant
func (s *ModuleService) GetTenantModules(tenantID uuid.UUID) ([]models.TenantModule, error) {
	var tenantModules []models.TenantModule
	query := `
		SELECT tm.*, m.name as module_name, m.display_name, m.description, m.version, m.category, m.icon
		FROM tenant_modules tm
		JOIN modules m ON tm.module_id = m.id
		WHERE tm.tenant_id = $1
		ORDER BY m.category, m.name`

	err := s.db.Select(&tenantModules, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant modules: %w", err)
	}
	return tenantModules, nil
}

// InstallModule installs a module for a tenant
func (s *ModuleService) InstallModule(tenantID, moduleID, userID uuid.UUID, version, config string) (*models.ModuleInstallation, error) {
	// Check if module exists and is active
	var module models.Module
	err := s.db.Get(&module, "SELECT * FROM modules WHERE id = $1 AND is_active = true", moduleID)
	if err != nil {
		return nil, fmt.Errorf("module not found or not active")
	}

	// Check if module is already installed
	var existingInstallation models.ModuleInstallation
	err = s.db.Get(&existingInstallation,
		"SELECT id FROM module_installations WHERE tenant_id = $1 AND module_id = $2 AND uninstalled_at IS NULL",
		tenantID, moduleID)
	if err == nil {
		return nil, fmt.Errorf("module is already installed")
	}

	// Check dependencies
	err = s.checkModuleDependencies(tenantID, moduleID)
	if err != nil {
		return nil, fmt.Errorf("dependency check failed: %w", err)
	}

	// Start transaction
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Create module installation record
	installation := &models.ModuleInstallation{
		ID:          uuid.New(),
		TenantID:    tenantID,
		ModuleID:    moduleID,
		Version:     version,
		Status:      models.ModuleStatusInstalling,
		Config:      config,
		InstallData: "{}",
		InstalledBy: userID,
		InstalledAt: time.Now(),
		UpdatedAt:   time.Now(),
	}

	installQuery := `
		INSERT INTO module_installations (id, tenant_id, module_id, version, status, config, install_data, installed_by, installed_at, updated_at)
		VALUES (:id, :tenant_id, :module_id, :version, :status, :config, :install_data, :installed_by, :installed_at, :updated_at)`

	_, err = tx.NamedExec(installQuery, installation)
	if err != nil {
		return nil, fmt.Errorf("failed to create installation record: %w", err)
	}

	// Create or update tenant module record
	tenantModule := &models.TenantModule{
		ID:          uuid.New(),
		TenantID:    tenantID,
		ModuleID:    moduleID,
		IsEnabled:   true,
		Config:      config,
		InstalledAt: time.Now(),
		UpdatedAt:   time.Now(),
	}

	tenantModuleQuery := `
		INSERT INTO tenant_modules (id, tenant_id, module_id, is_enabled, config, installed_at, updated_at)
		VALUES (:id, :tenant_id, :module_id, :is_enabled, :config, :installed_at, :updated_at)
		ON CONFLICT (tenant_id, module_id) 
		DO UPDATE SET is_enabled = :is_enabled, config = :config, updated_at = :updated_at`

	_, err = tx.NamedExec(tenantModuleQuery, tenantModule)
	if err != nil {
		return nil, fmt.Errorf("failed to create tenant module record: %w", err)
	}

	// Update installation status to installed
	_, err = tx.Exec("UPDATE module_installations SET status = $1, updated_at = $2 WHERE id = $3",
		models.ModuleStatusInstalled, time.Now(), installation.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update installation status: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	installation.Status = models.ModuleStatusInstalled
	return installation, nil
}

// UninstallModule uninstalls a module for a tenant
func (s *ModuleService) UninstallModule(tenantID, moduleID uuid.UUID) error {
	// Check if module is installed
	var installation models.ModuleInstallation
	err := s.db.Get(&installation,
		"SELECT * FROM module_installations WHERE tenant_id = $1 AND module_id = $2 AND uninstalled_at IS NULL",
		tenantID, moduleID)
	if err != nil {
		return fmt.Errorf("module is not installed")
	}

	// Check if other modules depend on this one
	dependentModules, err := s.getDependentModules(tenantID, moduleID)
	if err != nil {
		return fmt.Errorf("failed to check dependencies: %w", err)
	}

	if len(dependentModules) > 0 {
		return fmt.Errorf("cannot uninstall module: other modules depend on it")
	}

	// Start transaction
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Update installation status
	_, err = tx.Exec("UPDATE module_installations SET status = $1, uninstalled_at = $2, updated_at = $3 WHERE id = $4",
		models.ModuleStatusUninstalling, time.Now(), time.Now(), installation.ID)
	if err != nil {
		return fmt.Errorf("failed to update installation status: %w", err)
	}

	// Disable tenant module
	_, err = tx.Exec("UPDATE tenant_modules SET is_enabled = false, updated_at = $1 WHERE tenant_id = $2 AND module_id = $3",
		time.Now(), tenantID, moduleID)
	if err != nil {
		return fmt.Errorf("failed to disable tenant module: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// EnableModule enables a module for a tenant
func (s *ModuleService) EnableModule(tenantID, moduleID uuid.UUID) error {
	result, err := s.db.Exec("UPDATE tenant_modules SET is_enabled = true, updated_at = $1 WHERE tenant_id = $2 AND module_id = $3",
		time.Now(), tenantID, moduleID)
	if err != nil {
		return fmt.Errorf("failed to enable module: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("module not found for tenant")
	}

	return nil
}

// DisableModule disables a module for a tenant
func (s *ModuleService) DisableModule(tenantID, moduleID uuid.UUID) error {
	// Check if other modules depend on this one
	dependentModules, err := s.getDependentModules(tenantID, moduleID)
	if err != nil {
		return fmt.Errorf("failed to check dependencies: %w", err)
	}

	if len(dependentModules) > 0 {
		return fmt.Errorf("cannot disable module: other modules depend on it")
	}

	result, err := s.db.Exec("UPDATE tenant_modules SET is_enabled = false, updated_at = $1 WHERE tenant_id = $2 AND module_id = $3",
		time.Now(), tenantID, moduleID)
	if err != nil {
		return fmt.Errorf("failed to disable module: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("module not found for tenant")
	}

	return nil
}

// UpdateModuleConfig updates module configuration
func (s *ModuleService) UpdateModuleConfig(tenantID, moduleID uuid.UUID, config string) error {
	// Validate JSON config
	var configMap map[string]interface{}
	err := json.Unmarshal([]byte(config), &configMap)
	if err != nil {
		return fmt.Errorf("invalid JSON config: %w", err)
	}

	result, err := s.db.Exec("UPDATE tenant_modules SET config = $1, updated_at = $2 WHERE tenant_id = $3 AND module_id = $4",
		config, time.Now(), tenantID, moduleID)
	if err != nil {
		return fmt.Errorf("failed to update module config: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("module not found for tenant")
	}

	return nil
}

// checkModuleDependencies checks if all required dependencies are installed
func (s *ModuleService) checkModuleDependencies(tenantID, moduleID uuid.UUID) error {
	var dependencies []models.ModuleDependency
	err := s.db.Select(&dependencies,
		"SELECT * FROM module_dependencies WHERE module_id = $1 AND is_required = true",
		moduleID)
	if err != nil {
		return fmt.Errorf("failed to get dependencies: %w", err)
	}

	for _, dep := range dependencies {
		var tenantModule models.TenantModule
		err = s.db.Get(&tenantModule,
			"SELECT * FROM tenant_modules WHERE tenant_id = $1 AND module_id = $2 AND is_enabled = true",
			tenantID, dep.DependsOnID)
		if err != nil {
			var depModule models.Module
			s.db.Get(&depModule, "SELECT name FROM modules WHERE id = $1", dep.DependsOnID)
			return fmt.Errorf("required dependency '%s' is not installed or enabled", depModule.Name)
		}
	}

	return nil
}

// getDependentModules gets modules that depend on the given module
func (s *ModuleService) getDependentModules(tenantID, moduleID uuid.UUID) ([]models.TenantModule, error) {
	var dependentModules []models.TenantModule
	query := `
		SELECT tm.* FROM tenant_modules tm
		JOIN module_dependencies md ON tm.module_id = md.module_id
		WHERE tm.tenant_id = $1 AND md.depends_on_id = $2 AND tm.is_enabled = true AND md.is_required = true`

	err := s.db.Select(&dependentModules, query, tenantID, moduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependent modules: %w", err)
	}

	return dependentModules, nil
}
