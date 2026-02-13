-- +goose Up
-- +goose StatementBegin
CREATE TABLE system_otp_audits (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    otp_id uuid NOT NULL,
    target varchar(255) NOT NULL,
    purpose varchar(50) NOT NULL,

    attempt_number int NOT NULL,

    result varchar(20) NOT NULL,
    -- success | failed | expired | locked

    failure_reason varchar(50),
    -- invalid_code | expired | max_attempt | already_used

    ip_address inet,
    user_agent text,

    metadata jsonb DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_otp_audit_otp
        FOREIGN KEY (otp_id)
        REFERENCES system_otps(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_system_otp_audits_otp_id
ON system_otp_audits (otp_id);

CREATE INDEX IF NOT EXISTS idx_system_otp_audits_created_at
ON system_otp_audits (created_at);

DROP TRIGGER IF EXISTS update_system_otp_audits_updated_at ON system_otp_audits;
CREATE TRIGGER update_system_otp_audits_updated_at
BEFORE UPDATE ON system_otp_audits
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_system_otp_audits_updated_at ON system_otp_audits;
DROP TABLE IF EXISTS system_otp_audits;
-- +goose StatementEnd
