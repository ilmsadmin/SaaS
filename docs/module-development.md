# Module Development Guide - Zplus SaaS

## 1. Tổng quan Module System

Zplus SaaS được thiết kế với kiến trúc module linh hoạt, cho phép bật/tắt các tính năng theo nhu cầu từng tenant. Mỗi module là một microservice độc lập với database schema riêng.

### Module Architecture

```
Module Structure:
├── Core Interface       ← BaseModule interface
├── Database Schema      ← Tables specific to module
├── API Endpoints        ← REST/GraphQL endpoints
├── Business Logic       ← Services and handlers
├── Permissions         ← RBAC permissions
└── Frontend Components  ← UI components (optional)
```

### Available Modules

| Module | Description | Status |
|--------|-------------|--------|
| **CRM** | Customer Relationship Management | ✅ Active |
| **LMS** | Learning Management System | ✅ Active |
| **POS** | Point of Sale | ✅ Active |
| **HRM** | Human Resource Management | ✅ Active |
| **Checkin** | Attendance & Location Tracking | ✅ Active |
| **Analytics** | Business Intelligence | 🚧 Planned |
| **Chat** | Real-time Messaging | 🚧 Planned |

## 2. Module Interface & Lifecycle

### 2.1 BaseModule Interface

**Go Module Interface:**
```go
package modules

import (
    "context"
    "database/sql/driver"
)

type BaseModule interface {
    // Module identification
    GetName() string
    GetVersion() string
    GetDisplayName() string
    GetDescription() string
    GetCategory() ModuleCategory
    GetIcon() string
    
    // Dependencies
    GetDependencies() []string
    GetConflicts() []string
    
    // Lifecycle management
    Install(ctx context.Context, tenantID string) error
    Uninstall(ctx context.Context, tenantID string) error
    Enable(ctx context.Context, tenantID string) error
    Disable(ctx context.Context, tenantID string) error
    
    // Configuration
    GetDefaultConfig() map[string]interface{}
    ValidateConfig(config map[string]interface{}) error
    
    // API registration
    GetRoutes() []Route
    GetGraphQLSchema() string
    GetPermissions() []Permission
    
    // Database operations
    CreateSchema(ctx context.Context, tenantID string) error
    DropSchema(ctx context.Context, tenantID string) error
    MigrateSchema(ctx context.Context, tenantID string) error
    
    // Health check
    HealthCheck(ctx context.Context, tenantID string) error
}

type ModuleCategory string

const (
    CategoryCRM       ModuleCategory = "crm"
    CategoryLMS       ModuleCategory = "lms"
    CategoryPOS       ModuleCategory = "pos"
    CategoryHRM       ModuleCategory = "hrm"
    CategoryAnalytics ModuleCategory = "analytics"
    CategoryOther     ModuleCategory = "other"
)

type Route struct {
    Method      string
    Path        string
    Handler     interface{}
    Middlewares []interface{}
    Permission  string
}

type Permission struct {
    Name        string
    Resource    string
    Action      string
    Description string
}
```

### 2.2 Module Registration

**Module Registry:**
```go
package registry

import (
    "fmt"
    "sync"
)

type ModuleRegistry struct {
    modules map[string]BaseModule
    mu      sync.RWMutex
}

var GlobalRegistry = &ModuleRegistry{
    modules: make(map[string]BaseModule),
}

func (r *ModuleRegistry) Register(module BaseModule) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    name := module.GetName()
    if _, exists := r.modules[name]; exists {
        return fmt.Errorf("module %s already registered", name)
    }
    
    // Validate module
    if err := r.validateModule(module); err != nil {
        return fmt.Errorf("module validation failed: %v", err)
    }
    
    r.modules[name] = module
    return nil
}

func (r *ModuleRegistry) GetModule(name string) (BaseModule, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    module, exists := r.modules[name]
    if !exists {
        return nil, fmt.Errorf("module %s not found", name)
    }
    
    return module, nil
}

func (r *ModuleRegistry) GetAllModules() map[string]BaseModule {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    result := make(map[string]BaseModule)
    for name, module := range r.modules {
        result[name] = module
    }
    
    return result
}

func (r *ModuleRegistry) validateModule(module BaseModule) error {
    // Check required fields
    if module.GetName() == "" {
        return fmt.Errorf("module name is required")
    }
    
    if module.GetVersion() == "" {
        return fmt.Errorf("module version is required")
    }
    
    // Validate dependencies
    for _, dep := range module.GetDependencies() {
        if _, exists := r.modules[dep]; !exists {
            return fmt.Errorf("dependency %s not found", dep)
        }
    }
    
    return nil
}
```

## 3. Developing a New Module

### 3.1 Module Structure

**Project Structure for a new module:**
```
apps/backend/mymodule/
├── main.go                 # Entry point & module registration
├── module.go              # Module implementation
├── models/                # Database models
│   ├── model1.go
│   └── model2.go
├── services/              # Business logic
│   ├── service1.go
│   └── service2.go
├── handlers/              # HTTP handlers
│   ├── rest.go
│   └── graphql.go
├── migrations/            # Database migrations
│   ├── 001_initial.up.sql
│   └── 001_initial.down.sql
├── config/               # Module configuration
│   └── config.go
├── schemas/              # GraphQL schemas
│   └── schema.graphql
└── tests/                # Unit tests
    ├── module_test.go
    └── integration_test.go
```

### 3.2 Implementing a Sample Module

**Example: Simple Task Management Module**

**main.go:**
```go
package main

import (
    "log"
    
    "github.com/zplus/saas/apps/backend/shared/registry"
)

func main() {
    // Create and register module
    taskModule := NewTaskModule()
    
    if err := registry.GlobalRegistry.Register(taskModule); err != nil {
        log.Fatal("Failed to register task module:", err)
    }
    
    log.Println("Task module registered successfully")
}
```

**module.go:**
```go
package main

import (
    "context"
    "fmt"
    
    "github.com/zplus/saas/apps/backend/shared/modules"
)

type TaskModule struct {
    name        string
    version     string
    displayName string
    description string
}

func NewTaskModule() *TaskModule {
    return &TaskModule{
        name:        "tasks",
        version:     "1.0.0",
        displayName: "Task Management",
        description: "Simple task management with assignments and deadlines",
    }
}

func (m *TaskModule) GetName() string        { return m.name }
func (m *TaskModule) GetVersion() string     { return m.version }
func (m *TaskModule) GetDisplayName() string { return m.displayName }
func (m *TaskModule) GetDescription() string { return m.description }
func (m *TaskModule) GetCategory() modules.ModuleCategory { return modules.CategoryOther }
func (m *TaskModule) GetIcon() string        { return "task-icon.svg" }

func (m *TaskModule) GetDependencies() []string {
    return []string{} // No dependencies
}

func (m *TaskModule) GetConflicts() []string {
    return []string{} // No conflicts
}

func (m *TaskModule) GetDefaultConfig() map[string]interface{} {
    return map[string]interface{}{
        "max_tasks_per_user":     100,
        "allow_file_attachments": true,
        "notification_enabled":   true,
    }
}

func (m *TaskModule) ValidateConfig(config map[string]interface{}) error {
    if maxTasks, ok := config["max_tasks_per_user"].(float64); ok {
        if maxTasks < 1 || maxTasks > 1000 {
            return fmt.Errorf("max_tasks_per_user must be between 1 and 1000")
        }
    }
    return nil
}

func (m *TaskModule) GetRoutes() []modules.Route {
    return []modules.Route{
        {
            Method:     "GET",
            Path:       "/api/v1/tasks",
            Handler:    "GetTasks",
            Permission: "tasks:read",
        },
        {
            Method:     "POST",
            Path:       "/api/v1/tasks",
            Handler:    "CreateTask",
            Permission: "tasks:write",
        },
        {
            Method:     "PUT",
            Path:       "/api/v1/tasks/:id",
            Handler:    "UpdateTask",
            Permission: "tasks:write",
        },
        {
            Method:     "DELETE",
            Path:       "/api/v1/tasks/:id",
            Handler:    "DeleteTask",
            Permission: "tasks:delete",
        },
    }
}

func (m *TaskModule) GetGraphQLSchema() string {
    return `
        type Task {
            id: ID!
            title: String!
            description: String
            status: TaskStatus!
            priority: TaskPriority!
            assignedTo: User
            dueDate: DateTime
            createdAt: DateTime!
            updatedAt: DateTime!
        }
        
        enum TaskStatus {
            TODO
            IN_PROGRESS
            COMPLETED
            CANCELLED
        }
        
        enum TaskPriority {
            LOW
            MEDIUM
            HIGH
            URGENT
        }
        
        extend type Query {
            tasks(filter: TaskFilter): [Task!]!
            task(id: ID!): Task
        }
        
        extend type Mutation {
            createTask(input: CreateTaskInput!): Task!
            updateTask(id: ID!, input: UpdateTaskInput!): Task!
            deleteTask(id: ID!): Boolean!
        }
        
        input TaskFilter {
            status: TaskStatus
            priority: TaskPriority
            assignedTo: ID
        }
        
        input CreateTaskInput {
            title: String!
            description: String
            priority: TaskPriority!
            assignedTo: ID
            dueDate: DateTime
        }
        
        input UpdateTaskInput {
            title: String
            description: String
            status: TaskStatus
            priority: TaskPriority
            assignedTo: ID
            dueDate: DateTime
        }
    `
}

func (m *TaskModule) GetPermissions() []modules.Permission {
    return []modules.Permission{
        {
            Name:        "tasks:read",
            Resource:    "tasks",
            Action:      "read",
            Description: "Read tasks",
        },
        {
            Name:        "tasks:write",
            Resource:    "tasks",
            Action:      "write",
            Description: "Create and update tasks",
        },
        {
            Name:        "tasks:delete",
            Resource:    "tasks",
            Action:      "delete",
            Description: "Delete tasks",
        },
    }
}

func (m *TaskModule) Install(ctx context.Context, tenantID string) error {
    // Create database schema
    if err := m.CreateSchema(ctx, tenantID); err != nil {
        return fmt.Errorf("failed to create schema: %v", err)
    }
    
    // Run migrations
    if err := m.MigrateSchema(ctx, tenantID); err != nil {
        return fmt.Errorf("failed to migrate schema: %v", err)
    }
    
    // Create default permissions
    if err := m.createDefaultPermissions(ctx, tenantID); err != nil {
        return fmt.Errorf("failed to create permissions: %v", err)
    }
    
    return nil
}

func (m *TaskModule) Uninstall(ctx context.Context, tenantID string) error {
    // Drop database schema
    return m.DropSchema(ctx, tenantID)
}

func (m *TaskModule) Enable(ctx context.Context, tenantID string) error {
    // Enable module-specific features
    return m.updateModuleStatus(ctx, tenantID, true)
}

func (m *TaskModule) Disable(ctx context.Context, tenantID string) error {
    // Disable module-specific features
    return m.updateModuleStatus(ctx, tenantID, false)
}

func (m *TaskModule) CreateSchema(ctx context.Context, tenantID string) error {
    // Implementation to create database tables
    return nil
}

func (m *TaskModule) DropSchema(ctx context.Context, tenantID string) error {
    // Implementation to drop database tables
    return nil
}

func (m *TaskModule) MigrateSchema(ctx context.Context, tenantID string) error {
    // Run database migrations
    return nil
}

func (m *TaskModule) HealthCheck(ctx context.Context, tenantID string) error {
    // Check if module is working correctly
    return nil
}

// Helper methods
func (m *TaskModule) createDefaultPermissions(ctx context.Context, tenantID string) error {
    // Create default permissions for module
    return nil
}

func (m *TaskModule) updateModuleStatus(ctx context.Context, tenantID string, enabled bool) error {
    // Update module status in database
    return nil
}
```

### 3.3 Database Models

**models/task.go:**
```go
package models

import (
    "time"
    
    "gorm.io/gorm"
)

type Task struct {
    ID          string     `gorm:"primaryKey" json:"id"`
    TenantID    string     `gorm:"not null;index" json:"tenant_id"`
    Title       string     `gorm:"not null" json:"title"`
    Description *string    `json:"description"`
    Status      TaskStatus `gorm:"not null;default:'TODO'" json:"status"`
    Priority    TaskPriority `gorm:"not null;default:'MEDIUM'" json:"priority"`
    AssignedTo  *string    `gorm:"index" json:"assigned_to"`
    DueDate     *time.Time `json:"due_date"`
    CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
    UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
    
    // Relationships
    Assignee *User `gorm:"foreignKey:AssignedTo" json:"assignee,omitempty"`
}

type TaskStatus string

const (
    TaskStatusTodo       TaskStatus = "TODO"
    TaskStatusInProgress TaskStatus = "IN_PROGRESS"
    TaskStatusCompleted  TaskStatus = "COMPLETED"
    TaskStatusCancelled  TaskStatus = "CANCELLED"
)

type TaskPriority string

const (
    TaskPriorityLow    TaskPriority = "LOW"
    TaskPriorityMedium TaskPriority = "MEDIUM"
    TaskPriorityHigh   TaskPriority = "HIGH"
    TaskPriorityUrgent TaskPriority = "URGENT"
)

func (t *Task) BeforeCreate(tx *gorm.DB) error {
    if t.ID == "" {
        t.ID = generateUUID()
    }
    return nil
}

func (t *Task) TableName() string {
    return "tasks"
}
```

### 3.4 Business Logic Services

**services/task_service.go:**
```go
package services

import (
    "context"
    "fmt"
    
    "github.com/zplus/saas/apps/backend/mymodule/models"
    "gorm.io/gorm"
)

type TaskService struct {
    db *gorm.DB
}

func NewTaskService(db *gorm.DB) *TaskService {
    return &TaskService{db: db}
}

func (s *TaskService) GetTasks(ctx context.Context, tenantID string, filter *TaskFilter) ([]*models.Task, error) {
    var tasks []*models.Task
    
    query := s.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
    
    if filter != nil {
        if filter.Status != nil {
            query = query.Where("status = ?", *filter.Status)
        }
        if filter.Priority != nil {
            query = query.Where("priority = ?", *filter.Priority)
        }
        if filter.AssignedTo != nil {
            query = query.Where("assigned_to = ?", *filter.AssignedTo)
        }
    }
    
    if err := query.Preload("Assignee").Find(&tasks).Error; err != nil {
        return nil, fmt.Errorf("failed to get tasks: %v", err)
    }
    
    return tasks, nil
}

func (s *TaskService) GetTask(ctx context.Context, tenantID, taskID string) (*models.Task, error) {
    var task models.Task
    
    if err := s.db.WithContext(ctx).
        Where("id = ? AND tenant_id = ?", taskID, tenantID).
        Preload("Assignee").
        First(&task).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, fmt.Errorf("task not found")
        }
        return nil, fmt.Errorf("failed to get task: %v", err)
    }
    
    return &task, nil
}

func (s *TaskService) CreateTask(ctx context.Context, tenantID string, input *CreateTaskInput) (*models.Task, error) {
    task := &models.Task{
        TenantID:    tenantID,
        Title:       input.Title,
        Description: input.Description,
        Priority:    input.Priority,
        AssignedTo:  input.AssignedTo,
        DueDate:     input.DueDate,
        Status:      models.TaskStatusTodo,
    }
    
    if err := s.db.WithContext(ctx).Create(task).Error; err != nil {
        return nil, fmt.Errorf("failed to create task: %v", err)
    }
    
    // Preload relationships
    if err := s.db.WithContext(ctx).Preload("Assignee").First(task, task.ID).Error; err != nil {
        return nil, fmt.Errorf("failed to load task relationships: %v", err)
    }
    
    return task, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, tenantID, taskID string, input *UpdateTaskInput) (*models.Task, error) {
    var task models.Task
    
    if err := s.db.WithContext(ctx).
        Where("id = ? AND tenant_id = ?", taskID, tenantID).
        First(&task).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, fmt.Errorf("task not found")
        }
        return nil, fmt.Errorf("failed to get task: %v", err)
    }
    
    // Update fields
    updates := make(map[string]interface{})
    if input.Title != nil {
        updates["title"] = *input.Title
    }
    if input.Description != nil {
        updates["description"] = *input.Description
    }
    if input.Status != nil {
        updates["status"] = *input.Status
    }
    if input.Priority != nil {
        updates["priority"] = *input.Priority
    }
    if input.AssignedTo != nil {
        updates["assigned_to"] = *input.AssignedTo
    }
    if input.DueDate != nil {
        updates["due_date"] = *input.DueDate
    }
    
    if err := s.db.WithContext(ctx).Model(&task).Updates(updates).Error; err != nil {
        return nil, fmt.Errorf("failed to update task: %v", err)
    }
    
    // Preload relationships
    if err := s.db.WithContext(ctx).Preload("Assignee").First(&task, task.ID).Error; err != nil {
        return nil, fmt.Errorf("failed to load task relationships: %v", err)
    }
    
    return &task, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, tenantID, taskID string) error {
    result := s.db.WithContext(ctx).
        Where("id = ? AND tenant_id = ?", taskID, tenantID).
        Delete(&models.Task{})
    
    if result.Error != nil {
        return fmt.Errorf("failed to delete task: %v", result.Error)
    }
    
    if result.RowsAffected == 0 {
        return fmt.Errorf("task not found")
    }
    
    return nil
}

// Input types
type TaskFilter struct {
    Status     *models.TaskStatus   `json:"status"`
    Priority   *models.TaskPriority `json:"priority"`
    AssignedTo *string              `json:"assigned_to"`
}

type CreateTaskInput struct {
    Title       string                `json:"title" validate:"required,min=1,max=255"`
    Description *string               `json:"description"`
    Priority    models.TaskPriority   `json:"priority" validate:"required"`
    AssignedTo  *string               `json:"assigned_to"`
    DueDate     *time.Time            `json:"due_date"`
}

type UpdateTaskInput struct {
    Title       *string               `json:"title"`
    Description *string               `json:"description"`
    Status      *models.TaskStatus    `json:"status"`
    Priority    *models.TaskPriority  `json:"priority"`
    AssignedTo  *string               `json:"assigned_to"`
    DueDate     *time.Time            `json:"due_date"`
}
```

### 3.5 API Handlers

**handlers/rest.go:**
```go
package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/zplus/saas/apps/backend/mymodule/services"
)

type TaskHandler struct {
    service *services.TaskService
}

func NewTaskHandler(service *services.TaskService) *TaskHandler {
    return &TaskHandler{service: service}
}

func (h *TaskHandler) GetTasks(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    
    // Parse query parameters
    var filter services.TaskFilter
    if err := c.QueryParser(&filter); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid query parameters",
        })
    }
    
    tasks, err := h.service.GetTasks(c.Context(), tenantID, &filter)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.JSON(fiber.Map{
        "success": true,
        "data":    tasks,
    })
}

func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    
    var input services.CreateTaskInput
    if err := c.BodyParser(&input); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
    
    // Validate input
    if err := validateStruct(&input); err != nil {
        return c.Status(422).JSON(fiber.Map{
            "error": "Validation failed",
            "details": err.Error(),
        })
    }
    
    task, err := h.service.CreateTask(c.Context(), tenantID, &input)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.Status(201).JSON(fiber.Map{
        "success": true,
        "data":    task,
    })
}

func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    taskID := c.Params("id")
    
    var input services.UpdateTaskInput
    if err := c.BodyParser(&input); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
    
    task, err := h.service.UpdateTask(c.Context(), tenantID, taskID, &input)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.JSON(fiber.Map{
        "success": true,
        "data":    task,
    })
}

func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    taskID := c.Params("id")
    
    if err := h.service.DeleteTask(c.Context(), tenantID, taskID); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.JSON(fiber.Map{
        "success": true,
        "message": "Task deleted successfully",
    })
}
```

## 4. Database Migrations

### 4.1 Migration Files

**migrations/001_initial.up.sql:**
```sql
-- Create tasks table
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'TODO',
    priority VARCHAR(20) NOT NULL DEFAULT 'MEDIUM',
    assigned_to UUID,
    due_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_tasks_tenant_id ON tasks(tenant_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_priority ON tasks(priority);
CREATE INDEX idx_tasks_assigned_to ON tasks(assigned_to);
CREATE INDEX idx_tasks_due_date ON tasks(due_date);

-- Create constraints
ALTER TABLE tasks ADD CONSTRAINT chk_tasks_status 
    CHECK (status IN ('TODO', 'IN_PROGRESS', 'COMPLETED', 'CANCELLED'));

ALTER TABLE tasks ADD CONSTRAINT chk_tasks_priority 
    CHECK (priority IN ('LOW', 'MEDIUM', 'HIGH', 'URGENT'));

-- Create triggers
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_tasks_updated_at BEFORE UPDATE ON tasks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

**migrations/001_initial.down.sql:**
```sql
-- Drop triggers
DROP TRIGGER IF EXISTS update_tasks_updated_at ON tasks;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop table
DROP TABLE IF EXISTS tasks;
```

### 4.2 Migration Management

**Migration Runner:**
```go
package migrations

import (
    "context"
    "fmt"
    "io/ioutil"
    "path/filepath"
    
    "gorm.io/gorm"
)

type MigrationRunner struct {
    db            *gorm.DB
    migrationsDir string
}

func NewMigrationRunner(db *gorm.DB, migrationsDir string) *MigrationRunner {
    return &MigrationRunner{
        db:            db,
        migrationsDir: migrationsDir,
    }
}

func (r *MigrationRunner) RunMigrations(ctx context.Context, tenantID string) error {
    // Get all migration files
    files, err := filepath.Glob(filepath.Join(r.migrationsDir, "*.up.sql"))
    if err != nil {
        return fmt.Errorf("failed to find migration files: %v", err)
    }
    
    for _, file := range files {
        if err := r.runMigrationFile(ctx, file); err != nil {
            return fmt.Errorf("failed to run migration %s: %v", file, err)
        }
    }
    
    return nil
}

func (r *MigrationRunner) runMigrationFile(ctx context.Context, filename string) error {
    content, err := ioutil.ReadFile(filename)
    if err != nil {
        return err
    }
    
    return r.db.WithContext(ctx).Exec(string(content)).Error
}
```

## 5. Testing Module

### 5.1 Unit Tests

**tests/module_test.go:**
```go
package tests

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/zplus/saas/apps/backend/mymodule"
)

func TestTaskModule_Implementation(t *testing.T) {
    module := main.NewTaskModule()
    
    // Test basic properties
    assert.Equal(t, "tasks", module.GetName())
    assert.Equal(t, "1.0.0", module.GetVersion())
    assert.Equal(t, "Task Management", module.GetDisplayName())
    
    // Test permissions
    permissions := module.GetPermissions()
    assert.Len(t, permissions, 3)
    
    expectedPermissions := []string{"tasks:read", "tasks:write", "tasks:delete"}
    for i, perm := range permissions {
        assert.Equal(t, expectedPermissions[i], perm.Name)
    }
    
    // Test routes
    routes := module.GetRoutes()
    assert.Len(t, routes, 4)
    
    // Test default config
    config := module.GetDefaultConfig()
    assert.Contains(t, config, "max_tasks_per_user")
    assert.Equal(t, 100, int(config["max_tasks_per_user"].(float64)))
}

func TestTaskModule_ConfigValidation(t *testing.T) {
    module := main.NewTaskModule()
    
    // Valid config
    validConfig := map[string]interface{}{
        "max_tasks_per_user": 50,
    }
    assert.NoError(t, module.ValidateConfig(validConfig))
    
    // Invalid config - too high
    invalidConfig := map[string]interface{}{
        "max_tasks_per_user": 2000,
    }
    assert.Error(t, module.ValidateConfig(invalidConfig))
    
    // Invalid config - too low
    invalidConfig2 := map[string]interface{}{
        "max_tasks_per_user": 0,
    }
    assert.Error(t, module.ValidateConfig(invalidConfig2))
}
```

### 5.2 Integration Tests

**tests/integration_test.go:**
```go
package tests

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "github.com/zplus/saas/apps/backend/mymodule/services"
    "gorm.io/gorm"
)

type TaskServiceTestSuite struct {
    suite.Suite
    db      *gorm.DB
    service *services.TaskService
    ctx     context.Context
    tenantID string
}

func (suite *TaskServiceTestSuite) SetupTest() {
    // Setup test database
    suite.db = setupTestDB()
    suite.service = services.NewTaskService(suite.db)
    suite.ctx = context.Background()
    suite.tenantID = "test-tenant-123"
    
    // Create test tables
    suite.db.AutoMigrate(&models.Task{})
}

func (suite *TaskServiceTestSuite) TearDownTest() {
    // Clean up test data
    suite.db.Exec("DELETE FROM tasks WHERE tenant_id = ?", suite.tenantID)
}

func (suite *TaskServiceTestSuite) TestCreateTask() {
    input := &services.CreateTaskInput{
        Title:    "Test Task",
        Priority: models.TaskPriorityHigh,
    }
    
    task, err := suite.service.CreateTask(suite.ctx, suite.tenantID, input)
    
    assert.NoError(suite.T(), err)
    assert.NotNil(suite.T(), task)
    assert.Equal(suite.T(), "Test Task", task.Title)
    assert.Equal(suite.T(), models.TaskPriorityHigh, task.Priority)
    assert.Equal(suite.T(), models.TaskStatusTodo, task.Status)
}

func (suite *TaskServiceTestSuite) TestGetTasks() {
    // Create test tasks
    task1 := createTestTask(suite.db, suite.tenantID, "Task 1", models.TaskStatusTodo)
    task2 := createTestTask(suite.db, suite.tenantID, "Task 2", models.TaskStatusCompleted)
    
    // Get all tasks
    tasks, err := suite.service.GetTasks(suite.ctx, suite.tenantID, nil)
    assert.NoError(suite.T(), err)
    assert.Len(suite.T(), tasks, 2)
    
    // Filter by status
    filter := &services.TaskFilter{
        Status: &models.TaskStatusTodo,
    }
    tasks, err = suite.service.GetTasks(suite.ctx, suite.tenantID, filter)
    assert.NoError(suite.T(), err)
    assert.Len(suite.T(), tasks, 1)
    assert.Equal(suite.T(), task1.ID, tasks[0].ID)
}

func (suite *TaskServiceTestSuite) TestUpdateTask() {
    // Create test task
    task := createTestTask(suite.db, suite.tenantID, "Original Title", models.TaskStatusTodo)
    
    // Update task
    newTitle := "Updated Title"
    newStatus := models.TaskStatusInProgress
    input := &services.UpdateTaskInput{
        Title:  &newTitle,
        Status: &newStatus,
    }
    
    updatedTask, err := suite.service.UpdateTask(suite.ctx, suite.tenantID, task.ID, input)
    
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), "Updated Title", updatedTask.Title)
    assert.Equal(suite.T(), models.TaskStatusInProgress, updatedTask.Status)
}

func (suite *TaskServiceTestSuite) TestDeleteTask() {
    // Create test task
    task := createTestTask(suite.db, suite.tenantID, "Task to Delete", models.TaskStatusTodo)
    
    // Delete task
    err := suite.service.DeleteTask(suite.ctx, suite.tenantID, task.ID)
    assert.NoError(suite.T(), err)
    
    // Verify task is deleted
    _, err = suite.service.GetTask(suite.ctx, suite.tenantID, task.ID)
    assert.Error(suite.T(), err)
}

func TestTaskServiceTestSuite(t *testing.T) {
    suite.Run(t, new(TaskServiceTestSuite))
}

func createTestTask(db *gorm.DB, tenantID, title string, status models.TaskStatus) *models.Task {
    task := &models.Task{
        TenantID: tenantID,
        Title:    title,
        Status:   status,
        Priority: models.TaskPriorityMedium,
    }
    
    db.Create(task)
    return task
}
```

## 6. Frontend Integration

### 6.1 React Components

**components/TaskList.tsx:**
```typescript
import React, { useState, useEffect } from 'react';
import { useQuery, useMutation } from '@apollo/client';
import { GET_TASKS, CREATE_TASK, UPDATE_TASK, DELETE_TASK } from '../graphql/tasks';

interface Task {
  id: string;
  title: string;
  description?: string;
  status: 'TODO' | 'IN_PROGRESS' | 'COMPLETED' | 'CANCELLED';
  priority: 'LOW' | 'MEDIUM' | 'HIGH' | 'URGENT';
  assignedTo?: string;
  dueDate?: string;
  createdAt: string;
  updatedAt: string;
}

const TaskList: React.FC = () => {
  const [filter, setFilter] = useState({});
  
  const { data, loading, error, refetch } = useQuery(GET_TASKS, {
    variables: { filter },
  });
  
  const [createTask] = useMutation(CREATE_TASK, {
    onCompleted: () => refetch(),
  });
  
  const [updateTask] = useMutation(UPDATE_TASK, {
    onCompleted: () => refetch(),
  });
  
  const [deleteTask] = useMutation(DELETE_TASK, {
    onCompleted: () => refetch(),
  });
  
  if (loading) return <div>Loading tasks...</div>;
  if (error) return <div>Error: {error.message}</div>;
  
  const handleCreateTask = async (taskData: Partial<Task>) => {
    try {
      await createTask({
        variables: {
          input: taskData,
        },
      });
    } catch (err) {
      console.error('Failed to create task:', err);
    }
  };
  
  const handleUpdateTask = async (id: string, updates: Partial<Task>) => {
    try {
      await updateTask({
        variables: {
          id,
          input: updates,
        },
      });
    } catch (err) {
      console.error('Failed to update task:', err);
    }
  };
  
  const handleDeleteTask = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this task?')) {
      try {
        await deleteTask({
          variables: { id },
        });
      } catch (err) {
        console.error('Failed to delete task:', err);
      }
    }
  };
  
  return (
    <div className="task-list">
      <div className="task-list-header">
        <h2>Tasks</h2>
        <button onClick={() => handleCreateTask({ title: 'New Task', priority: 'MEDIUM' })}>
          Add Task
        </button>
      </div>
      
      <div className="task-filters">
        <select onChange={(e) => setFilter({ ...filter, status: e.target.value })}>
          <option value="">All Status</option>
          <option value="TODO">Todo</option>
          <option value="IN_PROGRESS">In Progress</option>
          <option value="COMPLETED">Completed</option>
        </select>
        
        <select onChange={(e) => setFilter({ ...filter, priority: e.target.value })}>
          <option value="">All Priority</option>
          <option value="LOW">Low</option>
          <option value="MEDIUM">Medium</option>
          <option value="HIGH">High</option>
          <option value="URGENT">Urgent</option>
        </select>
      </div>
      
      <div className="task-items">
        {data?.tasks?.map((task: Task) => (
          <TaskItem
            key={task.id}
            task={task}
            onUpdate={(updates) => handleUpdateTask(task.id, updates)}
            onDelete={() => handleDeleteTask(task.id)}
          />
        ))}
      </div>
    </div>
  );
};

interface TaskItemProps {
  task: Task;
  onUpdate: (updates: Partial<Task>) => void;
  onDelete: () => void;
}

const TaskItem: React.FC<TaskItemProps> = ({ task, onUpdate, onDelete }) => {
  const [isEditing, setIsEditing] = useState(false);
  const [editData, setEditData] = useState(task);
  
  const handleSave = () => {
    onUpdate(editData);
    setIsEditing(false);
  };
  
  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'URGENT': return 'red';
      case 'HIGH': return 'orange';
      case 'MEDIUM': return 'yellow';
      case 'LOW': return 'green';
      default: return 'gray';
    }
  };
  
  return (
    <div className={`task-item status-${task.status.toLowerCase()}`}>
      {isEditing ? (
        <div className="task-edit">
          <input
            value={editData.title}
            onChange={(e) => setEditData({ ...editData, title: e.target.value })}
          />
          <textarea
            value={editData.description || ''}
            onChange={(e) => setEditData({ ...editData, description: e.target.value })}
          />
          <select
            value={editData.status}
            onChange={(e) => setEditData({ ...editData, status: e.target.value as any })}
          >
            <option value="TODO">Todo</option>
            <option value="IN_PROGRESS">In Progress</option>
            <option value="COMPLETED">Completed</option>
            <option value="CANCELLED">Cancelled</option>
          </select>
          <button onClick={handleSave}>Save</button>
          <button onClick={() => setIsEditing(false)}>Cancel</button>
        </div>
      ) : (
        <div className="task-view">
          <div className="task-header">
            <h3>{task.title}</h3>
            <span 
              className="priority-badge" 
              style={{ backgroundColor: getPriorityColor(task.priority) }}
            >
              {task.priority}
            </span>
          </div>
          {task.description && <p>{task.description}</p>}
          <div className="task-meta">
            <span>Status: {task.status}</span>
            {task.dueDate && <span>Due: {new Date(task.dueDate).toLocaleDateString()}</span>}
          </div>
          <div className="task-actions">
            <button onClick={() => setIsEditing(true)}>Edit</button>
            <button onClick={onDelete}>Delete</button>
          </div>
        </div>
      )}
    </div>
  );
};

export default TaskList;
```

### 6.2 GraphQL Queries

**graphql/tasks.ts:**
```typescript
import { gql } from '@apollo/client';

export const GET_TASKS = gql`
  query GetTasks($filter: TaskFilter) {
    tasks(filter: $filter) {
      id
      title
      description
      status
      priority
      assignedTo
      dueDate
      createdAt
      updatedAt
      assignee {
        id
        name
        email
      }
    }
  }
`;

export const GET_TASK = gql`
  query GetTask($id: ID!) {
    task(id: $id) {
      id
      title
      description
      status
      priority
      assignedTo
      dueDate
      createdAt
      updatedAt
      assignee {
        id
        name
        email
      }
    }
  }
`;

export const CREATE_TASK = gql`
  mutation CreateTask($input: CreateTaskInput!) {
    createTask(input: $input) {
      id
      title
      description
      status
      priority
      assignedTo
      dueDate
      createdAt
      updatedAt
    }
  }
`;

export const UPDATE_TASK = gql`
  mutation UpdateTask($id: ID!, $input: UpdateTaskInput!) {
    updateTask(id: $id, input: $input) {
      id
      title
      description
      status
      priority
      assignedTo
      dueDate
      createdAt
      updatedAt
    }
  }
`;

export const DELETE_TASK = gql`
  mutation DeleteTask($id: ID!) {
    deleteTask(id: $id)
  }
`;
```

## 7. Module Deployment

### 7.1 Docker Configuration

**Dockerfile:**
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./apps/backend/mymodule

# Production stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/apps/backend/mymodule/migrations ./migrations

EXPOSE 8080

CMD ["./main"]
```

### 7.2 Kubernetes Deployment

**k8s/task-module.yaml:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: task-module
  namespace: zplus-prod
spec:
  replicas: 2
  selector:
    matchLabels:
      app: task-module
  template:
    metadata:
      labels:
        app: task-module
    spec:
      containers:
      - name: task-module
        image: zplus/task-module:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: zplus-config
              key: DB_HOST
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: zplus-secrets
              key: DB_PASSWORD
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: task-module-service
  namespace: zplus-prod
spec:
  selector:
    app: task-module
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
```

## 8. Best Practices

### 8.1 Module Development Guidelines

1. **Follow Interface Contract**: Implement all required methods
2. **Proper Error Handling**: Return meaningful errors
3. **Logging**: Log important events and errors
4. **Testing**: Comprehensive unit and integration tests
5. **Documentation**: Clear API documentation
6. **Security**: Validate all inputs and implement proper authorization
7. **Performance**: Optimize database queries and use caching
8. **Monitoring**: Include health checks and metrics

### 8.2 Database Best Practices

1. **Use Transactions**: For operations that modify multiple records
2. **Proper Indexing**: Index frequently queried columns
3. **Foreign Key Constraints**: Maintain data integrity
4. **Soft Deletes**: Consider soft deletes for important data
5. **Audit Trail**: Track changes to important entities
6. **Migration Scripts**: Version all database changes

### 8.3 API Best Practices

1. **RESTful Design**: Follow REST conventions
2. **Input Validation**: Validate all inputs
3. **Rate Limiting**: Implement appropriate rate limits
4. **Pagination**: Support pagination for list endpoints
5. **Error Responses**: Consistent error response format
6. **API Versioning**: Version your APIs
7. **Documentation**: Comprehensive API documentation

## 9. Troubleshooting

### 9.1 Common Issues

**Module Registration Failed:**
```bash
# Check if module implements all required methods
go vet ./apps/backend/mymodule/

# Check dependencies
go mod verify
```

**Database Migration Issues:**
```bash
# Check migration syntax
psql -f migrations/001_initial.up.sql --dry-run

# Manual rollback
psql -f migrations/001_initial.down.sql
```

**Permission Denied:**
```bash
# Check if permissions are properly registered
# Verify RBAC configuration
# Check user roles and permissions
```

### 9.2 Debugging

**Enable Debug Logging:**
```go
log.SetLevel(log.DebugLevel)
```

**Database Query Logging:**
```go
db := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})
```

**HTTP Request Logging:**
```go
app.Use(logger.New())
```

## 10. Advanced Topics

### 10.1 Module Dependencies

**Declaring Dependencies:**
```go
func (m *TaskModule) GetDependencies() []string {
    return []string{"users", "notifications"}
}
```

**Circular Dependency Detection:**
```go
func (r *ModuleRegistry) detectCircularDependencies() error {
    // Implementation to detect circular dependencies
    return nil
}
```

### 10.2 Module Communication

**Inter-module Communication:**
```go
type ModuleCommunicator interface {
    SendEvent(event ModuleEvent) error
    Subscribe(eventType string, handler EventHandler) error
}

type ModuleEvent struct {
    Type     string
    Source   string
    Data     interface{}
    TenantID string
}
```

### 10.3 Module Upgrades

**Version Migration:**
```go
func (m *TaskModule) Upgrade(ctx context.Context, tenantID string, fromVersion, toVersion string) error {
    // Implementation for module upgrade
    return nil
}
```

**Backward Compatibility:**
```go
func (m *TaskModule) IsCompatible(version string) bool {
    // Check if this module version is compatible with system version
    return true
}
```
