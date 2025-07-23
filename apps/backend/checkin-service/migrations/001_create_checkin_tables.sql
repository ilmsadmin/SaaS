-- Migration for Checkin Service
-- Creates tables for attendance tracking and checkin management

-- Create checkin_records table
CREATE TABLE IF NOT EXISTS checkin_records (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    employee_id INTEGER NOT NULL,
    employee_name VARCHAR(255) NOT NULL,
    checkin_type VARCHAR(50) NOT NULL CHECK (checkin_type IN ('checkin', 'checkout', 'break_start', 'break_end')),
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    location VARCHAR(500),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    ip_address INET,
    device_info TEXT,
    photo TEXT, -- URL or base64 encoded photo
    notes TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'approved' CHECK (status IN ('approved', 'pending', 'rejected')),
    approved_by INTEGER,
    approved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create attendance_policies table
CREATE TABLE IF NOT EXISTS attendance_policies (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    work_start_time TIME NOT NULL DEFAULT '09:00:00',
    work_end_time TIME NOT NULL DEFAULT '18:00:00',
    break_duration INTEGER NOT NULL DEFAULT 60, -- in minutes
    late_threshold INTEGER NOT NULL DEFAULT 15, -- in minutes
    early_threshold INTEGER NOT NULL DEFAULT 15, -- in minutes
    require_photo BOOLEAN NOT NULL DEFAULT FALSE,
    require_location BOOLEAN NOT NULL DEFAULT FALSE,
    allowed_radius INTEGER NOT NULL DEFAULT 100, -- in meters
    work_days JSONB NOT NULL DEFAULT '["monday", "tuesday", "wednesday", "thursday", "friday"]',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create attendance_summary table
CREATE TABLE IF NOT EXISTS attendance_summary (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    employee_id INTEGER NOT NULL,
    employee_name VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    checkin_time TIMESTAMP WITH TIME ZONE,
    checkout_time TIMESTAMP WITH TIME ZONE,
    work_hours DECIMAL(5, 2) NOT NULL DEFAULT 0,
    break_hours DECIMAL(5, 2) NOT NULL DEFAULT 0,
    overtime_hours DECIMAL(5, 2) NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'absent' CHECK (status IN ('present', 'absent', 'late', 'early_leave', 'partial')),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(tenant_id, employee_id, date)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_checkin_records_tenant_employee ON checkin_records(tenant_id, employee_id);
CREATE INDEX IF NOT EXISTS idx_checkin_records_tenant_timestamp ON checkin_records(tenant_id, timestamp);
CREATE INDEX IF NOT EXISTS idx_checkin_records_tenant_type ON checkin_records(tenant_id, checkin_type);
CREATE INDEX IF NOT EXISTS idx_checkin_records_tenant_status ON checkin_records(tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_checkin_records_employee_date ON checkin_records(employee_id, DATE(timestamp));

CREATE INDEX IF NOT EXISTS idx_attendance_policies_tenant ON attendance_policies(tenant_id);
CREATE INDEX IF NOT EXISTS idx_attendance_policies_tenant_active ON attendance_policies(tenant_id, is_active);

CREATE INDEX IF NOT EXISTS idx_attendance_summary_tenant_employee ON attendance_summary(tenant_id, employee_id);
CREATE INDEX IF NOT EXISTS idx_attendance_summary_tenant_date ON attendance_summary(tenant_id, date);
CREATE INDEX IF NOT EXISTS idx_attendance_summary_tenant_status ON attendance_summary(tenant_id, status);

-- Create triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_checkin_records_updated_at BEFORE UPDATE ON checkin_records FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_attendance_policies_updated_at BEFORE UPDATE ON attendance_policies FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_attendance_summary_updated_at BEFORE UPDATE ON attendance_summary FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert default attendance policy for demo
INSERT INTO attendance_policies (
    tenant_id, name, work_start_time, work_end_time, break_duration,
    late_threshold, early_threshold, require_photo, require_location,
    allowed_radius, work_days, is_active
) VALUES (
    'demo-tenant', 'Standard Policy', '09:00:00', '18:00:00', 60,
    15, 15, false, false, 100,
    '["monday", "tuesday", "wednesday", "thursday", "friday"]', true
) ON CONFLICT DO NOTHING;
