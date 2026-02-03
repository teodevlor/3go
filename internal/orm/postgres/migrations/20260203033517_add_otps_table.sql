-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE otps (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    target varchar(255) NOT NULL,
    -- email hoặc phone

    otp_code varchar(10) NOT NULL,

    purpose varchar(50) NOT NULL,
    -- register | reset_password | transaction | withdraw | login | verify_phone

    attempt_count int DEFAULT 0,
    max_attempt int DEFAULT 5,

    expires_at timestamp NOT NULL,
    used_at timestamp NULL,

    status varchar(20) DEFAULT 'active',
    -- active | used | expired | locked

    metadata jsonb DEFAULT '{}'::jsonb,

    created_at timestamp DEFAULT now(),
    updated_at timestamp DEFAULT now()
);

CREATE INDEX idx_otp_target_purpose
ON otps (target, purpose);

CREATE INDEX idx_otp_expires_at
ON otps (expires_at);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS otps;
-- +goose StatementEnd
