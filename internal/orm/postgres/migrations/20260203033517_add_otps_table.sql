-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE system_otps (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    target varchar(255) NOT NULL,
    -- email hoặc phone

    otp_code varchar(10) NOT NULL,

    purpose varchar(50) NOT NULL,
    -- register | reset_password | transaction | withdraw | login | verify_phone

    attempt_count int DEFAULT 0,
    max_attempt int DEFAULT 5,

    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,

    status varchar(20) DEFAULT 'active',
    -- active | used | expired | locked

    metadata jsonb DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_otp_target_purpose
ON system_otps (target, purpose);

CREATE INDEX idx_otp_expires_at
ON system_otps (expires_at);

DROP TRIGGER IF EXISTS update_system_otps_updated_at ON system_otps;
CREATE TRIGGER update_system_otps_updated_at
BEFORE UPDATE ON system_otps
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_system_otps_updated_at ON system_otps;
DROP TABLE IF EXISTS system_otps;
-- +goose StatementEnd
