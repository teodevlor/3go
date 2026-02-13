-- +goose Up
-- +goose StatementBegin

-- =====================
-- system_Sessions table (phiên đăng nhập cụ thể: account + device + app)
-- =====================
CREATE TABLE IF NOT EXISTS system_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    account_app_device_id UUID NOT NULL,
    
    refresh_token_hash VARCHAR(255) NOT NULL, -- Hash của refresh token (không lưu plain text)
    
    expires_at TIMESTAMPTZ NOT NULL,
    is_revoked BOOLEAN NOT NULL DEFAULT false,
    revoked_at TIMESTAMPTZ,
    revoked_reason VARCHAR(255),      -- logout, security, password_reset, ...
    
    last_active_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Thông tin bảo mật bổ sung
    ip_address VARCHAR(45),           -- IPv4 hoặc IPv6
    user_agent TEXT,
    
    metadata JSONB DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    -- Foreign keys
    CONSTRAINT fk_system_sessions_account_app_device
        FOREIGN KEY (account_app_device_id)
        REFERENCES account_app_devices(id)
        ON DELETE CASCADE
);

-- =====================
-- Indexes
-- =====================
CREATE UNIQUE INDEX IF NOT EXISTS idx_system_sessions_refresh_token_hash_unique
ON system_sessions(refresh_token_hash)
WHERE is_revoked = false AND deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_system_sessions_account_app_device_id
ON system_sessions(account_app_device_id);

CREATE INDEX IF NOT EXISTS idx_system_sessions_is_revoked
ON system_sessions(is_revoked);

CREATE INDEX IF NOT EXISTS idx_system_sessions_expires_at
ON system_sessions(expires_at);

CREATE INDEX IF NOT EXISTS idx_system_sessions_last_active_at
ON system_sessions(last_active_at);

CREATE INDEX IF NOT EXISTS idx_system_sessions_deleted_at
ON system_sessions(deleted_at)
WHERE deleted_at IS NOT NULL;

-- =====================
-- Trigger: auto update updated_at
-- =====================
DROP TRIGGER IF EXISTS update_system_sessions_updated_at ON system_sessions;
CREATE TRIGGER update_system_sessions_updated_at
BEFORE UPDATE ON system_sessions
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_system_sessions_updated_at ON system_sessions;

DROP INDEX IF EXISTS idx_system_sessions_deleted_at;
DROP INDEX IF EXISTS idx_system_sessions_last_active_at;
DROP INDEX IF EXISTS idx_system_sessions_expires_at;
DROP INDEX IF EXISTS idx_system_sessions_is_revoked;
DROP INDEX IF EXISTS idx_system_sessions_account_app_device_id;
DROP INDEX IF EXISTS idx_system_sessions_refresh_token_hash_unique;

DROP TABLE IF EXISTS system_sessions;

-- +goose StatementEnd
