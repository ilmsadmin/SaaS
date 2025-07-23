-- CRM Service Database Schema
-- Migration: 001_create_crm_tables.sql
-- Description: Creates tables for CRM service (customers, leads, opportunities, activities)

-- Create customers table
CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    company VARCHAR(255),
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100),
    zip_code VARCHAR(20),
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'prospect')),
    source VARCHAR(100),
    tags TEXT, -- Comma-separated tags
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create leads table
CREATE TABLE IF NOT EXISTS leads (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    company VARCHAR(255),
    title VARCHAR(255),
    source VARCHAR(100),
    status VARCHAR(50) DEFAULT 'new' CHECK (status IN ('new', 'qualified', 'contacted', 'converted', 'lost')),
    score INTEGER DEFAULT 0 CHECK (score >= 0 AND score <= 100),
    assigned_to INTEGER, -- User ID from auth service
    value DECIMAL(15,2) DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    converted_at TIMESTAMP WITH TIME ZONE
);

-- Create opportunities table
CREATE TABLE IF NOT EXISTS opportunities (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    customer_id INTEGER NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    value DECIMAL(15,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'USD',
    stage VARCHAR(50) DEFAULT 'prospecting' CHECK (stage IN ('prospecting', 'qualification', 'proposal', 'negotiation', 'closed-won', 'closed-lost')),
    probability INTEGER DEFAULT 10 CHECK (probability >= 0 AND probability <= 100),
    source VARCHAR(100),
    assigned_to INTEGER, -- User ID from auth service
    expected_date DATE NOT NULL,
    closed_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create contact_activities table
CREATE TABLE IF NOT EXISTS contact_activities (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    customer_id INTEGER REFERENCES customers(id) ON DELETE CASCADE,
    lead_id INTEGER REFERENCES leads(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('call', 'email', 'meeting', 'note', 'task')),
    subject VARCHAR(255) NOT NULL,
    description TEXT,
    duration INTEGER DEFAULT 0, -- in minutes
    user_id INTEGER NOT NULL, -- User ID from auth service
    scheduled_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure either customer_id or lead_id is set, but not both
    CONSTRAINT check_contact_reference CHECK (
        (customer_id IS NOT NULL AND lead_id IS NULL) OR 
        (customer_id IS NULL AND lead_id IS NOT NULL)
    )
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_customers_tenant_id ON customers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_customers_email ON customers(tenant_id, email);
CREATE INDEX IF NOT EXISTS idx_customers_status ON customers(tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_customers_created_at ON customers(tenant_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_leads_tenant_id ON leads(tenant_id);
CREATE INDEX IF NOT EXISTS idx_leads_email ON leads(tenant_id, email);
CREATE INDEX IF NOT EXISTS idx_leads_status ON leads(tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_leads_assigned_to ON leads(tenant_id, assigned_to);
CREATE INDEX IF NOT EXISTS idx_leads_created_at ON leads(tenant_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_opportunities_tenant_id ON opportunities(tenant_id);
CREATE INDEX IF NOT EXISTS idx_opportunities_customer_id ON opportunities(tenant_id, customer_id);
CREATE INDEX IF NOT EXISTS idx_opportunities_stage ON opportunities(tenant_id, stage);
CREATE INDEX IF NOT EXISTS idx_opportunities_assigned_to ON opportunities(tenant_id, assigned_to);
CREATE INDEX IF NOT EXISTS idx_opportunities_expected_date ON opportunities(tenant_id, expected_date);
CREATE INDEX IF NOT EXISTS idx_opportunities_created_at ON opportunities(tenant_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_activities_tenant_id ON contact_activities(tenant_id);
CREATE INDEX IF NOT EXISTS idx_activities_customer_id ON contact_activities(tenant_id, customer_id);
CREATE INDEX IF NOT EXISTS idx_activities_lead_id ON contact_activities(tenant_id, lead_id);
CREATE INDEX IF NOT EXISTS idx_activities_user_id ON contact_activities(tenant_id, user_id);
CREATE INDEX IF NOT EXISTS idx_activities_type ON contact_activities(tenant_id, type);
CREATE INDEX IF NOT EXISTS idx_activities_created_at ON contact_activities(tenant_id, created_at DESC);

-- Add triggers to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_customers_updated_at BEFORE UPDATE ON customers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_leads_updated_at BEFORE UPDATE ON leads
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_opportunities_updated_at BEFORE UPDATE ON opportunities
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_activities_updated_at BEFORE UPDATE ON contact_activities
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert sample data for development
INSERT INTO customers (tenant_id, name, email, phone, company, status, source, tags, notes) VALUES
('tenant-1', 'John Doe', 'john.doe@example.com', '+1-555-0123', 'Acme Corp', 'active', 'website', 'vip,enterprise', 'Important client from website signup'),
('tenant-1', 'Jane Smith', 'jane.smith@techcorp.com', '+1-555-0124', 'TechCorp', 'active', 'referral', 'tech,startup', 'Referred by existing client'),
('tenant-1', 'Mike Johnson', 'mike@startup.io', '+1-555-0125', 'Startup Inc', 'prospect', 'social', 'startup,saas', 'Found us on LinkedIn')
ON CONFLICT DO NOTHING;

INSERT INTO leads (tenant_id, name, email, phone, company, title, source, status, score, assigned_to, value, notes) VALUES
('tenant-1', 'Sarah Wilson', 'sarah.wilson@bigcorp.com', '+1-555-0126', 'BigCorp', 'Marketing Director', 'website', 'qualified', 75, 1, 25000.00, 'Interested in enterprise package'),
('tenant-1', 'Robert Brown', 'robert.brown@smallbiz.com', '+1-555-0127', 'Small Biz', 'Owner', 'cold-call', 'new', 25, 1, 5000.00, 'Initial contact made'),
('tenant-1', 'Lisa Davis', 'lisa.davis@mediumco.com', '+1-555-0128', 'Medium Co', 'IT Manager', 'referral', 'contacted', 60, 1, 15000.00, 'Had demo call, interested')
ON CONFLICT DO NOTHING;

INSERT INTO opportunities (tenant_id, customer_id, name, description, value, currency, stage, probability, source, assigned_to, expected_date) VALUES
('tenant-1', 1, 'Acme Corp Expansion', 'Expand services to include additional modules', 50000.00, 'USD', 'proposal', 70, 'existing-client', 1, '2025-08-15'),
('tenant-1', 2, 'TechCorp Implementation', 'Full CRM implementation for TechCorp', 35000.00, 'USD', 'negotiation', 80, 'referral', 1, '2025-07-30'),
('tenant-1', 1, 'Acme Training Package', 'Staff training and onboarding', 10000.00, 'USD', 'closed-won', 100, 'existing-client', 1, '2025-07-25')
ON CONFLICT DO NOTHING;
