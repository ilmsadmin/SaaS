-- Module system tables
CREATE TABLE IF NOT EXISTS modules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    display_name VARCHAR(100) NOT NULL,
    description TEXT,
    version VARCHAR(20) NOT NULL,
    category VARCHAR(50) NOT NULL,
    icon TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    dependencies TEXT NOT NULL DEFAULT '[]', -- JSON array
    permissions TEXT NOT NULL DEFAULT '[]', -- JSON array
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_modules_category ON modules(category);
CREATE INDEX idx_modules_is_active ON modules(is_active);

-- Module dependencies
CREATE TABLE IF NOT EXISTS module_dependencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    module_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    depends_on_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    min_version VARCHAR(20) NOT NULL,
    is_required BOOLEAN NOT NULL DEFAULT true,
    conflicts_with BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    UNIQUE(module_id, depends_on_id)
);

CREATE INDEX idx_module_dependencies_module_id ON module_dependencies(module_id);
CREATE INDEX idx_module_dependencies_depends_on_id ON module_dependencies(depends_on_id);

-- Module installations for tenants
CREATE TABLE IF NOT EXISTS module_installations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    module_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    version VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'installing',
    config TEXT NOT NULL DEFAULT '{}', -- JSON config
    install_data TEXT NOT NULL DEFAULT '{}', -- JSON data
    error_message TEXT,
    installed_by UUID NOT NULL,
    installed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    uninstalled_at TIMESTAMPTZ,
    
    CONSTRAINT chk_installation_status CHECK (status IN ('installing', 'installed', 'failed', 'updating', 'uninstalling')),
    UNIQUE(tenant_id, module_id, uninstalled_at) -- Allow reinstalling
);

CREATE INDEX idx_module_installations_tenant_id ON module_installations(tenant_id);
CREATE INDEX idx_module_installations_module_id ON module_installations(module_id);
CREATE INDEX idx_module_installations_status ON module_installations(status);

-- Module permissions
CREATE TABLE IF NOT EXISTS module_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    module_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    permission VARCHAR(100) NOT NULL,
    description VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    UNIQUE(module_id, permission)
);

CREATE INDEX idx_module_permissions_module_id ON module_permissions(module_id);
CREATE INDEX idx_module_permissions_category ON module_permissions(category);

-- Tenant configurations
CREATE TABLE IF NOT EXISTS tenant_configurations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE UNIQUE,
    custom_domain VARCHAR(255),
    ssl_enabled BOOLEAN NOT NULL DEFAULT false,
    custom_css TEXT,
    custom_javascript TEXT,
    branding_config JSONB NOT NULL DEFAULT '{}',
    security_config JSONB NOT NULL DEFAULT '{}',
    notification_config JSONB NOT NULL DEFAULT '{}',
    integration_config JSONB NOT NULL DEFAULT '{}',
    feature_flags JSONB NOT NULL DEFAULT '{}',
    data_retention_days INTEGER NOT NULL DEFAULT 365,
    allowed_ips JSONB,
    two_factor_required BOOLEAN NOT NULL DEFAULT false,
    password_policy JSONB NOT NULL DEFAULT '{"minLength": 8, "requireUppercase": true, "requireLowercase": true, "requireNumbers": true, "requireSpecialChars": false}',
    session_timeout_mins INTEGER NOT NULL DEFAULT 480,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tenant_configurations_tenant_id ON tenant_configurations(tenant_id);
CREATE INDEX idx_tenant_configurations_custom_domain ON tenant_configurations(custom_domain);

-- Module dependencies
CREATE TABLE IF NOT EXISTS module_dependencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    module_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    depends_on_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    min_version VARCHAR(20) NOT NULL,
    is_required BOOLEAN NOT NULL DEFAULT true,
    conflicts_with BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    UNIQUE(module_id, depends_on_id)
);

CREATE INDEX idx_module_dependencies_module_id ON module_dependencies(module_id);
CREATE INDEX idx_module_dependencies_depends_on_id ON module_dependencies(depends_on_id);

-- Module marketplace
CREATE TABLE IF NOT EXISTS module_marketplace (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    module_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE UNIQUE,
    publisher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    is_free BOOLEAN NOT NULL DEFAULT true,
    rating DECIMAL(3,2) NOT NULL DEFAULT 0.00,
    total_downloads INTEGER NOT NULL DEFAULT 0,
    screenshots JSONB NOT NULL DEFAULT '[]',
    documentation TEXT,
    support_email VARCHAR(255),
    homepage VARCHAR(255),
    repository VARCHAR(255),
    license VARCHAR(100) NOT NULL DEFAULT 'MIT',
    tags JSONB NOT NULL DEFAULT '[]',
    is_verified BOOLEAN NOT NULL DEFAULT false,
    is_featured BOOLEAN NOT NULL DEFAULT false,
    published_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_module_marketplace_module_id ON module_marketplace(module_id);
CREATE INDEX idx_module_marketplace_publisher_id ON module_marketplace(publisher_id);
CREATE INDEX idx_module_marketplace_is_verified ON module_marketplace(is_verified);
CREATE INDEX idx_module_marketplace_is_featured ON module_marketplace(is_featured);

-- Module installations
CREATE TABLE IF NOT EXISTS module_installations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    module_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    version VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'installing',
    config JSONB NOT NULL DEFAULT '{}',
    install_data JSONB NOT NULL DEFAULT '{}',
    error_message TEXT,
    installed_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    installed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    uninstalled_at TIMESTAMPTZ,
    
    CONSTRAINT chk_installation_status CHECK (status IN ('installing', 'installed', 'failed', 'updating', 'uninstalling'))
);

CREATE INDEX idx_module_installations_tenant_id ON module_installations(tenant_id);
CREATE INDEX idx_module_installations_module_id ON module_installations(module_id);
CREATE INDEX idx_module_installations_status ON module_installations(status);

-- Module permissions
CREATE TABLE IF NOT EXISTS module_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    module_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    permission VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    category VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    UNIQUE(module_id, permission)
);

CREATE INDEX idx_module_permissions_module_id ON module_permissions(module_id);
CREATE INDEX idx_module_permissions_category ON module_permissions(category);

-- Seed default modules
INSERT INTO modules (id, name, display_name, description, version, category, icon, is_active, dependencies, permissions) VALUES
(gen_random_uuid(), 'crm', 'Customer Relationship Management', 'Manage customers, leads, and sales pipeline', '1.0.0', 'business', 'users', true, '[]', '["crm.read", "crm.write", "crm.delete"]'),
(gen_random_uuid(), 'hrm', 'Human Resource Management', 'Employee management and HR processes', '1.0.0', 'hr', 'user-group', true, '[]', '["hrm.read", "hrm.write", "hrm.delete"]'),
(gen_random_uuid(), 'pos', 'Point of Sale', 'Sales and inventory management system', '1.0.0', 'retail', 'shopping-cart', true, '[]', '["pos.read", "pos.write", "pos.delete"]'),
(gen_random_uuid(), 'lms', 'Learning Management System', 'Course and training management', '1.0.0', 'education', 'academic-cap', true, '[]', '["lms.read", "lms.write", "lms.delete"]'),
(gen_random_uuid(), 'checkin', 'Check-in System', 'Attendance and check-in tracking', '1.0.0', 'attendance', 'clock', true, '["hrm"]', '["checkin.read", "checkin.write"]'),
(gen_random_uuid(), 'payment', 'Payment Processing', 'Payment gateway and transaction management', '1.0.0', 'finance', 'credit-card', true, '[]', '["payment.read", "payment.write"]'),
(gen_random_uuid(), 'accounting', 'Accounting System', 'Financial accounting and reporting', '1.0.0', 'finance', 'calculator', true, '["payment"]', '["accounting.read", "accounting.write"]'),
(gen_random_uuid(), 'ecommerce', 'E-commerce Platform', 'Online store and product management', '1.0.0', 'retail', 'shopping-bag', true, '["pos", "payment"]', '["ecommerce.read", "ecommerce.write", "ecommerce.delete"]');

-- Add unique constraint for tenant_modules if not exists
ALTER TABLE tenant_modules ADD CONSTRAINT IF NOT EXISTS unique_tenant_module UNIQUE (tenant_id, module_id);
