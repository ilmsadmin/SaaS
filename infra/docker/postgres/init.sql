-- Initialize Zplus SaaS Database
-- This script creates the necessary database structure for the multi-tenant platform

-- Create main application database if not exists
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'zplus_saas') THEN
        CREATE DATABASE zplus_saas;
    END IF;
END
$$;

-- Connect to the main database
\c zplus_saas;

-- Enable necessary extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Create basic schemas for multi-tenancy
CREATE SCHEMA IF NOT EXISTS public;
CREATE SCHEMA IF NOT EXISTS system;
CREATE SCHEMA IF NOT EXISTS shared;

-- System tables for multi-tenancy management
CREATE TABLE IF NOT EXISTS system.tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    subdomain VARCHAR(100) UNIQUE NOT NULL,
    custom_domain VARCHAR(255),
    schema_name VARCHAR(100) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    plan VARCHAR(100) DEFAULT 'basic',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_tenants_subdomain ON system.tenants(subdomain);
CREATE INDEX IF NOT EXISTS idx_tenants_schema ON system.tenants(schema_name);
CREATE INDEX IF NOT EXISTS idx_tenants_status ON system.tenants(status);

-- Insert default system tenant
INSERT INTO system.tenants (name, subdomain, schema_name, status, plan) 
VALUES ('System Admin', 'admin', 'system', 'active', 'enterprise')
ON CONFLICT (subdomain) DO NOTHING;

-- Create trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_tenants_updated_at 
    BEFORE UPDATE ON system.tenants 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Success message
SELECT 'Zplus SaaS Database initialized successfully!' as status;
