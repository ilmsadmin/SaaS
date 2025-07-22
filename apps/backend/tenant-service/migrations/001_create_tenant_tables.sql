-- Migration: 001_create_tenant_tables.sql
-- Description: Create tables for tenant management system

-- Create tenants table (system level)
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    subdomain VARCHAR(100) UNIQUE NOT NULL,
    domain VARCHAR(255) UNIQUE,
    logo TEXT,
    status VARCHAR(20) DEFAULT 'trial' CHECK (status IN ('active', 'suspended', 'trial')),
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create plans table
CREATE TABLE IF NOT EXISTS plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    billing_cycle VARCHAR(20) NOT NULL CHECK (billing_cycle IN ('monthly', 'yearly')),
    max_users INTEGER,
    max_storage BIGINT, -- in bytes
    features JSONB DEFAULT '[]',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create subscriptions table
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    plan_id UUID REFERENCES plans(id) ON DELETE RESTRICT,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'cancelled', 'expired')),
    trial_end_at TIMESTAMP WITH TIME ZONE,
    current_period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    current_period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    cancel_at_period_end BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create modules table
CREATE TABLE IF NOT EXISTS modules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    version VARCHAR(20) NOT NULL,
    category VARCHAR(50),
    icon VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    dependencies JSONB DEFAULT '[]',
    permissions JSONB DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create tenant_modules table
CREATE TABLE IF NOT EXISTS tenant_modules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    module_id UUID REFERENCES modules(id) ON DELETE CASCADE,
    is_enabled BOOLEAN DEFAULT true,
    config JSONB DEFAULT '{}',
    installed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(tenant_id, module_id)
);

-- Create payments table
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscription_id UUID REFERENCES subscriptions(id) ON DELETE CASCADE,
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'completed', 'failed', 'refunded')),
    payment_method VARCHAR(50),
    external_payment_id VARCHAR(255),
    paid_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_tenants_subdomain ON tenants(subdomain);
CREATE INDEX IF NOT EXISTS idx_tenants_domain ON tenants(domain);
CREATE INDEX IF NOT EXISTS idx_tenants_status ON tenants(status);

CREATE INDEX IF NOT EXISTS idx_subscriptions_tenant_id ON subscriptions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON subscriptions(status);
CREATE INDEX IF NOT EXISTS idx_subscriptions_expires ON subscriptions(current_period_end);

CREATE INDEX IF NOT EXISTS idx_tenant_modules_tenant ON tenant_modules(tenant_id);
CREATE INDEX IF NOT EXISTS idx_tenant_modules_module ON tenant_modules(module_id);

CREATE INDEX IF NOT EXISTS idx_payments_subscription ON payments(subscription_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_tenants_updated_at BEFORE UPDATE ON tenants
    FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE TRIGGER update_plans_updated_at BEFORE UPDATE ON plans
    FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE TRIGGER update_subscriptions_updated_at BEFORE UPDATE ON subscriptions
    FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE TRIGGER update_modules_updated_at BEFORE UPDATE ON modules
    FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE TRIGGER update_tenant_modules_updated_at BEFORE UPDATE ON tenant_modules
    FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- Insert default plans
INSERT INTO plans (id, name, description, price, currency, billing_cycle, max_users, max_storage, features, is_active) VALUES
(
    '00000000-0000-0000-0000-000000000001',
    'Free',
    'Perfect for getting started',
    0.00,
    'USD',
    'monthly',
    5,
    1073741824, -- 1GB
    '["Basic CRM", "5 Users", "1GB Storage", "Email Support"]',
    true
),
(
    '00000000-0000-0000-0000-000000000002',
    'Starter',
    'Great for small teams',
    29.00,
    'USD',
    'monthly',
    15,
    10737418240, -- 10GB
    '["All CRM Features", "15 Users", "10GB Storage", "Priority Support", "Basic Analytics"]',
    true
),
(
    '00000000-0000-0000-0000-000000000003',
    'Professional',
    'Best for growing businesses',
    99.00,
    'USD',
    'monthly',
    50,
    107374182400, -- 100GB
    '["All Features", "50 Users", "100GB Storage", "24/7 Support", "Advanced Analytics", "API Access"]',
    true
),
(
    '00000000-0000-0000-0000-000000000004',
    'Enterprise',
    'For large organizations',
    299.00,
    'USD',
    'monthly',
    NULL, -- Unlimited
    NULL, -- Unlimited
    '["All Features", "Unlimited Users", "Unlimited Storage", "Dedicated Support", "Custom Integration", "SLA"]',
    true
)
ON CONFLICT (id) DO NOTHING;

-- Insert default modules
INSERT INTO modules (id, name, display_name, description, version, category, icon, is_active) VALUES
(
    '00000000-0000-0000-0000-000000000101',
    'crm',
    'Customer Relationship Management',
    'Manage customers, leads, and sales pipeline',
    '1.0.0',
    'business',
    'users',
    true
),
(
    '00000000-0000-0000-0000-000000000102',
    'hrm',
    'Human Resource Management',
    'Employee management and HR processes',
    '1.0.0',
    'business',
    'user-tie',
    true
),
(
    '00000000-0000-0000-0000-000000000103',
    'pos',
    'Point of Sale',
    'Sales transactions and inventory management',
    '1.0.0',
    'business',
    'cash-register',
    true
),
(
    '00000000-0000-0000-0000-000000000104',
    'lms',
    'Learning Management System',
    'Course creation and student management',
    '1.0.0',
    'education',
    'graduation-cap',
    true
),
(
    '00000000-0000-0000-0000-000000000105',
    'checkin',
    'Check-in Management',
    'Employee attendance and time tracking',
    '1.0.0',
    'business',
    'clock',
    true
),
(
    '00000000-0000-0000-0000-000000000106',
    'payment',
    'Payment Processing',
    'Handle payments and billing',
    '1.0.0',
    'finance',
    'credit-card',
    true
)
ON CONFLICT (id) DO NOTHING;

-- Insert demo tenant for development
INSERT INTO tenants (id, name, subdomain, status) VALUES 
(
    '00000000-0000-0000-0000-000000000201',
    'Demo Company',
    'demo',
    'active'
)
ON CONFLICT (subdomain) DO NOTHING;

-- Create demo subscription
INSERT INTO subscriptions (id, tenant_id, plan_id, status, current_period_start, current_period_end) VALUES
(
    '00000000-0000-0000-0000-000000000301',
    '00000000-0000-0000-0000-000000000201', -- demo tenant
    '00000000-0000-0000-0000-000000000002', -- starter plan
    'active',
    NOW(),
    NOW() + INTERVAL '1 month'
)
ON CONFLICT (id) DO NOTHING;

-- Enable all modules for demo tenant
INSERT INTO tenant_modules (tenant_id, module_id, is_enabled) 
SELECT 
    '00000000-0000-0000-0000-000000000201' as tenant_id,
    id as module_id,
    true as is_enabled
FROM modules
ON CONFLICT (tenant_id, module_id) DO NOTHING;
