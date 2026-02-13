-- +goose Up
-- +goose StatementBegin

-- Create enum type for system_admins.department
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'admin_department') THEN
        CREATE TYPE admin_department AS ENUM (
            'employee',
            'admin',
            'seller',
            'marketer'
        );
    END IF;
END;
$$;

-- Create system_admins table
CREATE TABLE IF NOT EXISTS system_admins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    full_name VARCHAR(255),
    department admin_department,
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMPTZ,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_system_admins_email
    ON system_admins (email);

CREATE INDEX IF NOT EXISTS idx_system_admins_is_active
    ON system_admins (is_active);

CREATE INDEX IF NOT EXISTS idx_system_admins_department
    ON system_admins (department);

-- Trigger: auto update updated_at
DROP TRIGGER IF EXISTS update_system_admins_updated_at ON system_admins;
CREATE TRIGGER update_system_admins_updated_at
BEFORE UPDATE ON system_admins
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_system_admins_updated_at ON system_admins;

DROP INDEX IF EXISTS idx_system_admins_department;
DROP INDEX IF EXISTS idx_system_admins_is_active;
DROP INDEX IF EXISTS idx_system_admins_email;

DROP TABLE IF EXISTS system_admins;

DROP TYPE IF EXISTS admin_department;

-- +goose StatementEnd
