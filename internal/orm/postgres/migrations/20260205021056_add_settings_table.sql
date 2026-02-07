-- +goose Up
-- +goose StatementBegin

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum type for settings.type
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

-- Create settings table
CREATE TABLE IF NOT EXISTS settings (
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

    CONSTRAINT fk_settings_account
        FOREIGN KEY (account_id)
        REFERENCES accounts(id)
        ON DELETE CASCADE
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_settings_is_active
    ON settings (is_active);

CREATE INDEX IF NOT EXISTS idx_settings_account_id
    ON settings (account_id);

CREATE INDEX IF NOT EXISTS idx_settings_key
    ON settings (key);

DROP TRIGGER IF EXISTS update_settings_updated_at ON settings;
CREATE TRIGGER update_settings_updated_at
BEFORE UPDATE ON settings
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_settings_updated_at ON settings;
DROP TABLE IF EXISTS settings;
DROP TYPE IF EXISTS setting_type;
-- +goose StatementEnd
