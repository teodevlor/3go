-- +goose Up
-- +goose StatementBegin

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum type for system_settings.type
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'setting_type') THEN
        CREATE TYPE setting_type AS ENUM (
            'web_system',
            'app_user',
            'app_driver'
        );
    END IF;
END;
$$;

-- Create system_settings table
CREATE TABLE IF NOT EXISTS system_settings (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    account_id uuid NOT NULL,

    key varchar(255) NOT NULL,
    value jsonb,
    "type" setting_type NOT NULL,
    description text,
    is_active boolean DEFAULT true,

    metadata jsonb DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_system_settings_account
        FOREIGN KEY (account_id)
        REFERENCES accounts(id)
        ON DELETE CASCADE
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_system_settings_is_active
    ON system_settings (is_active);

CREATE INDEX IF NOT EXISTS idx_system_settings_account_id
    ON system_settings (account_id);

CREATE INDEX IF NOT EXISTS idx_system_settings_key
    ON system_settings (key);

DROP TRIGGER IF EXISTS update_system_settings_updated_at ON system_settings;
CREATE TRIGGER update_system_settings_updated_at
BEFORE UPDATE ON system_settings
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_system_settings_updated_at ON system_settings;
DROP TABLE IF EXISTS system_settings;
DROP TYPE IF EXISTS setting_type;
-- +goose StatementEnd
